package transaction

import (
	"context"
	"database/sql"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/mysql"
)

//nolint:gochecknoglobals
var (
	pool       *sql.DB
	IsUnitTest bool
)

type ctxKey struct{}

func Init(p *mysql.Pool) {
	pool = p.DB
}

type Transaction struct {
	*sql.Tx
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func TryExtractTX(ctx context.Context) Executor { //nolint:ireturn
	tx, ok := ctx.Value(ctxKey{}).(*Transaction)
	if !ok {
		return pool
	}

	return tx
}
