package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) GetTeamByID(ctx context.Context, teamID int64) (domain.Team, error) {
	const sql = `SELECT id, name, created_by, created_at FROM teams WHERE id = ?`

	txOrPool := transaction.TryExtractTX(ctx)
	var t domain.Team

	err := txOrPool.QueryRowContext(ctx, sql, teamID).Scan(&t.ID, &t.Name, &t.CreatedBy, &t.CreatedAt)
	if errors.Is(err, ErrNoRows) {
		return t, apperr.ErrNotFound
	}

	if err != nil {
		return t, fmt.Errorf("getting team: %w", err)
	}

	return t, nil
}
