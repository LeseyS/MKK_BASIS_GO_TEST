package redis

import (
	"context"
	"fmt"
)

func (r *Redis) InvalidateTeamTasks(ctx context.Context, teamID int64) error {
	pattern := fmt.Sprintf("tasks:team:%d:*", teamID)

	var cursor uint64
	for {
		keys, next, err := r.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("redis.Scan: %w", err)
		}

		if len(keys) > 0 {
			if err := r.redis.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("redis.Del: %w", err)
			}
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	return nil
}
