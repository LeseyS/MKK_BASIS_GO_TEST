package mysql

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) TeamListForUser(ctx context.Context, userID int64) ([]domain.Team, error) {
	const sql = `
		SELECT t.id, t.name, t.created_by, t.created_at
		FROM teams t
		JOIN team_members tm ON tm.team_id = t.id
		WHERE tm.user_id = ?
		ORDER BY t.created_at DESC`

	args := []any{
		userID,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	var teams []domain.Team

	rows, err := txOrPool.QueryContext(ctx, sql, args...)
	if err != nil {
		return teams, fmt.Errorf("txOrPool.QueryContext: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedBy, &t.CreatedAt); err != nil {
			return teams, fmt.Errorf("rows.Scan: %w", err)
		}
		teams = append(teams, t)
	}

	return teams, rows.Err()
}
