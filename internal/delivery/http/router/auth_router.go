package router

import (
	"github.com/devlucas-java/luca-s3/internal/delivery/http/adapter"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/handler"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/middleware"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/repository"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/security/jwt"
	"github.com/go-chi/chi"
)

type AuthRouter struct {
	authHandler    *handler.AuthHandler
	jwtService     *jwt.JWTService
	userRepository repository.UserRepository
}

func NewAuthRouter(handler *handler.AuthHandler, jwt *jwt.JWTService, repo repository.UserRepository) *AuthRouter {
	return &AuthRouter{
		authHandler:    handler,
		jwtService:     jwt,
		userRepository: repo,
	}
}

func (r *AuthRouter) RegisterRouters(c chi.Router) {

	c.Post("/login", adapter.Adapt(r.authHandler.Login))
	c.Post("/register", adapter.Adapt(r.authHandler.Register))

	c.Route("/", func(protect chi.Router) {
		protect.Use(middleware.AuthMiddleware(r.jwtService, r.userRepository))

		protect.Put("/password", adapter.Adapt(r.authHandler.ChangePassword))
	})
}
