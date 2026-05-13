package module

import (
	"github.com/devlucas-java/luca-s3/internal/application/service"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/handler"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/router"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/database"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/security/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/gocql/gocql"
)

func InitModuleUser(cql *gocql.Session, jwtService *jwt.JWTService) chi.Router {
	db := database.NewUserDB(cql)
	userMapper := mapper.NewUserMapper()
	svc := service.NewUserService(db, userMapper)
	h := handler.NewUserHandler(svc)
	r := router.NewUserRouter(h, jwtService, db)

	c := chi.NewRouter()
	r.RegisterRouters(c)
	return c
}
