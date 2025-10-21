package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/luisfernandomoraes/order-packing-api/internal/config"
	"github.com/luisfernandomoraes/order-packing-api/internal/middleware"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

func setupIntegrationServer(t *testing.T) (*httptest.Server, httpClient) {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")

	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	originalWD, err := os.Getwd()
	require.NoError(t, err, "getwd failed")

	require.NoError(t, os.Chdir(projectRoot))
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	cfg := config.Config{
		Port:             "0",
		DefaultPackSizes: []int{250, 500, 1000},
		ReadTimeout:      time.Second,
		WriteTimeout:     time.Second,
		IdleTimeout:      time.Second,
	}

	srv := New(cfg)
	ts := httptest.NewServer(srv.setupRoutes())
	t.Cleanup(ts.Close)

	return ts, ts.Client()
}

func TestIntegration_Endpoints(t *testing.T) {
	ts, client := setupIntegrationServer(t)

	t.Run("health endpoint returns status", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

		assert.Equal(t, "healthy", body["status"])
		assert.Equal(t, "Order Packing Calculator API", body["app"])
	})

	t.Run("calculate POST handles valid payload", func(t *testing.T) {
		payload := map[string]interface{}{
			"order":      250,
			"pack_sizes": []int{250, 500, 1000},
		}
		resp := doJSONRequest(t, client, http.MethodPost, ts.URL+"/api/calculate", payload)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Order      int            `json:"order"`
			TotalItems int            `json:"total_items"`
			Packs      map[string]int `json:"packs"`
			TotalPacks int            `json:"total_packs"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

		assert.Equal(t, 250, body.Order)
		assert.Equal(t, 250, body.TotalItems)
		assert.Equal(t, 1, body.TotalPacks)
		assert.Equal(t, map[string]int{"250": 1}, body.Packs)
	})

	t.Run("calculate POST rejects negative order", func(t *testing.T) {
		payload := map[string]interface{}{
			"order":      -5,
			"pack_sizes": []int{250, 500, 1000},
		}
		resp := doJSONRequest(t, client, http.MethodPost, ts.URL+"/api/calculate", payload)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("calculate POST rejects empty pack sizes", func(t *testing.T) {
		payload := map[string]interface{}{
			"order":      250,
			"pack_sizes": []int{},
		}
		resp := doJSONRequest(t, client, http.MethodPost, ts.URL+"/api/calculate", payload)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("calculate POST rejects negative pack sizes", func(t *testing.T) {
		payload := map[string]interface{}{
			"order":      250,
			"pack_sizes": []int{250, -500, 1000},
		}
		resp := doJSONRequest(t, client, http.MethodPost, ts.URL+"/api/calculate", payload)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("calculate POST with malformed JSON", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/calculate", bytes.NewBufferString("{invalid"))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("calculate method not allowed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/calculate", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("static UI served", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), "<!DOCTYPE html>")
	})

	t.Run("CORS preflight returns immediately", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodOptions, ts.URL+"/api/calculate", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	})
}

func TestIntegration_RecoveryMiddleware(t *testing.T) {
	handler := middleware.Chain(
		func(http.ResponseWriter, *http.Request) {
			panic("boom")
		},
		middleware.CORS,
		middleware.Logging,
		middleware.Recovery,
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "Internal server error")
}

func doJSONRequest(t *testing.T, client httpClient, method, url string, payload interface{}) *http.Response {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}
