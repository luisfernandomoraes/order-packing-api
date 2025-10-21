package handlers

import (
	"net/http"

	"github.com/luisfernandomoraes/order-packing-api/internal/domain"
	"github.com/luisfernandomoraes/order-packing-api/internal/response"
)

type CalculateHandler struct{}

func NewCalculateHandler() *CalculateHandler {
	return &CalculateHandler{}
}

// CalculateRequest represents the request body for calculate endpoint
type CalculateRequest struct {
	Order     int   `json:"order" example:"501" minimum:"0"`
	PackSizes []int `json:"pack_sizes" example:"250,500,1000,2000,5000"`
}

// CalculateResponse represents the response from calculate endpoint
type CalculateResponse struct {
	Order      int         `json:"order" example:"501"`
	TotalItems int         `json:"total_items" example:"750"`
	Packs      map[int]int `json:"packs" example:"250:1,500:1"`
	PackSizes  []int       `json:"pack_sizes" example:"250,500,1000,2000,5000"`
	Surplus    int         `json:"surplus" example:"249"`
	TotalPacks int         `json:"total_packs" example:"2"`
}

// Handle godoc
// @Summary Calculate optimal package combination
// @Description Calculates the best package combination to fulfill an order, minimizing items shipped and number of packages
// @Tags calculate
// @Accept json
// @Produce json
// @Param request body CalculateRequest true "Order quantity and package sizes"
// @Success 200 {object} CalculateResponse
// @Failure 400 {object} map[string]string "Bad Request - Invalid order, negative value, or invalid pack sizes"
// @Failure 405 {object} map[string]string "Method Not Allowed"
// @Router /api/calculate [post]
func (h *CalculateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CalculateRequest

	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Order < 0 {
		response.Error(w, http.StatusBadRequest, "Order must be positive")
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

	calculator := domain.NewPackCalculator(req.PackSizes)
	result := calculator.Calculate(req.Order)

	responseData := CalculateResponse{
		Order:      result.Order,
		TotalItems: result.TotalItems,
		Packs:      result.Packs,
		PackSizes:  result.PackSizes,
		Surplus:    result.GetSurplus(),
		TotalPacks: result.GetTotalPackCount(),
	}

	response.JSON(w, http.StatusOK, responseData)
}
