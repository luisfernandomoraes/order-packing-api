package middleware

import "net/http"

// Chain composes multiple HTTP middlewares into a single handler.
//
// Middlewares are applied from left to right (first middleware in the list
// is the outermost wrapper). This means the first middleware will be the
// first to receive the request and the last to send the response.
//
// Usage:
//
//	mux.HandleFunc("/api/endpoint", middleware.Chain(
//	    handlers.MyHandler,
//	    middleware.CORS,
//	    middleware.Logging,
//	    middleware.Recovery,
//	))
//
// The execution flow with the example above would be:
//
//	Request → CORS → Logging → Recovery → MyHandler → Recovery → Logging → CORS → Response
func Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
