package middleware

import (
	"github.com/gin-gonic/gin"
)

// Content-Security-Policy для SPA-фронтенда:
//   - 'self'                    	— бандлы Vite, API, проксированные аттачменты
//   - https://giscus.app        	— скрипт и iframe комментариев
//   - fonts.googleapis.com      	— CSS Google Fonts
//   - fonts.gstatic.com         	— файлы шрифтов
//   - 'unsafe-inline' в style-src 	— inline-стили, которые ставит Vue (SFC-биндинги)
const contentSecurityPolicy = "default-src 'self'; " +
	"script-src 'self' https://giscus.app; " +
	"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
	"img-src 'self' data: blob:; " +
	"font-src 'self' https://fonts.gstatic.com; " +
	"connect-src 'self'; " +
	"frame-src https://giscus.app; " +
	"worker-src 'self'; " +
	"manifest-src 'self'; " +
	"object-src 'none'; " +
	"base-uri 'self'; " +
	"form-action 'none'; " +
	"frame-ancestors 'none'; " +
	"upgrade-insecure-requests"

const permissionsPolicy = "camera=(), microphone=(), geolocation=(), " +
	"payment=(), usb=(), bluetooth=(), accelerometer=(), gyroscope=(), " +
	"magnetometer=(), fullscreen=(self)"

const hstsValue = "max-age=63072000; includeSubDomains; preload"

func isHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	proto := c.GetHeader("X-Forwarded-Proto")
	return proto == "https" || proto == "HTTPS"
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()

		h.Set("Content-Security-Policy", contentSecurityPolicy)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", permissionsPolicy)
		h.Set("Cross-Origin-Opener-Policy", "same-origin")

		if isHTTPS(c) {
			h.Set("Strict-Transport-Security", hstsValue)
		}

		c.Next()
	}
}
