package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go-beresin/internal/transport/router"
	"go-beresin/pkg/database"

	// Import Swagger generated docs
	_ "go-beresin/docs"
)

// @title Go Beresin API
// @version 1.0.0
// @description High-performance backend API skeleton for Go Beresin featuring Fiber, Redis-based Rate Limiting, JWT auth with RBAC, and Multipart file uploading.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@go-beresin.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <your_token>" to authenticate.
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] No .env file found, using system environment variables")
	}

	// 1. Configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// 2. Initialize Redis Client
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = "!Abcd1234"
	}
	log.Printf("[INFO] Connecting to Redis at %s...", redisAddr)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,  // Use default DB
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("[WARNING] Could not connect to Redis: %v. Rate limiter will be bypassed (fail-open).", err)
		// Set to nil so middleware knows to bypass rate limiting
		rdb = nil
	} else {
		log.Println("[INFO] Successfully connected to Redis.")
	}

	// 3. Initialize PostgreSQL connection pool
	dbPool, err := database.InitPool(ctx)
	if err != nil {
		log.Fatalf("[FATAL] Could not initialize PostgreSQL pool: %v", err)
	}

	// 4. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName:      "Go Beresin API v1.0.0",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// 5. Setup Routing
	router.SetupRoutes(app, rdb, dbPool)

	// 6. Channel to catch errors during startup
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		log.Printf("[INFO] Server is starting on port %s...", port)
		if err := app.Listen(":" + port); err != nil {
			serverErrors <- err
		}
	}()

	// 7. Graceful Shutdown
	// Channel to listen for terminate signals from OS
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("[FATAL] Server failed to start: %v", err)

	case sig := <-shutdownChannel:
		log.Printf("[INFO] Shutdown signal received (%v). Starting graceful shutdown...", sig)

		// Trigger Fiber graceful shutdown
		log.Println("[INFO] Shutting down HTTP server...")
		if err := app.ShutdownWithTimeout(15 * time.Second); err != nil {
			log.Printf("[WARNING] Graceful server shutdown failed: %v, forcing close", err)
		}

		// Close PostgreSQL pool
		log.Println("[INFO] Closing PostgreSQL connection pool...")
		dbPool.Close()

		// Close Redis connection
		if rdb != nil {
			log.Println("[INFO] Closing Redis client...")
			if err := rdb.Close(); err != nil {
				log.Printf("[WARNING] Error closing Redis client: %v", err)
			}
		}

		log.Println("[INFO] Server stopped gracefully. Clean exit.")
	}
}
