package middleware

import (
	"net/http"

	"github.com/devlucas-java/luca-s3/internal/domain/entity"
)

func RoleMiddleware(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			auth, ok := r.Context().Value(AuthKey).(*entity.User)
			if !ok || auth == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			allowed := false

			for _, r := range roles {
				for _, ur := range auth.Roles {
					if r == ur {
						allowed = true
						break
					}
				}
				if allowed {
					break
				}
			}

			if !allowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
