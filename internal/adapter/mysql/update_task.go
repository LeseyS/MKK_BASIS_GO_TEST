package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) UpdateTask(ctx context.Context, in dto.UpdateTaskIn) error {
	txOrPool := transaction.TryExtractTX(ctx)

	const selectSQL = `SELECT title, description, status, assignee_id
			FROM tasks WHERE id = ? FOR UPDATE`

	var (
		curTitle       string
		curDescription string
		curStatus      string
		curAssigneeID  *int64
	)

	err := txOrPool.QueryRowContext(ctx, selectSQL, in.TaskID).
		Scan(&curTitle, &curDescription, &curStatus, &curAssigneeID)
	if errors.Is(err, ErrNoRows) {
		return apperr.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("locking task row: %w", err)
	}

	type change struct {
		field    string
		oldValue *string
		newValue *string
	}
	var changes []change

	setClauses := ""
	args := []any{}

	addStr := func(field string, cur *string, next *string) {
		if next == nil || *next == *cur {
			return
		}
		old := *cur
		changes = append(changes, change{field: field, oldValue: &old, newValue: next})
		if setClauses != "" {
			setClauses += ", "
		}
		setClauses += field + " = ?"
		args = append(args, *next)
		*cur = *next
	}

	addStr("title", &curTitle, in.Title)
	addStr("description", &curDescription, in.Description)
	addStr("status", &curStatus, in.Status)

	if in.AssigneeID != nil && (curAssigneeID == nil || *curAssigneeID != *in.AssigneeID) {
		var old *string
		if curAssigneeID != nil {
			s := fmt.Sprintf("%d", *curAssigneeID)
			old = &s
		}
		newVal := fmt.Sprintf("%d", *in.AssigneeID)
		changes = append(changes, change{field: "assignee_id", oldValue: old, newValue: &newVal})

		if setClauses != "" {
			setClauses += ", "
		}
		setClauses += "assignee_id = ?"
		args = append(args, *in.AssigneeID)
	}

	if len(changes) == 0 {
		return nil
	}

	updateSQL := "UPDATE tasks SET " + setClauses + ", updated_at = NOW() WHERE id = ?"
	args = append(args, in.TaskID)

	if _, err := txOrPool.ExecContext(ctx, updateSQL, args...); err != nil {
		return fmt.Errorf("updating task: %w", err)
	}

	const historySQL = `INSERT INTO task_history (task_id, changed_by, field, old_value, new_value)
			VALUES (?, ?, ?, ?, ?)`

	for _, c := range changes {
		if _, err := txOrPool.ExecContext(ctx, historySQL,
			in.TaskID, in.ActorID, c.field, c.oldValue, c.newValue); err != nil {
			return fmt.Errorf("inserting task_history row (%s): %w", c.field, err)
		}
	}

	return nil
}
