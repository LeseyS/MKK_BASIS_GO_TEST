//go:build integration

package test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"

	"database/sql"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/config"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/app"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/httpserver"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/jwtutil"
	mysqlpkg "github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/mysql"
	redispkg "github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/redis"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/task_client"
)

// Run: make integration-test

const (
	dbName     = "taskdb"
	dbUser     = "taskuser"
	dbPassword = "taskpass"
	httpPort   = "18080"
)

var ctx = context.Background()

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	api *task_client.Client
	db  *sql.DB

	mysqlContainer *tcmysql.MySQLContainer
	redisContainer *tcredis.RedisContainer
}

func (s *Suite) SetupSuite() {
	s.Assertions = s.Require()

	log.Logger = zerolog.Nop()

	mysqlCfg := s.startMySQL()
	redisCfg := s.startRedis()

	s.applyMigrations()

	c := config.Config{
		HTTP:  httpserver.Config{Port: httpPort, RateLimitPerMinute: 100_000},
		MYSQL: mysqlCfg,
		Redis: redisCfg,
		JWT:   jwtutil.JWTConfig{Secret: "integration-secret", TTL: time.Hour},
	}

	go func() {
		if err := app.Run(ctx, c); err != nil {
			panic(fmt.Sprintf("app.Run: %v", err))
		}
	}()

	s.api = task_client.New(task_client.Config{Host: "localhost", Port: httpPort})
	s.NoError(s.api.WaitReady(ctx, 30*time.Second))
}

func (s *Suite) startMySQL() mysqlpkg.Config {
	container, err := tcmysql.Run(ctx, "mysql:8.0",
		tcmysql.WithDatabase(dbName),
		tcmysql.WithUsername(dbUser),
		tcmysql.WithPassword(dbPassword),
	)
	s.NoError(err)
	s.mysqlContainer = container

	host, err := container.Host(ctx)
	s.NoError(err)

	port, err := container.MappedPort(ctx, "3306/tcp")
	s.NoError(err)

	return mysqlpkg.Config{
		User:     dbUser,
		Password: dbPassword,
		Host:     host,
		Port:     port.Port(),
		DBName:   dbName,
	}
}

func (s *Suite) startRedis() redispkg.Config {
	container, err := tcredis.Run(ctx, "redis:alpine")
	s.NoError(err)
	s.redisContainer = container

	host, err := container.Host(ctx)
	s.NoError(err)

	port, err := container.MappedPort(ctx, "6379/tcp")
	s.NoError(err)

	return redispkg.Config{Addr: fmt.Sprintf("%s:%s", host, port.Port())}
}

func (s *Suite) applyMigrations() {
	host, err := s.mysqlContainer.Host(ctx)
	s.NoError(err)

	port, err := s.mysqlContainer.MappedPort(ctx, "3306/tcp")
	s.NoError(err)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		dbUser, dbPassword, host, port.Port(), dbName)

	db, err := sql.Open("mysql", dsn)
	s.NoError(err)
	s.db = db

	driver, err := migratemysql.WithInstance(db, &migratemysql.Config{})
	s.NoError(err)

	path, err := filepath.Abs("../migration/mysql")
	s.NoError(err)

	m, err := migrate.NewWithDatabaseInstance("file://"+path, "mysql", driver)
	s.NoError(err)

	s.NoError(m.Up())
}

func (s *Suite) TearDownSuite() {
	if s.mysqlContainer != nil {
		s.NoError(testcontainers.TerminateContainer(s.mysqlContainer))
	}
	if s.redisContainer != nil {
		s.NoError(testcontainers.TerminateContainer(s.redisContainer))
	}
}

func (s *Suite) SetupTest() {}

func (s *Suite) TearDownTest() {}
