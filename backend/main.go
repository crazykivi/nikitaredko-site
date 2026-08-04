package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"nikitaredko-backend/cache"
	"nikitaredko-backend/handlers"
)

var serveStatic = flag.Bool("s", false, "Serve static files from ./dist")

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	if os.Getenv("GIN_MODE") == "release" || os.Getenv("PRODUCTION") == "true" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	_ = r.SetTrustedProxies(nil)
	if trustedProxies := os.Getenv("TRUSTED_PROXIES"); trustedProxies != "" {
		proxies := strings.Split(trustedProxies, ",")
		for i := range proxies {
			proxies[i] = strings.TrimSpace(proxies[i])
		}
		_ = r.SetTrustedProxies(proxies)
	}

	corsOrigins := []string{"http://localhost:5173", "http://localhost:3000"}
	if envOrigins := os.Getenv("ALLOW_CORS"); envOrigins != "" {
		corsOrigins = []string{}
		for _, origin := range strings.Split(envOrigins, ",") {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				corsOrigins = append(corsOrigins, trimmed)
			}
		}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	cacheManager := cache.New()
	articleHandler := handlers.NewArticleHandler(cacheManager)
	usesHandler := handlers.NewUsesHandler(articleHandler, cacheManager)

	api := r.Group("/api")
	{
		api.GET("/collections", articleHandler.ListCollections)
		api.GET("/articles", articleHandler.ListArticles)
		api.GET("/articles/structured", articleHandler.ListArticlesStructured)
		api.GET("/articles/:id", articleHandler.GetArticle)
		api.GET("/articles/search", articleHandler.SearchArticles)
		api.GET("/articles/feed", articleHandler.GetArticlesFeed)
		api.GET("/rss.xml", articleHandler.GetRSS)
		api.GET("/sitemap.xml", articleHandler.GetSitemap)
		api.GET("/uses", usesHandler.GetUses)

		// ПРОКСИ ДЛЯ КАРТИНОК
		api.GET("/attachments.redirect", articleHandler.ProxyOutlineAttachment)

		// CACHE
		api.POST("/webhook/outline", cacheManager.WebhookHandler)
		api.GET("/cache/health", cacheManager.HealthCheck)
	}

	shouldServeStatic := *serveStatic
	if !shouldServeStatic {
		if envVal := os.Getenv("SERVE_STATIC"); envVal == "true" {
			shouldServeStatic = true
		}
	}

	if shouldServeStatic {
		if _, err := os.Stat("./dist"); err == nil {
			r.Static("/assets", "./dist/assets")
			r.StaticFile("/favicon.svg", "./dist/favicon.svg")

			r.NoRoute(func(c *gin.Context) {
				c.File("./dist/index.html")
			})
			log.Println("[Static] Serving frontend from ./dist")
		} else {
			log.Println("[Static] No ./dist folder found, API-only mode")
		}
	} else {
		log.Println("[Static] Static serving disabled")
	}

	log.Printf("Server starting on port %s (mode: %s)", port, gin.Mode())
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
