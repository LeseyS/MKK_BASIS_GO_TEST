package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) GetMemberRole(ctx context.Context, teamID, userID int64) (domain.Role, error) {
	var role domain.Role
	const sql = `SELECT role FROM team_members WHERE team_id = ? AND user_id = ?`

	txOrPool := transaction.TryExtractTX(ctx)
	err := txOrPool.QueryRowContext(ctx, sql, teamID, userID).Scan(&role)
	if errors.Is(err, ErrNoRows) {
		return role, apperr.ErrNotFound
	}

	if err != nil {
		return role, fmt.Errorf("getting member role: %w", err)
	}

	return role, nil
}
