package router

import (
	"github.com/devlucas-java/luca-s3/internal/delivery/http/adapter"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/handler"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/middleware"
	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/repository"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/security/jwt"
	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
	userHandler    *handler.UserHandler
	jwtService     *jwt.JWTService
	userRepository repository.UserRepository
}

func NewUserRouter(
	handler *handler.UserHandler,
	jwt *jwt.JWTService,
	repo repository.UserRepository,
) *UserRouter {
	return &UserRouter{
		userHandler:    handler,
		jwtService:     jwt,
		userRepository: repo,
	}
}

func (r *UserRouter) RegisterRouters(c chi.Router) {
	c.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware(r.jwtService, r.userRepository))

		protected.Get("/me", adapter.Adapt(r.userHandler.GetMe))
		protected.Put("/me", adapter.Adapt(r.userHandler.UpdateMe))
		protected.Delete("/me", adapter.Adapt(r.userHandler.DeleteMe))

		protected.Group(func(admin chi.Router) {
			admin.Use(middleware.RoleMiddleware([]string{enums.RoleAdmin.String()}))

			admin.Get("/{id}", adapter.Adapt(r.userHandler.GetByID))
			admin.Put("/{id}", adapter.Adapt(r.userHandler.UpdateByID))
			admin.Delete("/{id}", adapter.Adapt(r.userHandler.DeleteByID))
		})
	})
}
