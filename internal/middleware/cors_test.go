package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	t.Run("sets all CORS headers", func(t *testing.T) {
		handler := CORS(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type", rr.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("OPTIONS request short-circuits with 200", func(t *testing.T) {
		nextCalled := false
		handler := CORS(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusCreated)
		})

		req := httptest.NewRequest(http.MethodOptions, "/api/calculate", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, nextCalled, "OPTIONS should not call next handler")
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("GET request calls next handler", func(t *testing.T) {
		nextCalled := false
		handler := CORS(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/api/calculate", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.True(t, nextCalled, "GET should call next handler")
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("POST request calls next handler", func(t *testing.T) {
		nextCalled := false
		handler := CORS(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusCreated)
		})

		req := httptest.NewRequest(http.MethodPost, "/api/calculate", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.True(t, nextCalled, "POST should call next handler")
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("preserves additional headers set by handler", func(t *testing.T) {
		handler := CORS(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Custom-Header", "value")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, "value", rr.Header().Get("X-Custom-Header"))
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("CORS headers set by middleware appear in response", func(t *testing.T) {
		handler := CORS(func(w http.ResponseWriter, r *http.Request) {
			// Handler sets custom headers which should be preserved
			w.Header().Set("X-Custom-Header", "custom-value")
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		// CORS headers should be present
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		// Custom header should also be present
		assert.Equal(t, "custom-value", rr.Header().Get("X-Custom-Header"))
	})
}
