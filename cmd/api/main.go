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

	"github.com/luisfernandomoraes/order-packing-api/internal/config"
	"github.com/luisfernandomoraes/order-packing-api/internal/domain"
	"github.com/luisfernandomoraes/order-packing-api/internal/server"
)

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
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}
