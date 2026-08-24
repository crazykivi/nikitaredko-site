package cache

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	json "github.com/goccy/go-json"

	"github.com/gin-gonic/gin"
)

// максимальный возраст подписи для защиты от replay-атак
const maxSignatureAge = 5 * time.Minute

type OutlineWebhookEvent struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Model     string          `json:"model"`  // "document", "collection" и т.д.
	Action    string          `json:"action"` // "create", "update", "delete", "publish", "unpublish"
	Data      json.RawMessage `json:"data"`
	Timestamp string          `json:"timestamp"`
}
type DocumentData struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	CollectionID string `json:"collectionId"`
}

func (c *Cache) WebhookHandler(ctx *gin.Context) {
	expectedSecret := os.Getenv("OUTLINE_WEBHOOK_SECRET")

	const maxWebhookBodySize = 1 << 20 // 1 MB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxWebhookBodySize)
	rawBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Payload too large"})
		return
	}

	if expectedSecret != "" && !validWebhookRequest(ctx, rawBody, expectedSecret) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid signature"})
		return
	}

	var event OutlineWebhookEvent
	if err := json.Unmarshal(rawBody, &event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	log.Printf("[Webhook] Received: %s %s (model: %s, action: %s)",
		event.ID, event.Name, event.Model, event.Action)
	c.handleEvent(event)
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "event": event.ID})
}

// Проверка подлинности запроса:
//   - ?secret=... в query — для ручной проверки (curl, тесты);
//   - Outline-Signature — штатная подпись Outline:
//     HMAC-SHA256(secret, "<t>.<body>") в формате "t=<millis>,s=<hex>".
func validWebhookRequest(ctx *gin.Context, rawBody []byte, expectedSecret string) bool {
	if subtle.ConstantTimeCompare([]byte(ctx.Query("secret")), []byte(expectedSecret)) == 1 {
		return true
	}

	sigHeader := ctx.GetHeader("Outline-Signature")
	if sigHeader == "" {
		log.Printf("[Webhook] Rejected: missing Outline-Signature header")
		return false
	}
	if !verifyOutlineSignature(sigHeader, rawBody, expectedSecret) {
		log.Printf("[Webhook] Rejected: invalid Outline-Signature")
		return false
	}
	return true
}

func verifyOutlineSignature(header string, rawBody []byte, secret string) bool {
	var ts, sig string
	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		switch {
		case strings.HasPrefix(part, "t="):
			ts = strings.TrimPrefix(part, "t=")
		case strings.HasPrefix(part, "s="):
			sig = strings.TrimPrefix(part, "s=")
		}
	}
	if ts == "" || sig == "" {
		return false
	}

	tsMillis, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return false
	}
	if age := time.Since(time.UnixMilli(tsMillis)); age < -maxSignatureAge || age > maxSignatureAge {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts + "."))
	mac.Write(rawBody)
	expected := hex.EncodeToString(mac.Sum(nil))

	return subtle.ConstantTimeCompare([]byte(expected), []byte(sig)) == 1
}

func (c *Cache) handleEvent(event OutlineWebhookEvent) {
	switch event.Model {
	case "document", "collection":
		switch event.Action {
		case "create", "update", "delete", "publish", "unpublish", "archive", "unarchive":
			log.Printf("[Webhook] Invalidating cache for %s %s", event.Model, event.Action)
			c.Flush()
			return
		}
	}
	log.Printf("[Webhook] Invalidating cache (default)")
	c.Flush()
}
func (c *Cache) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"enabled": c.IsEnabled(),
		"ttl":     c.ttl.String(),
		"keys":    len(c.Keys()),
	})
}
