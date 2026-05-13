package middleware

import (
	"net/http"
	"slices"

	"github.com/devlucas-java/luca-s3/internal/delivery/http/response"
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/internal/domain/errors"
)

func RoleMiddleware(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth, ok := r.Context().Value(AuthKey).(*entity.User)
			if !ok || auth == nil {
				response.ResponseError(w, errors.ErrUnauthorized(nil))
				return
			}

			allowed := false
			for _, role := range roles {
				if slices.Contains(auth.Roles, role) {
					allowed = true
					break
				}
			}

			if !allowed {
				response.ResponseError(w, errors.ErrForbidden(nil))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
