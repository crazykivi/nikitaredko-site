package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"nikitaredko-backend/cache"
	"nikitaredko-backend/handlers"
)

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

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
