package mysql

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) CreateTeamMembers(ctx context.Context, teamID, userID int64, role domain.Role) error {
	const query = `INSERT INTO team_members (team_id, user_id, role) VALUES (?, ?, ?)`

	txOrPool := transaction.TryExtractTX(ctx)
	_, err := txOrPool.ExecContext(ctx, query, teamID, userID, role)
	if err != nil {
		return fmt.Errorf("create team members: %w", err)
	}

	return nil
}
