package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	goredis "github.com/redis/go-redis/v9"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
)

func (r *Redis) GetTaskList(ctx context.Context, f dto.ListTasks) (dto.ListTasksOut, error) {
	raw, err := r.redis.Get(ctx, taskListKey(f)).Bytes()
	if errors.Is(err, goredis.Nil) {
		return dto.ListTasksOut{}, apperr.ErrNotFound
	}
	if err != nil {
		return dto.ListTasksOut{}, fmt.Errorf("redis.Get: %w", err)
	}

	var out dto.ListTasksOut
	if err := json.Unmarshal(raw, &out); err != nil {
		return dto.ListTasksOut{}, fmt.Errorf("unmarshalling task list: %w", err)
	}

	return out, nil
}
