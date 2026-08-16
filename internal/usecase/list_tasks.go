package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/rs/zerolog/log"
)

const (
	defaultTasksLimit = 20
	maxTasksLimit     = 100
)

func (uc *UseCase) ListTasks(ctx context.Context, in dto.ListTasks) (dto.ListTasksOut, error) {
	if _, err := uc.mysql.GetMemberRole(ctx, in.TeamID, in.UserID); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return dto.ListTasksOut{}, apperr.ErrForbidden
		}
		return dto.ListTasksOut{}, fmt.Errorf("uc.mysql.GetMemberRole: %w", err)
	}

	if in.Limit <= 0 {
		in.Limit = defaultTasksLimit
	}
	if in.Limit > maxTasksLimit {
		in.Limit = maxTasksLimit
	}
	if in.Offset < 0 {
		in.Offset = 0
	}

	if cached, err := uc.redis.GetTaskList(ctx, in); err == nil {
		cached.Limit, cached.Offset = in.Limit, in.Offset
		return cached, nil
	} else if !errors.Is(err, apperr.ErrNotFound) {
		log.Error().Err(err).Msg("usecase: redis.GetTaskList")
	}

	tasks, total, err := uc.mysql.ListTasks(ctx, in)
	if err != nil {
		return dto.ListTasksOut{}, fmt.Errorf("uc.mysql.ListTasks: %w", err)
	}

	out := dto.ListTasksOut{Tasks: tasks, Total: total, Limit: in.Limit, Offset: in.Offset}

	if err := uc.redis.SetTaskList(ctx, in, out); err != nil {
		log.Error().Err(err).Msg("usecase: redis.SetTaskList")
	}

	return out, nil
}
