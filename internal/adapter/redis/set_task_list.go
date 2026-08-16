package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
)

func (r *Redis) SetTaskList(ctx context.Context, f dto.ListTasks, out dto.ListTasksOut) error {
	raw, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("marshalling task list: %w", err)
	}

	if err := r.redis.Set(ctx, taskListKey(f), raw, taskListTTL).Err(); err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}

func taskListKey(f dto.ListTasks) string {
	status := f.Status
	if status == "" {
		status = "any"
	}

	assignee := "any"
	if f.AssigneeID != 0 {
		assignee = fmt.Sprintf("%d", f.AssigneeID)
	}

	return fmt.Sprintf("tasks:team:%d:status:%s:assignee:%s:limit:%d:offset:%d",
		f.TeamID, status, assignee, f.Limit, f.Offset)
}
