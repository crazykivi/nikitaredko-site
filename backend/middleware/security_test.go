package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, "<html></html>")
	})
	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func TestSecurityHeadersSet(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	tests := []struct {
		name   string
		header string
		want   string
	}{
		{"nosniff", "X-Content-Type-Options", "nosniff"},
		{"frame deny", "X-Frame-Options", "DENY"},
		{"referrer policy", "Referrer-Policy", "strict-origin-when-cross-origin"},
		{"coop", "Cross-Origin-Opener-Policy", "same-origin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := w.Header().Get(tt.header); got != tt.want {
				t.Errorf("%s = %q, want %q", tt.header, got, tt.want)
			}
		})
	}

	if got := w.Header().Get("Permissions-Policy"); got == "" {
		t.Error("Permissions-Policy is empty")
	}
}

func TestCSPHeader(t *testing.T) {
	r := newTestRouter()

	t.Run("present on html response", func(t *testing.T) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

		csp := w.Header().Get("Content-Security-Policy")
		if csp == "" {
			t.Fatal("Content-Security-Policy missing on html response")
		}

		for _, directive := range []string{
			"default-src 'self'",
			"frame-ancestors 'none'",
			"object-src 'none'",
			"base-uri 'self'",
			"script-src 'self' https://giscus.app",
			"frame-src https://giscus.app",
			"font-src 'self' https://fonts.gstatic.com",
		} {
			if !strings.Contains(csp, directive) {
				t.Errorf("CSP missing directive %q\nCSP: %s", directive, csp)
			}
		}
	})

	t.Run("present on json response too", func(t *testing.T) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/json", nil))

		if w.Header().Get("Content-Security-Policy") == "" {
			t.Error("Content-Security-Policy missing on json response")
		}
	})
}

func TestHSTSOnlyOverHTTPS(t *testing.T) {
	r := newTestRouter()

	t.Run("absent on plain http", func(t *testing.T) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

		if got := w.Header().Get("Strict-Transport-Security"); got != "" {
			t.Errorf("HSTS must be absent on http, got %q", got)
		}
	})

	t.Run("present when behind https proxy", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Forwarded-Proto", "https")
		r.ServeHTTP(w, req)

		want := "max-age=63072000; includeSubDomains; preload"
		if got := w.Header().Get("Strict-Transport-Security"); got != want {
			t.Errorf("HSTS = %q, want %q", got, want)
		}
	})

	t.Run("present on direct TLS request", func(t *testing.T) {
		srv := httptest.NewTLSServer(r)
		defer srv.Close()

		client := srv.Client()
		resp, err := client.Get(srv.URL + "/test")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if got := resp.Header.Get("Strict-Transport-Security"); got == "" {
			t.Error("HSTS missing on direct TLS request")
		}
	})
}
