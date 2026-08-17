package mysql

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) TeamStats(ctx context.Context, userID int64) ([]dto.TeamStats, error) {
	const query = `
		SELECT
			t.id,
			t.name,
			COUNT(DISTINCT tm.user_id) AS members_count,
			COUNT(DISTINCT CASE
				WHEN tsk.status = 'done' AND tsk.updated_at >= NOW() - INTERVAL 7 DAY
				THEN tsk.id
			END) AS done_last_7_days
		FROM teams t
		JOIN team_members tm ON tm.team_id = t.id
		LEFT JOIN tasks tsk ON tsk.team_id = t.id
		WHERE t.id IN (SELECT team_id FROM team_members WHERE user_id = ?)
		GROUP BY t.id, t.name
		ORDER BY done_last_7_days DESC, t.name`

	txOrPool := transaction.TryExtractTX(ctx)

	rows, err := txOrPool.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("querying team stats: %w", err)
	}
	defer rows.Close()

	stats := make([]dto.TeamStats, 0)
	for rows.Next() {
		var s dto.TeamStats
		if err := rows.Scan(&s.TeamID, &s.Name, &s.MembersCount, &s.DoneLast7Days); err != nil {
			return nil, fmt.Errorf("scanning team stats row: %w", err)
		}
		stats = append(stats, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating team stats rows: %w", err)
	}

	return stats, nil
}
