package mysql

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) CreateTeam(ctx context.Context, in dto.CreateTeamIn) (int64, error) {
	const sql = `INSERT INTO teams (name, created_by) VALUES (?, ?)`

	args := []any{
		in.Name,
		in.OwnerID,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	res, err := txOrPool.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("inserting user: %w", err)
	}

	return res.LastInsertId()
}
