package handlers

import (
	"net/http"

	"github.com/luisfernandomoraes/order-packing-api/internal/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	responseData := map[string]string{
		"status": "healthy",
		"app":    "Order Packing Calculator API",
	}
	response.JSON(w, http.StatusOK, responseData)
}
