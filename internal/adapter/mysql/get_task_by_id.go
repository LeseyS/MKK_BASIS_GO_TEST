package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) GetTaskByID(ctx context.Context, id int64) (domain.Task, error) {
	const query = `SELECT id, team_id, title, description, status, assignee_id, created_by, created_at, updated_at
		FROM tasks WHERE id = ?`

	txOrPool := transaction.TryExtractTX(ctx)

	var t domain.Task
	err := txOrPool.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.TeamID, &t.Title, &t.Description, &t.Status,
		&t.AssigneeID, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return t, apperr.ErrNotFound
	}

	if err != nil {
		return t, fmt.Errorf("scanning task: %w", err)
	}

	return t, nil
}
