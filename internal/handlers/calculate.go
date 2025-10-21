package handlers

import (
	"net/http"

	"github.com/luisfernandomoraes/order-packing-api/internal/domain"
	"github.com/luisfernandomoraes/order-packing-api/internal/response"
)

type CalculateHandler struct {
	calculator *domain.PackCalculator
}

func NewCalculateHandler(calculator *domain.PackCalculator) *CalculateHandler {
	return &CalculateHandler{
		calculator: calculator,
	}
}

func (h *CalculateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Order int `json:"order"`
	}

	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Order < 0 {
		response.Error(w, http.StatusBadRequest, "Order must be positive")
		return
	}

	result := h.calculator.Calculate(req.Order)

	responseData := map[string]interface{}{
		"order":       result.Order,
		"total_items": result.TotalItems,
		"packs":       result.Packs,
		"pack_sizes":  result.PackSizes,
		"surplus":     result.GetSurplus(),
		"total_packs": result.GetTotalPackCount(),
	}

	response.JSON(w, http.StatusOK, responseData)
}
