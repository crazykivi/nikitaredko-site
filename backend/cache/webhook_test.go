package cache

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWebhookHandler_ValidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Unsetenv("OUTLINE_WEBHOOK_SECRET")

	c := New()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/webhook", c.WebhookHandler)

	payload := `{"id":"evt-1","name":"documents.update","model":"document","action":"update","data":{},"timestamp":"2025-01-01T00:00:00Z"}`
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

func TestWebhookHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Unsetenv("OUTLINE_WEBHOOK_SECRET")

	c := New()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/webhook", c.WebhookHandler)

	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestWebhookHandler_WithSecret_Valid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("OUTLINE_WEBHOOK_SECRET", "mysecret")
	defer os.Unsetenv("OUTLINE_WEBHOOK_SECRET")

	c := New()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/webhook", c.WebhookHandler)

	payload := `{"id":"evt-1","name":"test","model":"document","action":"update","data":{},"timestamp":""}`
	req := httptest.NewRequest("POST", "/webhook?secret=mysecret", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

func TestWebhookHandler_WithSecret_Invalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("OUTLINE_WEBHOOK_SECRET", "mysecret")
	defer os.Unsetenv("OUTLINE_WEBHOOK_SECRET")

	c := New()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/webhook", c.WebhookHandler)

	payload := `{"id":"evt-1","name":"test","model":"document","action":"update","data":{},"timestamp":""}`
	req := httptest.NewRequest("POST", "/webhook?secret=wrongsecret", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403", w.Code)
	}
}

func TestWebhookHandler_FlushesCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Unsetenv("OUTLINE_WEBHOOK_SECRET")

	c := New()
	c.Set("articles_list", []string{"a", "b"})

	if _, found := c.Get("articles_list"); !found {
		t.Fatal("cache should have data before webhook")
	}

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/webhook", c.WebhookHandler)

	payload := `{"id":"evt-1","name":"documents.update","model":"document","action":"update","data":{},"timestamp":""}`
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if _, found := c.Get("articles_list"); found {
		t.Error("cache should be flushed after webhook")
	}
}

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c := New()
	c.Set("test", "value")

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/health", c.HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte(`"enabled":true`)) {
		t.Error("health check should show enabled=true")
	}
}
