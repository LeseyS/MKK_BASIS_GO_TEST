package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/rs/zerolog/log"
)

func (uc *UseCase) CreateTask(ctx context.Context, t dto.CreateTaskIn) (dto.CreateTaskOut, error) {
	if err := uc.requireMembership(ctx, t.TeamID, t.ActorID); err != nil {
		return dto.CreateTaskOut{}, err
	}

	if t.AssigneeID != nil {
		if err := uc.requireMembership(ctx, t.TeamID, *t.AssigneeID); err != nil {
			if errors.Is(err, apperr.ErrForbidden) {
				return dto.CreateTaskOut{}, apperr.ErrInvalidAssignee
			}
			return dto.CreateTaskOut{}, err
		}
	}

	newTask, err := domain.NewTask(t.TeamID, t.Title, t.Description,
		domain.TaskStatus(t.Status), t.ActorID, t.AssigneeID)
	if err != nil {
		return dto.CreateTaskOut{}, fmt.Errorf("domain.NewTask: %w", err)
	}

	task, err := uc.mysql.CreateTask(ctx, newTask)
	if err != nil {
		return dto.CreateTaskOut{}, fmt.Errorf("uc.mysql.CreateTask: %w", err)
	}

	if err := uc.redis.InvalidateTeamTasks(ctx, task.TeamID); err != nil {
		log.Error().Err(err).Msg("usecase: redis.InvalidateTeamTasks")
	}

	return dto.CreateTaskOut{Task: task}, nil
}

func (uc *UseCase) requireMembership(ctx context.Context, teamID int64, userID int64) error {
	_, err := uc.mysql.GetMemberRole(ctx, teamID, userID)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return apperr.ErrForbidden
		}
		return fmt.Errorf("uc.mysql.GetMemberRole: %w", err)
	}

	return nil
}
