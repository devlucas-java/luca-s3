package cassandra

import (
	"fmt"
	"time"

	"github.com/devlucas-java/luca-s3/configs"
	"github.com/devlucas-java/luca-s3/pkg/logger"
	"github.com/gocql/gocql"
)

var session *gocql.Session
var log = logger.Instance()

func InitSession(cfg *configs.Config) error {

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("erros in validate config: %w", err)
	}

	auth := gocql.PasswordAuthenticator{
		Username: cfg.CassandraUsername,
		Password: cfg.CassandraPassword,
	}

	if err := Bootstrap(cfg); err != nil {
		return err
	}

	cluster := gocql.NewCluster(cfg.CassandraHosts...)
	cluster.Keyspace = cfg.CassandraKeyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 15 * time.Second
	cluster.ConnectTimeout = 15 * time.Second
	cluster.Authenticator = auth

	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to create Cassandra session: %w", err)
	}
	log.Debugf("Connected to Cassandra (keyspace: %s)", cfg.CassandraKeyspace)

	return nil
}
