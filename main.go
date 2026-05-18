package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Response represents a standard JSON response structure.
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Determine the port from environment variables or use a default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new ServeMux for routing
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/health", handleHealthCheck)

	// Create the HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      loggerMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Channel to listen for errors during server startup
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("[INFO] Server is starting on port %s...", port)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for interrupt or terminate signals from the OS
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGTERM)

	// Block until a signal or a server error is received
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server failed to start: %v", err)
		}

	case sig := <-shutdownChannel:
		log.Printf("[INFO] Shutdown signal received (%v). Starting graceful shutdown...", sig)

		// Create a context with a timeout for the shutdown process
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Attempt to gracefully shut down the server
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("[ERROR] Graceful shutdown failed: %v", err)
			log.Printf("[INFO] Forcing server shutdown...")
			if err := server.Close(); err != nil {
				log.Fatalf("[FATAL] Could not force close server: %v", err)
			}
		}
		log.Println("[INFO] Server stopped gracefully. Clean exit.")
	}
}

// loggerMiddleware logs incoming HTTP requests details.
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[HTTP] %s %s %s - %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// handleHome handles requests to the root path.
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		respondWithJSON(w, http.StatusNotFound, Response{
			Status:  "error",
			Message: "Resource not found",
		})
		return
	}

	data := map[string]string{
		"project":     "Go Beresin",
		"version":     "1.0.0",
		"description": "Project Go ini telah berhasil diinisialisasi dan siap digunakan!",
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Welcome to Go Beresin API!",
		Data:    data,
	})
}

// handleHealthCheck handles basic server health monitoring.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Server is healthy",
	})
}

// respondWithJSON helper writes a JSON response to the client.
func respondWithJSON(w http.ResponseWriter, statusCode int, payload Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("[ERROR] Failed to encode JSON response: %v", err)
	}
}
