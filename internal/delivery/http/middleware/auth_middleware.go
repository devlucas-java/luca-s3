package middleware

import (
	"context"
	"net/http"

	"github.com/devlucas-java/luca-s3/internal/delivery/http/response"
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/internal/domain/errors"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/repository"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/luca-s3/pkg/logger"
	"github.com/go-chi/jwtauth"
	"github.com/gocql/gocql"
)

type contextKey string

const AuthKey contextKey = "auth_context"

var log = logger.Instance()

func AuthMiddleware(jwtService *jwt.JWTService, userRepository repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := jwtauth.TokenFromHeader(r)

			claims, err := jwtService.Validate(token)
			if err != nil {
				response.ResponseError(w, errors.ErrUnauthorized(err))
				return
			}

			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				response.ResponseError(w, errors.ErrUnauthorized(nil))
				return
			}

			userID, err := gocql.ParseUUID(userIDStr)
			if err != nil {
				response.ResponseError(w, errors.ErrInvalidUUID(err))
				return
			}

			email, _ := claims["email"].(string)

			user, err := userRepository.FindByID(userID)
			if err != nil || user == nil {
				log.Errorf("error finding user by ID %s: %v", userID, err)
				response.ResponseError(w, errors.ErrForbidden(err))
				return
			}

			var roles []string
			if rolesClaim, ok := claims["roles"].([]any); ok {
				for _, ro := range rolesClaim {
					if roleStr, ok := ro.(string); ok {
						roles = append(roles, roleStr)
					}
				}
			}

			auth := &entity.User{
				ID:    userID,
				Email: email,
				Roles: roles,
			}

			ctx := context.WithValue(r.Context(), AuthKey, auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
