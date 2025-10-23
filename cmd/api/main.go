// Package main starts the HTTP server for the Order Packing API.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/luisfernandomoraes/order-packing-api/docs" // Swagger docs
	"github.com/luisfernandomoraes/order-packing-api/internal/config"
	"github.com/luisfernandomoraes/order-packing-api/internal/domain"
	"github.com/luisfernandomoraes/order-packing-api/internal/server"
)

// @title Order Packing Calculator API
// @version 1.0
// @description A REST API to calculate the optimal package combination to fulfill orders, minimizing items shipped and number of packages.
// @termsOfService http://swagger.io/terms/

// @contact.name Luis Fernando Moraes
// @contact.email luisfernandomoraes@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Initialize domain services
	calculator := domain.NewPackCalculator(cfg.DefaultPackSizes)

	// Create and start server
	srv := server.New(cfg, calculator)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.Port)
		log.Printf("📦 Default pack sizes: %v", cfg.DefaultPackSizes)
		log.Printf("🌐 API: http://localhost:%s/api", cfg.Port)
		log.Printf("📚 Swagger docs: http://localhost:%s/swagger/index.html", cfg.Port)
		log.Printf("💚 Health: http://localhost:%s/health", cfg.Port)
		log.Printf("🎨 UI: http://localhost:%s", cfg.Port)

		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
		cancel()
		os.Exit(1)
	}

	cancel()
	log.Println("✅ Server stopped gracefully")
}
