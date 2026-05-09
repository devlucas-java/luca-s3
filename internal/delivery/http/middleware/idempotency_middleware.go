package middleware

import "net/http"

func IdempotencyMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// logic to cache with redis

			next.ServeHTTP(w, r)
		})
	}
}
