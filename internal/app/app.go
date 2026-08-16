package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/config"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/adapter/mysql"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/adapter/redis"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/controller/http"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/usecase"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/email_service"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/httpserver"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/jwtutil"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/metrics"
	mysqlPKG "github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/mysql"
	redislib "github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/redis"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/router"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
	"github.com/rs/zerolog/log"
)

func Run(ctx context.Context, c config.Config) error {
	// MySQL
	mySQLPool, err := mysqlPKG.New(ctx, c.MYSQL)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}

	transaction.Init(mySQLPool)

	// Redis
	redisClient, err := redislib.New(c.Redis)
	if err != nil {
		return fmt.Errorf("redislib.New: %w", err)
	}

	// JWT
	jwtIssuer := jwtutil.NewIssuer(c.JWT.Secret, c.JWT.TTL)

	// EmailService
	emailService := email_service.NewEmailService()

	// UseCase
	uc := usecase.New(
		mysql.New(),
		redis.New(redisClient),
		jwtIssuer,
		emailService,
	)

	// Metrics
	httpMetrics := metrics.NewHTTPServer()

	// HTTP
	r := router.New()
	http.Router(r, uc, httpMetrics, jwtIssuer)
	httpServer := httpserver.New(r, c.HTTP)

	log.Info().Msg("app: started")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	log.Info().Msg("app: got signal to stop")

	// Controllers close
	httpServer.Close()

	// Adapters close
	redisClient.Close()
	mySQLPool.Close()

	log.Info().Msg("app: stopped")

	return nil
}
