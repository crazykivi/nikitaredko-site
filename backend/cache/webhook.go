package cache

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

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
	if expectedSecret != "" {
		secret := ctx.Query("secret")
		if secret != expectedSecret {
			sigHeader := ctx.GetHeader("X-Outline-Signature")
			if sigHeader != expectedSecret && secret != expectedSecret {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid secret"})
				return
			}
		}
	}

	var event OutlineWebhookEvent
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	log.Printf("[Webhook] Received: %s %s (model: %s, action: %s)",
		event.ID, event.Name, event.Model, event.Action)

	c.handleEvent(event)

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "event": event.ID})
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
