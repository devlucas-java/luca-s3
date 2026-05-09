package main

import (
	"net/http"

	"github.com/devlucas-java/luca-s3/configs"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/cassandra"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/luca-s3/internal/module"
	"github.com/devlucas-java/luca-s3/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	logger.SetLogLevel(logger.DEBUG)
	log := logger.Instance()

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf(" FATAL in load config: %v", err)
	}

	if err := cassandra.InitSession(cfg); err != nil {
		log.Fatalf("Failed to initialize Cassandra session: %v", err)
	}
	defer cassandra.CloseSession()

	log.Debug("Server started successfully")

	jwt := jwt.NewJWTService(cfg.JwtSecret)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authRouters := module.InitModuleAuth(cassandra.GetSession(), jwt)
	r.Mount("/auth", authRouters)

	http.ListenAndServe(cfg.ServerPort, r)
}
