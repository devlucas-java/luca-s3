package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	CassandraHosts    []string `mapstructure:"CASSANDRA_HOSTS"`
	CassandraKeyspace string   `mapstructure:"CASSANDRA_KEYSPACE"`
	CassandraUsername string   `mapstructure:"CASSANDRA_USERNAME"`
	CassandraPassword string   `mapstructure:"CASSANDRA_PASSWORD"`

	MinIOEndpoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinIOAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinIOSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinIOUseSSL    bool   `mapstructure:"MINIO_USE_SSL"`

	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	ServerPort string `mapstructure:"SERVER_PORT"`
	JwtSecret  string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigFile(".default.env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("REDIS_DB", 0)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	hostsRaw := v.GetString("CASSANDRA_HOSTS")
	var hosts []string
	for _, h := range strings.Split(hostsRaw, ",") {
		h = strings.TrimSpace(h)
		if h != "" {
			hosts = append(hosts, h)
		}
	}

	cfg := &Config{
		CassandraHosts:    hosts,
		CassandraKeyspace: v.GetString("CASSANDRA_KEYSPACE"),
		CassandraUsername: v.GetString("CASSANDRA_USERNAME"),
		CassandraPassword: v.GetString("CASSANDRA_PASSWORD"),

		MinIOEndpoint:  v.GetString("MINIO_ENDPOINT"),
		MinIOAccessKey: v.GetString("MINIO_ACCESS_KEY"),
		MinIOSecretKey: v.GetString("MINIO_SECRET_KEY"),
		MinIOUseSSL:    v.GetBool("MINIO_USE_SSL"),

		RedisAddr:     v.GetString("REDIS_ADDR"),
		RedisPassword: v.GetString("REDIS_PASSWORD"),
		RedisDB:       v.GetInt("REDIS_DB"),

		ServerPort: v.GetString("SERVER_PORT"),
		JwtSecret:  v.GetString("JWT_SECRET"),
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if len(c.CassandraHosts) == 0 {
		return fmt.Errorf("CASSANDRA_HOSTS is required")
	}
	if c.CassandraKeyspace == "" {
		return fmt.Errorf("CASSANDRA_KEYSPACE is required")
	}
	if c.CassandraUsername == "" {
		return fmt.Errorf("CASSANDRA_USERNAME is required")
	}
	if c.CassandraPassword == "" {
		return fmt.Errorf("CASSANDRA_PASSWORD is required")
	}

	if c.MinIOEndpoint == "" {
		return fmt.Errorf("MINIO_ENDPOINT is required")
	}
	if c.MinIOAccessKey == "" {
		return fmt.Errorf("MINIO_ACCESS_KEY is required")
	}
	if c.MinIOSecretKey == "" {
		return fmt.Errorf("MINIO_SECRET_KEY is required")
	}

	if c.RedisAddr == "" {
		return fmt.Errorf("REDIS_ADDR is required")
	}
	if c.RedisDB < 0 {
		return fmt.Errorf("REDIS_DB must be >= 0")
	}

	if c.ServerPort == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}

	return nil
}
