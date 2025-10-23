package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHealthHandler(t *testing.T) {
	handler := NewHealthHandler()
	assert.NotNil(t, handler)
}

func TestHealthHandler_Handle(t *testing.T) {
	tests := []struct {
		name              string
		method            string
		expectedStatus    int
		expectedStatusMsg string
		expectedAppName   string
		shouldHaveError   bool
	}{
		{
			name:              "should return healthy status with GET",
			method:            http.MethodGet,
			expectedStatus:    http.StatusOK,
			expectedStatusMsg: "healthy",
			expectedAppName:   "Order Packing Calculator API",
			shouldHaveError:   false,
		},
		{
			name:              "should return healthy status with POST",
			method:            http.MethodPost,
			expectedStatus:    http.StatusOK,
			expectedStatusMsg: "healthy",
			expectedAppName:   "Order Packing Calculator API",
			shouldHaveError:   false,
		},
		{
			name:              "should return healthy status with PUT",
			method:            http.MethodPut,
			expectedStatus:    http.StatusOK,
			expectedStatusMsg: "healthy",
			expectedAppName:   "Order Packing Calculator API",
			shouldHaveError:   false,
		},
		{
			name:              "should return healthy status with DELETE",
			method:            http.MethodDelete,
			expectedStatus:    http.StatusOK,
			expectedStatusMsg: "healthy",
			expectedAppName:   "Order Packing Calculator API",
			shouldHaveError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHealthHandler()
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			handler.Handle(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if !tt.shouldHaveError {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedStatusMsg, response["status"])
				assert.Equal(t, tt.expectedAppName, response["app"])
			}
		})
	}
}

func TestHealthHandler_ResponseFormat(t *testing.T) {
	t.Run("should return valid JSON format", func(t *testing.T) {
		handler := NewHealthHandler()
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		var response map[string]string
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "status")
		assert.Contains(t, response, "app")
		assert.Len(t, response, 2)
	})
}
