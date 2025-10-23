// Package server implements the HTTP server for the order packing API.
package server

import (
	"context"
	"net/http"

	"github.com/luisfernandomoraes/order-packing-api/internal/config"
	"github.com/luisfernandomoraes/order-packing-api/internal/domain"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	calculator *domain.PackCalculator
	config     config.Config
}

// New creates a new Server instance
func New(cfg config.Config, calculator *domain.PackCalculator) *Server {
	srv := &Server{
		calculator: calculator,
		config:     cfg,
	}

	srv.httpServer = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      srv.setupRoutes(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return srv
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
