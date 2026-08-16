package main

import (
	"context"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/config"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/app"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config.New")
	}

	logCloser := logger.Init(c.Logger)
	defer logCloser.Close()

	ctx := context.Background()

	err = app.Run(ctx, c)
	if err != nil {
		log.Error().Err(err).Msg("app.Run")
	}
}
