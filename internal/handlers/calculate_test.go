package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCalculateHandler(t *testing.T) {
	handler := NewCalculateHandler()
	assert.NotNil(t, handler)
}

func TestCalculateHandler_Handle_MethodRouting(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "should reject GET method",
			method:         http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "should reject PUT method",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "should reject DELETE method",
			method:         http.MethodDelete,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "should reject PATCH method",
			method:         http.MethodPatch,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewCalculateHandler()
			req := httptest.NewRequest(tt.method, "/calculate", nil)
			w := httptest.NewRecorder()

			handler.Handle(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var errorResponse map[string]string
			err := json.NewDecoder(w.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedError, errorResponse["error"])
		})
	}
}

func TestCalculateHandler_HandlePost(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        map[string]interface{}
		expectedStatus     int
		expectedOrder      int
		expectedTotalItems int
		shouldHaveError    bool
		expectedError      string
	}{
		{
			name: "should calculate order from JSON body",
			requestBody: map[string]interface{}{
				"order":      501,
				"pack_sizes": []int{250, 500, 1000},
			},
			expectedStatus:     http.StatusOK,
			expectedOrder:      501,
			expectedTotalItems: 750,
			shouldHaveError:    false,
		},
		{
			name: "should handle order zero",
			requestBody: map[string]interface{}{
				"order":      0,
				"pack_sizes": []int{250, 500, 1000},
			},
			expectedStatus:     http.StatusOK,
			expectedOrder:      0,
			expectedTotalItems: 0,
			shouldHaveError:    false,
		},
		{
			name: "should reject negative order",
			requestBody: map[string]interface{}{
				"order":      -100,
				"pack_sizes": []int{250, 500, 1000},
			},
			expectedStatus:  http.StatusBadRequest,
			shouldHaveError: true,
			expectedError:   "Order must be positive",
		},
		{
			name: "should handle large order",
			requestBody: map[string]interface{}{
				"order":      12001,
				"pack_sizes": []int{250, 500, 1000, 2000, 5000},
			},
			expectedStatus:     http.StatusOK,
			expectedOrder:      12001,
			expectedTotalItems: 12250,
			shouldHaveError:    false,
		},
		{
			name: "should handle exact match order",
			requestBody: map[string]interface{}{
				"order":      1000,
				"pack_sizes": []int{250, 500, 1000},
			},
			expectedStatus:     http.StatusOK,
			expectedOrder:      1000,
			expectedTotalItems: 1000,
			shouldHaveError:    false,
		},
		{
			name: "should reject empty pack sizes",
			requestBody: map[string]interface{}{
				"order":      501,
				"pack_sizes": []int{},
			},
			expectedStatus:  http.StatusBadRequest,
			shouldHaveError: true,
			expectedError:   "Pack sizes cannot be empty",
		},
		{
			name: "should reject negative pack size",
			requestBody: map[string]interface{}{
				"order":      501,
				"pack_sizes": []int{250, -500, 1000},
			},
			expectedStatus:  http.StatusBadRequest,
			shouldHaveError: true,
			expectedError:   "All pack sizes must be positive",
		},
		{
			name: "should reject zero pack size",
			requestBody: map[string]interface{}{
				"order":      501,
				"pack_sizes": []int{250, 0, 1000},
			},
			expectedStatus:  http.StatusBadRequest,
			shouldHaveError: true,
			expectedError:   "All pack sizes must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewCalculateHandler()

			bodyBytes, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Handle(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.shouldHaveError {
				var errorResponse map[string]string
				err := json.NewDecoder(w.Body).Decode(&errorResponse)
				require.NoError(t, err)

				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, errorResponse["error"])
				}
			} else {
				var response map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)

				assert.Equal(t, float64(tt.expectedOrder), response["order"])
				assert.Equal(t, float64(tt.expectedTotalItems), response["total_items"])
				assert.Contains(t, response, "packs")
				assert.Contains(t, response, "pack_sizes")
				assert.Contains(t, response, "surplus")
				assert.Contains(t, response, "total_packs")
			}
		})
	}
}

func TestCalculateHandler_HandlePost_InvalidJSON(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
	}{
		{
			name:        "should reject malformed JSON",
			requestBody: "{invalid json}",
		},
		{
			name:        "should reject empty body",
			requestBody: "",
		},
		{
			name:        "should reject non-JSON content",
			requestBody: "plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewCalculateHandler()

			req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Handle(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var errorResponse map[string]string
			err := json.NewDecoder(w.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.Equal(t, "Invalid request body", errorResponse["error"])
		})
	}
}

func TestCalculateHandler_ResponseFormat(t *testing.T) {
	t.Run("should return all expected fields in response", func(t *testing.T) {
		handler := NewCalculateHandler()

		bodyBytes, _ := json.Marshal(map[string]interface{}{
			"order":      501,
			"pack_sizes": []int{250, 500, 1000},
		})
		req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		// Verify all expected fields are present
		assert.Contains(t, response, "order")
		assert.Contains(t, response, "total_items")
		assert.Contains(t, response, "packs")
		assert.Contains(t, response, "pack_sizes")
		assert.Contains(t, response, "surplus")
		assert.Contains(t, response, "total_packs")

		// Verify response structure
		assert.Equal(t, float64(501), response["order"])
		assert.Equal(t, float64(750), response["total_items"])
		assert.IsType(t, map[string]interface{}{}, response["packs"])
		assert.IsType(t, []interface{}{}, response["pack_sizes"])
		assert.Equal(t, float64(249), response["surplus"])
		assert.Equal(t, float64(2), response["total_packs"])
	})
}
