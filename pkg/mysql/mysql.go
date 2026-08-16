package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

type Config struct {
	User     string `envconfig:"MYSQL_USER"     required:"true"`
	Password string `envconfig:"MYSQL_PASSWORD" required:"true"`
	Port     string `envconfig:"MYSQL_PORT"     required:"true"`
	Host     string `envconfig:"MYSQL_HOST"     required:"true"`
	DBName   string `envconfig:"MYSQL_DB_NAME"  required:"true"`
}

type Pool struct {
	*sql.DB
}

func New(ctx context.Context, c Config) (*Pool, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	return &Pool{DB: db}, nil
}

func (p *Pool) Close() {
	if err := p.DB.Close(); err != nil {
		log.Error().Err(err).Msg("mysql: close error")
		return
	}

	log.Info().Msg("mysql: closed")
}
