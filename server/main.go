package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-url-shortner/handlers"
	"go-url-shortner/services"
	"go-url-shortner/storage"
	"go-url-shortner/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
"go-url-shortner/middleware"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load configuration
	cfg := utils.Load()

	// Initialize Redis storage
	storage, err := storage.NewRedisStorage(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to initialize Redis storage: %v", err)
	}
	defer storage.Close()

	// Initialize AI service (optional)
	var aiService services.AISlugServiceInterface
	if cfg.OpenAIAPIKey != "" {
		aiService = services.NewAISlugService(cfg.OpenAIAPIKey)
		log.Println("AI slug generation enabled")
	} else {
		log.Println("AI slug generation disabled - no API key provided")
	}

	// Initialize URL service
	urlService := services.NewURLService(storage, aiService, cfg.ServerHost, cfg.ServerPort)

	// Initialize handlers
	urlHandler := handlers.NewURLHandler(urlService)

	// Setup Gin router
	router := gin.Default()
	// Allow frontend origin
	router.Use(middleware.CORSMiddleware("http://localhost:3000"))

	// Routes
	router.POST("/api/urls", urlHandler.CreateShortURL)
	router.GET("/:shortCode", urlHandler.RedirectToURL)
	router.GET("/health", urlHandler.HealthCheck)

	// Create HTTP server
       server := &http.Server{
	       Addr:    "0.0.0.0:" + cfg.ServerPort,
	       Handler: router,
       }

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s:%s", cfg.ServerHost, cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
