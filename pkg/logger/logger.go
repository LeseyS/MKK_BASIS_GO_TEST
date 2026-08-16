package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Level         string `default:"error" envconfig:"LOGGER_LEVEL"`
	PrettyConsole bool   `default:"false" envconfig:"LOGGER_PRETTY_CONSOLE"`
	LokiURL       string `default:""      envconfig:"LOKI_URL"`
	LokiLabels    string `default:""      envconfig:"LOKI_LABELS"`
}

func Init(c Config) io.Closer {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if level, err := zerolog.ParseLevel(c.Level); err == nil && level != zerolog.NoLevel {
		zerolog.SetGlobalLevel(level)
	}

	var writers []io.Writer

	if c.PrettyConsole {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	} else {
		writers = append(writers, os.Stderr)
	}

	log.Logger = zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Logger()

	log.Info().Msg("logger initialized")

	return io.NopCloser(nil)
}
