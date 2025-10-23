package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecovery(t *testing.T) {
	t.Run("recovers from string panic", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Recovery(func(w http.ResponseWriter, r *http.Request) {
			panic("something went wrong")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var body map[string]string
		err := json.NewDecoder(rr.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "Internal server error", body["error"])

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "Panic recovered: something went wrong")
		assert.Contains(t, logOutput, "goroutine") // Stack trace should be present
	})

	t.Run("recovers from error panic", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		testErr := errors.New("critical error")
		handler := Recovery(func(w http.ResponseWriter, r *http.Request) {
			panic(testErr)
		})

		req := httptest.NewRequest(http.MethodPost, "/api/calculate", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var body map[string]string
		err := json.NewDecoder(rr.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "Internal server error", body["error"])

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "Panic recovered: critical error")
	})

	t.Run("recovers from custom type panic", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		type customPanic struct {
			Code    int
			Message string
		}

		handler := Recovery(func(w http.ResponseWriter, r *http.Request) {
			panic(customPanic{Code: 42, Message: "custom panic"})
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "Panic recovered")
		assert.Contains(t, logOutput, "42")
		assert.Contains(t, logOutput, "custom panic")
	})

	t.Run("recovers from nil panic", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Recovery(func(w http.ResponseWriter, r *http.Request) {
			panic(nil)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		// Note: panic(nil) doesn't actually trigger recover, so handler completes normally
		// This test verifies that the middleware handles this edge case gracefully
		handler(rr, req)

		// With panic(nil), recover() returns nil, so no recovery happens
		// This is expected Go behavior
	})

	t.Run("does not interfere when no panic", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Recovery(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		})

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, `{"status":"ok"}`, rr.Body.String())

		logOutput := logBuffer.String()
		assert.NotContains(t, logOutput, "Panic recovered")
	})

	t.Run("includes stack trace in log", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Recovery(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "Panic recovered: test panic")
		// Stack trace should include goroutine info and file paths
		assert.Contains(t, logOutput, "goroutine")
		assert.Contains(t, logOutput, "middleware/recovery")
	})

	t.Run("sets correct content type in error response", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Recovery(func(w http.ResponseWriter, r *http.Request) {
			panic("test")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	})

	t.Run("panic after partial response write", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Recovery(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("partial"))
			panic("late panic")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		// Status was already written as 200, recovery can't change it
		assert.Equal(t, http.StatusOK, rr.Code)

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "Panic recovered: late panic")
	})

	t.Run("panic with integer value", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Recovery(func(w http.ResponseWriter, r *http.Request) {
			panic(42)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "Panic recovered: 42")
	})

	t.Run("preserves request context through recovery", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Recovery(func(w http.ResponseWriter, r *http.Request) {
			// Verify request is still valid when panic occurs
			assert.NotNil(t, r.Context())
			panic("test")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
