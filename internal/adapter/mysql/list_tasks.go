package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) ListTasks(ctx context.Context, f dto.ListTasks) ([]domain.Task, int64, error) {
	txOrPool := transaction.TryExtractTX(ctx)

	var where strings.Builder
	where.WriteString("team_id = ?")
	args := []interface{}{f.TeamID}

	if f.Status != "" {
		where.WriteString(" AND status = ?")
		args = append(args, f.Status)
	}
	if f.AssigneeID != 0 {
		where.WriteString(" AND assignee_id = ?")
		args = append(args, f.AssigneeID)
	}

	var total int64
	countQuery := "SELECT COUNT(*) FROM tasks WHERE " + where.String()
	if err := txOrPool.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting tasks: %w", err)
	}

	if total == 0 {
		return []domain.Task{}, 0, nil
	}

	query := `SELECT id, team_id, title, description, status, assignee_id, created_by, created_at, updated_at
		FROM tasks WHERE ` + where.String() + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, f.Limit, f.Offset)

	rows, err := txOrPool.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0, f.Limit)
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.TeamID, &t.Title, &t.Description, &t.Status,
			&t.AssigneeID, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning task row: %w", err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating task rows: %w", err)
	}

	return tasks, total, nil
}
