package main

import (
	"fmt"
	"net/http"

	"github.com/devlucas-java/luca-s3/configs"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/cassandra"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/luca-s3/internal/module"
	"github.com/devlucas-java/luca-s3/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger.SetLogLevel(logger.DEBUG)
	log := logger.Instance()

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cassandra.InitSession(cfg); err != nil {
		log.Fatalf("failed to initialize Cassandra session: %v", err)
	}
	defer cassandra.CloseSession()

	jwtService := jwt.NewJWTService(cfg.JwtSecret)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/auth", module.InitModuleAuth(cassandra.GetSession(), jwtService))
	r.Mount("/users", module.InitModuleUser(cassandra.GetSession(), jwtService))

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Infof("server listening on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
