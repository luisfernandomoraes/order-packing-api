package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	t.Run("logs request with default 200 status when handler writes body only", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Logging(func(w http.ResponseWriter, r *http.Request) {
			// Handler writes body without calling WriteHeader
			_, _ = w.Write([]byte("OK"))
		})

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rr := httptest.NewRecorder()

		handler(rr, req)

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "GET")
		assert.Contains(t, logOutput, "/health")
		assert.Contains(t, logOutput, "200") // Default status
		assert.Contains(t, logOutput, "192.168.1.1:12345")
	})

	t.Run("logs request with explicit status code", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Logging(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("Created"))
		})

		req := httptest.NewRequest(http.MethodPost, "/api/calculate", nil)
		req.RemoteAddr = "10.0.0.1:54321"
		rr := httptest.NewRecorder()

		handler(rr, req)

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "POST")
		assert.Contains(t, logOutput, "/api/calculate")
		assert.Contains(t, logOutput, "201")
		assert.Contains(t, logOutput, "10.0.0.1:54321")
	})

	t.Run("logs error status codes", func(t *testing.T) {
		testCases := []struct {
			name       string
			statusCode int
			expected   string
		}{
			{"bad request", http.StatusBadRequest, "400"},
			{"not found", http.StatusNotFound, "404"},
			{"internal server error", http.StatusInternalServerError, "500"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var logBuffer bytes.Buffer
				log.SetOutput(&logBuffer)
				t.Cleanup(func() {
					log.SetOutput(os.Stderr)
				})

				handler := Logging(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.statusCode)
				})

				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rr := httptest.NewRecorder()

				handler(rr, req)

				logOutput := logBuffer.String()
				assert.Contains(t, logOutput, tc.expected)
			})
		}
	})

	t.Run("logs duration in reasonable format", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Logging(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		logOutput := logBuffer.String()
		// Duration should be logged with time unit (ns, µs, ms, s)
		assert.Regexp(t, `\d+(\.\d+)?(ns|µs|ms|s)`, logOutput)
	})

	t.Run("captures status from WriteHeader", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Logging(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "404")
	})

	t.Run("logs full request URI including query params", func(t *testing.T) {
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		t.Cleanup(func() {
			log.SetOutput(os.Stderr)
		})

		handler := Logging(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/api/calculate?order=250&pack=1000", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "/api/calculate?order=250&pack=1000")
	})

	t.Run("responseWriter properly delegates Write", func(t *testing.T) {
		handler := Logging(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("test response"))
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, "test response", rr.Body.String())
	})

	t.Run("responseWriter preserves headers", func(t *testing.T) {
		handler := Logging(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Custom", "value")
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, "value", rr.Header().Get("X-Custom"))
	})

	t.Run("logs all HTTP methods", func(t *testing.T) {
		methods := []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodOptions,
		}

		for _, method := range methods {
			t.Run(method, func(t *testing.T) {
				var logBuffer bytes.Buffer
				log.SetOutput(&logBuffer)
				t.Cleanup(func() {
					log.SetOutput(os.Stderr)
				})

				handler := Logging(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})

				req := httptest.NewRequest(method, "/", nil)
				rr := httptest.NewRecorder()

				handler(rr, req)

				logOutput := logBuffer.String()
				assert.Contains(t, logOutput, method)
			})
		}
	})
}
