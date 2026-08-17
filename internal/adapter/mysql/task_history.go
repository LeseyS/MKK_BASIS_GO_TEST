package mysql

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) TaskHistory(ctx context.Context, in dto.TaskHistoryIn) ([]dto.TaskHistoryEntry, int64, error) {
	txOrPool := transaction.TryExtractTX(ctx)

	const countQuery = `SELECT COUNT(*) FROM task_history WHERE task_id = ?`

	var total int64
	if err := txOrPool.QueryRowContext(ctx, countQuery, in.TaskID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting task history: %w", err)
	}

	if total == 0 {
		return []dto.TaskHistoryEntry{}, 0, nil
	}

	const query = `
		SELECT h.id, h.field, h.old_value, h.new_value, h.changed_by, u.name, u.email, h.changed_at
		FROM task_history h
		JOIN users u ON u.id = h.changed_by
		WHERE h.task_id = ?
		ORDER BY h.changed_at DESC, h.id DESC
		LIMIT ? OFFSET ?`

	rows, err := txOrPool.QueryContext(ctx, query, in.TaskID, in.Limit, in.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("querying task history: %w", err)
	}
	defer rows.Close()

	entries := make([]dto.TaskHistoryEntry, 0, in.Limit)
	for rows.Next() {
		var e dto.TaskHistoryEntry
		if err := rows.Scan(&e.ID, &e.Field, &e.OldValue, &e.NewValue,
			&e.ChangedBy, &e.ChangedByName, &e.ChangedByEmail, &e.ChangedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning task history row: %w", err)
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating task history rows: %w", err)
	}

	return entries, total, nil
}
