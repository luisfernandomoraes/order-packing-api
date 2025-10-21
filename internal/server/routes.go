package server

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/luisfernandomoraes/order-packing-api/internal/handlers"
	"github.com/luisfernandomoraes/order-packing-api/internal/middleware"
)

func (s *Server) setupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Create handlers
	calculateHandler := handlers.NewCalculateHandler(s.calculator)
	packSizesHandler := handlers.NewPackSizesHandler(s.calculator)
	healthHandler := handlers.NewHealthHandler()

	// Swagger documentation
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// API routes with middleware
	mux.HandleFunc("/api/calculate", middleware.Chain(
		calculateHandler.Handle,
		middleware.CORS,
		middleware.Logging,
		middleware.Recovery,
	))

	mux.HandleFunc("/api/pack-sizes", middleware.Chain(
		packSizesHandler.Handle,
		middleware.CORS,
		middleware.Logging,
		middleware.Recovery,
	))

	mux.HandleFunc("/health", middleware.Chain(
		healthHandler.Handle,
		middleware.CORS,
		middleware.Recovery,
	))

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	return mux
}
