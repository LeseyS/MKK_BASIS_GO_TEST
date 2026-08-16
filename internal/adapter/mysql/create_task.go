package mysql

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) CreateTask(ctx context.Context, t domain.Task) (domain.Task, error) {
	const query = `INSERT INTO tasks (team_id, title, description, status, assignee_id, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`

	txOrPool := transaction.TryExtractTX(ctx)

	res, err := txOrPool.ExecContext(ctx, query,
		t.TeamID, t.Title, t.Description, t.Status, t.AssigneeID, t.CreatedBy)
	if err != nil {
		return domain.Task{}, fmt.Errorf("exec insert: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.Task{}, fmt.Errorf("res.LastInsertId: %w", err)
	}

	task, err := my.GetTaskByID(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("my.GetTaskByID: %w", err)
	}

	return task, nil
}
