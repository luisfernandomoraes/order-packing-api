package handlers

import (
	"net/http"

	"github.com/luisfernandomoraes/order-packing-api/internal/domain"
	"github.com/luisfernandomoraes/order-packing-api/internal/response"
)

type PackSizesHandler struct {
	calculator *domain.PackCalculator
}

func NewPackSizesHandler(calculator *domain.PackCalculator) *PackSizesHandler {
	return &PackSizesHandler{
		calculator: calculator,
	}
}

func (h *PackSizesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *PackSizesHandler) handleGet(w http.ResponseWriter, _ *http.Request) {
	responseData := map[string]interface{}{
		"pack_sizes": h.calculator.GetPackSizes(),
	}
	response.JSON(w, http.StatusOK, responseData)
}

func (h *PackSizesHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PackSizes []int `json:"pack_sizes"`
	}

	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.PackSizes) == 0 {
		response.Error(w, http.StatusBadRequest, "Pack sizes cannot be empty")
		return
	}

	for _, size := range req.PackSizes {
		if size <= 0 {
			response.Error(w, http.StatusBadRequest, "All pack sizes must be positive")
			return
		}
	}

	h.calculator.UpdatePackSizes(req.PackSizes)

	responseData := map[string]interface{}{
		"message":    "Pack sizes updated successfully",
		"pack_sizes": h.calculator.GetPackSizes(),
	}

	response.JSON(w, http.StatusOK, responseData)
}
