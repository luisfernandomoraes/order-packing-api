package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/luisfernandomoraes/order-packing-api/internal/response"
)

func Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
				response.Error(w, http.StatusInternalServerError, "Internal server error")
			}
		}()

		next(w, r)
	}
}
