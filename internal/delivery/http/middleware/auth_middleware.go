package middleware

import (
	"context"
	"net/http"

	"github.com/devlucas-java/luca-s3/internal/domain/entity"
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
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}
			userID, err := gocql.ParseUUID(userIDStr)
			if err != nil {
				http.Error(w, "Invalid duser ID in token", http.StatusUnauthorized)
				return
			}

			email, _ := claims["email"].(string)

			user, err := userRepository.FindByID(userID)
			if err != nil || user == nil {
				log.Errorf("Error finding duser by ID %s: %v", userID, err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			var roles []string
			if rolesClaim, ok := claims["roles"].([]interface{}); ok {
				for _, r := range rolesClaim {
					if roleStr, ok := r.(string); ok {
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
