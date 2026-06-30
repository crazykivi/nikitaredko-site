package main

import (
	"flag"
	"log"
	"os"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	cacheManager := cache.New()
	articleHandler := handlers.NewArticleHandler(cacheManager)

	api := r.Group("/api")
	{
		api.GET("/collections", articleHandler.ListCollections)
		api.GET("/articles", articleHandler.ListArticles)
		api.GET("/articles/structured", articleHandler.ListArticlesStructured)
		api.GET("/articles/:id", articleHandler.GetArticle)

		// CACHE
		api.POST("/webhook/outline", cacheManager.WebhookHandler)
		api.GET("/cache/health", cacheManager.HealthCheck)
	}

	if *serveStatic {
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
		log.Println("[Static] Static serving disabled via SERVE_STATIC=false")
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
