package server

import (
	"context"
	"net/http"

	"github.com/luisfernandomoraes/order-packing-api/internal/config"
)

type Server struct {
	httpServer *http.Server
	config     config.Config
}

func New(cfg config.Config) *Server {
	srv := &Server{
		config: cfg,
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

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
