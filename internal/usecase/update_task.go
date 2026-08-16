package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
	"github.com/rs/zerolog/log"
)

func (uc *UseCase) UpdateTask(ctx context.Context, in dto.UpdateTaskIn) (out dto.UpdateTaskOut, err error) {
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

	if in.AssigneeID != nil {
		if err := uc.requireMembership(ctx, task.TeamID, *in.AssigneeID); err != nil {
			if errors.Is(err, apperr.ErrForbidden) {
				return out, apperr.ErrInvalidAssignee
			}
			return out, err
		}
	}

	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		if err := uc.mysql.UpdateTask(ctx, in); err != nil {
			if errors.Is(err, apperr.ErrNotFound) {
				return apperr.ErrNotFound
			}
			return fmt.Errorf("uc.mysql.UpdateTask: %w", err)
		}

		updated, err := uc.mysql.GetTaskByID(ctx, in.TaskID)
		if err != nil {
			return fmt.Errorf("uc.mysql.GetTaskByID (post-update): %w", err)
		}

		if err := uc.redis.InvalidateTeamTasks(ctx, updated.TeamID); err != nil {
			log.Error().Err(err).Msg("usecase: redis.InvalidateTeamTasks")
		}

		out.Task = updated
		return nil
	})

	return out, err
}
