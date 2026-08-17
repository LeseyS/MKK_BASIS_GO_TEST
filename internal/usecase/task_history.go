package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
)

func (uc *UseCase) TaskHistory(ctx context.Context, in dto.TaskHistoryIn) (dto.TaskHistoryOut, error) {
	var out dto.TaskHistoryOut

	task, err := uc.mysql.GetTaskByID(ctx, in.TaskID)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return out, apperr.ErrNotFound
		}
		return out, fmt.Errorf("uc.mysql.GetTaskByID: %w", err)
	}

	if err := uc.requireMembership(ctx, task.TeamID, in.ActorID); err != nil {
		return out, err
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

	entries, total, err := uc.mysql.TaskHistory(ctx, in)
	if err != nil {
		return out, fmt.Errorf("uc.mysql.TaskHistory: %w", err)
	}

	out = dto.TaskHistoryOut{Entries: entries, Total: total, Limit: in.Limit, Offset: in.Offset}

	return out, nil
}
