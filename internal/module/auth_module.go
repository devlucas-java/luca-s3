package module

import (
	"github.com/devlucas-java/luca-s3/internal/application/service"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/handler"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/router"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/database"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/security/jwt"
	"github.com/go-chi/chi"
	"github.com/gocql/gocql"
)

func InitModuleAuth(cql *gocql.Session, jwt *jwt.JWTService) chi.Router {

	db := database.NewUserDB(cql)
	userMapper := mapper.NewUserMapper()
	service := service.NewAuthService(db, jwt, userMapper)
	handler := handler.NewAuthHandler(service)

	router := router.NewAuthRouter(handler, jwt, db)

	r := chi.NewRouter()

	router.RegisterRouters(r)

	return r
}
