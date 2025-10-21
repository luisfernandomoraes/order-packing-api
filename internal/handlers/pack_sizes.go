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

// PackSizesRequest represents the request body for updating pack sizes
type PackSizesRequest struct {
	PackSizes []int `json:"pack_sizes" example:"100,250,500,1000"`
}

// PackSizesResponse represents the response from pack sizes endpoints
type PackSizesResponse struct {
	PackSizes []int `json:"pack_sizes" example:"250,500,1000,2000,5000"`
}

// PackSizesUpdateResponse represents the response from update pack sizes endpoint
type PackSizesUpdateResponse struct {
	Message   string `json:"message" example:"Pack sizes updated successfully"`
	PackSizes []int  `json:"pack_sizes" example:"250,500,1000,2000,5000"`
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

// handleGet godoc
// @Summary Get current package sizes
// @Description Returns the currently configured package sizes
// @Tags pack-sizes
// @Produce json
// @Success 200 {object} PackSizesResponse
// @Router /api/pack-sizes [get]
func (h *PackSizesHandler) handleGet(w http.ResponseWriter, _ *http.Request) {
	responseData := PackSizesResponse{
		PackSizes: h.calculator.GetPackSizes(),
	}
	response.JSON(w, http.StatusOK, responseData)
}

// handlePost godoc
// @Summary Update package sizes
// @Description Updates the available package sizes used for calculations
// @Tags pack-sizes
// @Accept json
// @Produce json
// @Param request body PackSizesRequest true "New pack sizes"
// @Success 200 {object} PackSizesUpdateResponse
// @Failure 400 {object} map[string]string "Bad Request - Empty array or non-positive values"
// @Router /api/pack-sizes [post]
func (h *PackSizesHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var req PackSizesRequest

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

	responseData := PackSizesUpdateResponse{
		Message:   "Pack sizes updated successfully",
		PackSizes: h.calculator.GetPackSizes(),
	}

	response.JSON(w, http.StatusOK, responseData)
}
