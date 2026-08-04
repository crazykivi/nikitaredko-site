package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"nikitaredko-backend/cache"
	"nikitaredko-backend/handlers"
)

var serveStatic = flag.Bool("s", false, "Serve static files from ./dist")

func main() {
	loadConfig()

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
	aboutHandler := handlers.NewAboutHandler(articleHandler, cacheManager)

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
		api.GET("/about", aboutHandler.GetAbout)

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
			r.NoRoute(articleHandler.ServeFrontend)
			log.Println("[Static] Serving frontend from ./dist")
		} else {
			log.Println("[Static] No ./dist folder found, API-only mode")
		}
	} else {
		log.Println("[Static] Static serving disabled")
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s (mode: %s)", port, gin.Mode())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}

func loadConfig() {
	isProduction := os.Getenv("GIN_MODE") == "release" || os.Getenv("PRODUCTION") == "true"

	if isProduction {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		return
	}

	if err := godotenv.Load(); err != nil {
		log.Println("[Config] .env not found, using system environment variables")
	}
}
