package config

import (
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/httpserver"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/jwtutil"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/logger"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/mysql"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/redis"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTP   httpserver.Config
	Logger logger.Config
	MYSQL  mysql.Config
	Redis  redis.Config
	JWT    jwtutil.JWTConfig
}

func New() (Config, error) {
	var config Config

	err := godotenv.Load(".env")
	if err != nil {
		return config, fmt.Errorf("godotenv.Load: %w", err)
	}

	err = envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
