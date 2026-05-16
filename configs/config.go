package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	MinIOEndpoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinIOAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinIOSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinIOUseSSL    bool   `mapstructure:"MINIO_USE_SSL"`

	ServerPort string `mapstructure:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".default.env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("REDIS_ADDR", "redis:6379")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("SERVER_PORT", "50051")
	v.SetDefault("FFMPEG_WORK_DIR", os.TempDir())

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return &Config{
		RedisAddr:     v.GetString("REDIS_ADDR"),
		RedisPassword: v.GetString("REDIS_PASSWORD"),
		RedisDB:       v.GetInt("REDIS_DB"),

		MinIOEndpoint:  v.GetString("MINIO_ENDPOINT"),
		MinIOAccessKey: v.GetString("MINIO_ACCESS_KEY"),
		MinIOSecretKey: v.GetString("MINIO_SECRET_KEY"),
		MinIOUseSSL:    v.GetBool("MINIO_USE_SSL"),

		ServerPort: v.GetString("SERVER_PORT"),
	}, nil
}

// Validate checks required fields.
func (c *Config) Validate() error {
	required := map[string]string{
		"REDIS_ADDR":       c.RedisAddr,
		"MINIO_ENDPOINT":   c.MinIOEndpoint,
		"MINIO_ACCESS_KEY": c.MinIOAccessKey,
		"MINIO_SECRET_KEY": c.MinIOSecretKey,
		"SERVER_PORT":      c.ServerPort,
	}
	for k, v := range required {
		if v == "" {
			return fmt.Errorf("%s is required", k)
		}
	}
	return nil
}
