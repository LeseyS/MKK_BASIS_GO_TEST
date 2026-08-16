package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) GetTeam(ctx context.Context, teamID int64) (domain.Team, error) {
	const sql = `SELECT id, name, created_by, created_at FROM teams WHERE id = ?`

	txOrPool := transaction.TryExtractTX(ctx)

	var team domain.Team

	err := txOrPool.QueryRowContext(ctx, sql, teamID).Scan(&team.ID, &team.Name, &team.CreatedBy, &team.CreatedAt)
	if errors.Is(err, ErrNoRows) {
		return team, apperr.ErrNotFound
	}

	if err != nil {
		return team, fmt.Errorf("txOrPool.QueryRowContext: %w", err)
	}

	return team, nil
}
