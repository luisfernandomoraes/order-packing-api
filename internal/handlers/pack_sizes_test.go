package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luisfernandomoraes/order-packing-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPackSizesHandler(t *testing.T) {
	calculator := domain.NewPackCalculator([]int{250, 500, 1000})
	handler := NewPackSizesHandler(calculator)
	assert.NotNil(t, handler)
}

func TestPackSizesHandler_Handle_MethodRouting(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "should handle GET method",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "should handle POST method",
			method:         http.MethodPost,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "should reject PUT method",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "should reject DELETE method",
			method:         http.MethodDelete,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "should reject PATCH method",
			method:         http.MethodPatch,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := domain.NewPackCalculator([]int{250, 500, 1000})
			handler := NewPackSizesHandler(calculator)
			req := httptest.NewRequest(tt.method, "/pack-sizes", nil)
			w := httptest.NewRecorder()

			handler.Handle(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestPackSizesHandler_HandleGet(t *testing.T) {
	tests := []struct {
		name              string
		packSizes         []int
		expectedPackSizes []int
		expectedStatus    int
	}{
		{
			name:              "should return default pack sizes",
			packSizes:         []int{250, 500, 1000, 2000, 5000},
			expectedPackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedStatus:    http.StatusOK,
		},
		{
			name:              "should return sorted pack sizes",
			packSizes:         []int{5000, 250, 1000, 500},
			expectedPackSizes: []int{250, 500, 1000, 5000},
			expectedStatus:    http.StatusOK,
		},
		{
			name:              "should return single pack size",
			packSizes:         []int{100},
			expectedPackSizes: []int{100},
			expectedStatus:    http.StatusOK,
		},
		{
			name:              "should return empty pack sizes",
			packSizes:         []int{},
			expectedPackSizes: []int{},
			expectedStatus:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := domain.NewPackCalculator(tt.packSizes)
			handler := NewPackSizesHandler(calculator)

			req := httptest.NewRequest(http.MethodGet, "/pack-sizes", nil)
			w := httptest.NewRecorder()

			handler.Handle(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)

			assert.Contains(t, response, "pack_sizes")

			packSizesResponse := response["pack_sizes"].([]interface{})
			actualPackSizes := make([]int, len(packSizesResponse))
			for i, v := range packSizesResponse {
				actualPackSizes[i] = int(v.(float64))
			}

			assert.Equal(t, tt.expectedPackSizes, actualPackSizes)
		})
	}
}

func TestPackSizesHandler_HandlePost(t *testing.T) {
	tests := []struct {
		name              string
		initialPackSizes  []int
		requestBody       map[string]interface{}
		expectedStatus    int
		expectedPackSizes []int
		shouldHaveError   bool
		expectedError     string
	}{
		{
			name:              "should update pack sizes successfully",
			initialPackSizes:  []int{250, 500, 1000},
			requestBody:       map[string]interface{}{"pack_sizes": []int{100, 200, 300}},
			expectedStatus:    http.StatusOK,
			expectedPackSizes: []int{100, 200, 300},
			shouldHaveError:   false,
		},
		{
			name:              "should sort updated pack sizes",
			initialPackSizes:  []int{250, 500, 1000},
			requestBody:       map[string]interface{}{"pack_sizes": []int{500, 100, 300}},
			expectedStatus:    http.StatusOK,
			expectedPackSizes: []int{100, 300, 500},
			shouldHaveError:   false,
		},
		{
			name:             "should reject empty pack sizes",
			initialPackSizes: []int{250, 500, 1000},
			requestBody:      map[string]interface{}{"pack_sizes": []int{}},
			expectedStatus:   http.StatusBadRequest,
			shouldHaveError:  true,
			expectedError:    "Pack sizes cannot be empty",
		},
		{
			name:             "should reject negative pack size",
			initialPackSizes: []int{250, 500, 1000},
			requestBody:      map[string]interface{}{"pack_sizes": []int{100, -50, 200}},
			expectedStatus:   http.StatusBadRequest,
			shouldHaveError:  true,
			expectedError:    "All pack sizes must be positive",
		},
		{
			name:             "should reject zero pack size",
			initialPackSizes: []int{250, 500, 1000},
			requestBody:      map[string]interface{}{"pack_sizes": []int{100, 0, 200}},
			expectedStatus:   http.StatusBadRequest,
			shouldHaveError:  true,
			expectedError:    "All pack sizes must be positive",
		},
		{
			name:              "should handle single pack size update",
			initialPackSizes:  []int{250, 500, 1000},
			requestBody:       map[string]interface{}{"pack_sizes": []int{750}},
			expectedStatus:    http.StatusOK,
			expectedPackSizes: []int{750},
			shouldHaveError:   false,
		},
		{
			name:              "should handle large pack sizes",
			initialPackSizes:  []int{250, 500},
			requestBody:       map[string]interface{}{"pack_sizes": []int{10000, 50000, 100000}},
			expectedStatus:    http.StatusOK,
			expectedPackSizes: []int{10000, 50000, 100000},
			shouldHaveError:   false,
		},
		{
			name:              "should handle duplicate pack sizes",
			initialPackSizes:  []int{250, 500},
			requestBody:       map[string]interface{}{"pack_sizes": []int{100, 100, 200}},
			expectedStatus:    http.StatusOK,
			expectedPackSizes: []int{100, 100, 200},
			shouldHaveError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := domain.NewPackCalculator(tt.initialPackSizes)
			handler := NewPackSizesHandler(calculator)

			bodyBytes, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/pack-sizes", bytes.NewBuffer(bodyBytes))
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

				assert.Contains(t, response, "message")
				assert.Contains(t, response, "pack_sizes")
				assert.Equal(t, "Pack sizes updated successfully", response["message"])

				packSizesResponse := response["pack_sizes"].([]interface{})
				actualPackSizes := make([]int, len(packSizesResponse))
				for i, v := range packSizesResponse {
					actualPackSizes[i] = int(v.(float64))
				}

				assert.Equal(t, tt.expectedPackSizes, actualPackSizes)

				// Verify the calculator was actually updated
				assert.Equal(t, tt.expectedPackSizes, calculator.GetPackSizes())
			}
		})
	}
}

func TestPackSizesHandler_HandlePost_InvalidJSON(t *testing.T) {
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
			calculator := domain.NewPackCalculator([]int{250, 500, 1000})
			handler := NewPackSizesHandler(calculator)

			req := httptest.NewRequest(http.MethodPost, "/pack-sizes", bytes.NewBufferString(tt.requestBody))
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

func TestPackSizesHandler_ResponseFormat(t *testing.T) {
	t.Run("GET should return valid JSON format", func(t *testing.T) {
		calculator := domain.NewPackCalculator([]int{250, 500, 1000})
		handler := NewPackSizesHandler(calculator)

		req := httptest.NewRequest(http.MethodGet, "/pack-sizes", nil)
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "pack_sizes")
		assert.Len(t, response, 1)
	})

	t.Run("POST should return valid JSON format with message", func(t *testing.T) {
		calculator := domain.NewPackCalculator([]int{250, 500, 1000})
		handler := NewPackSizesHandler(calculator)

		bodyBytes, _ := json.Marshal(map[string]interface{}{"pack_sizes": []int{100, 200}})
		req := httptest.NewRequest(http.MethodPost, "/pack-sizes", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "message")
		assert.Contains(t, response, "pack_sizes")
		assert.Len(t, response, 2)
	})
}

func TestPackSizesHandler_Concurrency(t *testing.T) {
	t.Run("should handle concurrent GET requests safely", func(t *testing.T) {
		calculator := domain.NewPackCalculator([]int{250, 500, 1000})
		handler := NewPackSizesHandler(calculator)

		done := make(chan bool)

		for i := 0; i < 10; i++ {
			go func() {
				req := httptest.NewRequest(http.MethodGet, "/pack-sizes", nil)
				w := httptest.NewRecorder()
				handler.Handle(w, req)
				assert.Equal(t, http.StatusOK, w.Code)
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("should handle concurrent POST requests safely", func(t *testing.T) {
		calculator := domain.NewPackCalculator([]int{250, 500, 1000})
		handler := NewPackSizesHandler(calculator)

		done := make(chan bool)

		for i := 0; i < 5; i++ {
			go func(idx int) {
				body := map[string]interface{}{"pack_sizes": []int{100 * (idx + 1), 200 * (idx + 1)}}
				bodyBytes, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/pack-sizes", bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				handler.Handle(w, req)
				assert.Equal(t, http.StatusOK, w.Code)
				done <- true
			}(i)
		}

		for i := 0; i < 5; i++ {
			<-done
		}
	})
}