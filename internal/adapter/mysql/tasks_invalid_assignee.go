package mysql

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) TasksInvalidAssignee(ctx context.Context, userID int64) ([]dto.TaskInvalidAssignee, error) {
	const query = `
		SELECT tsk.id, tsk.team_id, t.name, tsk.title, tsk.status, tsk.assignee_id, u.email
		FROM tasks tsk
		JOIN teams t ON t.id = tsk.team_id
		JOIN users u ON u.id = tsk.assignee_id
		WHERE tsk.assignee_id IS NOT NULL
		  AND tsk.team_id IN (SELECT team_id FROM team_members WHERE user_id = ?)
		  AND NOT EXISTS (
			  SELECT 1
			  FROM team_members tm
			  WHERE tm.team_id = tsk.team_id AND tm.user_id = tsk.assignee_id
		  )
		ORDER BY tsk.team_id, tsk.id`

	txOrPool := transaction.TryExtractTX(ctx)

	rows, err := txOrPool.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("querying tasks with invalid assignee: %w", err)
	}
	defer rows.Close()

	tasks := make([]dto.TaskInvalidAssignee, 0)
	for rows.Next() {
		var t dto.TaskInvalidAssignee
		if err := rows.Scan(&t.TaskID, &t.TeamID, &t.TeamName, &t.Title,
			&t.Status, &t.AssigneeID, &t.AssigneeEmail); err != nil {
			return nil, fmt.Errorf("scanning invalid assignee row: %w", err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating invalid assignee rows: %w", err)
	}

	return tasks, nil
}
