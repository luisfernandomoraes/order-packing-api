package handlers

import (
	"net/http"
	"strconv"

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
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CalculateHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	orderStr := r.URL.Query().Get("order")
	order, err := strconv.Atoi(orderStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid order parameter")
		return
	}

	h.calculate(w, order)
}

func (h *CalculateHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Order int `json:"order"`
	}

	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	h.calculate(w, req.Order)
}

func (h *CalculateHandler) calculate(w http.ResponseWriter, order int) {
	if order < 0 {
		response.Error(w, http.StatusBadRequest, "Order must be positive")
		return
	}

	result := h.calculator.Calculate(order)

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
