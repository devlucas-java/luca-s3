package cassandra

import (
	"fmt"
	"time"

	"github.com/devlucas-java/luca-s3/configs"
	"github.com/gocql/gocql"
)

func Bootstrap(cfg *configs.Config) error {

	auth := gocql.PasswordAuthenticator{
		Username: cfg.CassandraUsername,
		Password: cfg.CassandraPassword,
	}

	bootstrapCluster := gocql.NewCluster(cfg.CassandraHosts...)
	bootstrapCluster.Consistency = gocql.Quorum
	bootstrapCluster.Timeout = 15 * time.Second
	bootstrapCluster.ConnectTimeout = 15 * time.Second
	bootstrapCluster.Authenticator = auth

	bootstrapSession, err := bootstrapCluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to create bootstrap Cassandra session: %w", err)
	}
	log.Debug("Connected to Cassandra (bootstrap — no keyspace)")

	if err := runMigrations(bootstrapSession, "./internal/infrastructure/cassandra/migrations"); err != nil {
		bootstrapSession.Close()
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	bootstrapSession.Close()
	return nil
}
