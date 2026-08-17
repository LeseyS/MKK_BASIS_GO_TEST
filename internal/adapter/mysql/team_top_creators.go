package mysql

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

const topCreatorsPerTeam = 3

func (my *MySQL) TeamTopCreators(ctx context.Context, userID int64) ([]dto.TeamTopCreator, error) {
	const query = `
		WITH creators AS (
			SELECT
				tsk.team_id,
				tsk.created_by AS user_id,
				COUNT(*) AS tasks_created
			FROM tasks tsk
			WHERE tsk.created_at >= NOW() - INTERVAL 1 MONTH
			  AND tsk.team_id IN (SELECT team_id FROM team_members WHERE user_id = ?)
			GROUP BY tsk.team_id, tsk.created_by
		),
		ranked AS (
			SELECT
				c.team_id,
				c.user_id,
				c.tasks_created,
				ROW_NUMBER() OVER (
					PARTITION BY c.team_id
					ORDER BY c.tasks_created DESC, c.user_id ASC
				) AS rank_position
			FROM creators c
		)
		SELECT r.team_id, t.name, r.user_id, u.name, u.email, r.tasks_created, r.rank_position
		FROM ranked r
		JOIN teams t ON t.id = r.team_id
		JOIN users u ON u.id = r.user_id
		WHERE r.rank_position <= ?
		ORDER BY t.name, r.rank_position`

	txOrPool := transaction.TryExtractTX(ctx)

	rows, err := txOrPool.QueryContext(ctx, query, userID, topCreatorsPerTeam)
	if err != nil {
		return nil, fmt.Errorf("querying top creators: %w", err)
	}
	defer rows.Close()

	top := make([]dto.TeamTopCreator, 0)
	for rows.Next() {
		var c dto.TeamTopCreator
		if err := rows.Scan(&c.TeamID, &c.TeamName, &c.UserID, &c.UserName,
			&c.UserEmail, &c.TasksCreated, &c.Position); err != nil {
			return nil, fmt.Errorf("scanning top creator row: %w", err)
		}
		top = append(top, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating top creator rows: %w", err)
	}

	return top, nil
}
