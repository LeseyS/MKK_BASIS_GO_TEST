package mysql

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) AddMemberTeam(ctx context.Context, in dto.AddMemberTeam) error {
	const query = `INSERT INTO team_members (team_id, user_id, role) VALUES (?, ?, ?)`

	txOrPool := transaction.TryExtractTX(ctx)
	_, err := txOrPool.ExecContext(ctx, query, in.TeamID, in.UserID, in.Role)
	if err != nil {
		if isDuplicateErr(err) {
			return apperr.ErrDuplicate
		}
		return fmt.Errorf("adding member: %w", err)
	}

	return nil
}
