package handlers

import (
	"net/http"

	"github.com/luisfernandomoraes/order-packing-api/internal/response"
)

// HealthHandler handles the /health endpoint
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Handle godoc
// @Summary Health check endpoint
// @Description Returns the health status of the API
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "status: healthy"
// @Router /health [get]
func (h *HealthHandler) Handle(w http.ResponseWriter, _ *http.Request) {
	responseData := map[string]string{
		"status": "healthy",
		"app":    "Order Packing Calculator API",
	}
	response.JSON(w, http.StatusOK, responseData)
}
