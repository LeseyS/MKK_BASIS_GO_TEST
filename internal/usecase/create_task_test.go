package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTask_NotAMember(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return "", apperr.ErrNotFound
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.CreateTask(context.Background(), dto.CreateTaskIn{
		TeamID: 1, ActorID: 42, Title: "task",
	})

	require.ErrorIs(t, err, apperr.ErrForbidden)
}

func TestCreateTask_AssigneeOutsideTeam(t *testing.T) {
	const actorID, assigneeID = int64(1), int64(2)

	my := &mockMySQL{
		getMemberRole: func(_ context.Context, _ int64, userID int64) (domain.Role, error) {
			if userID == assigneeID {
				return "", apperr.ErrNotFound
			}
			return domain.RoleMember, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	assignee := assigneeID
	_, err := uc.CreateTask(context.Background(), dto.CreateTaskIn{
		TeamID: 1, ActorID: actorID, Title: "task", AssigneeID: &assignee,
	})

	require.ErrorIs(t, err, apperr.ErrInvalidAssignee)
	assert.NotErrorIs(t, err, apperr.ErrForbidden)
}

func TestCreateTask_InvalidStatus(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.CreateTask(context.Background(), dto.CreateTaskIn{
		TeamID: 1, ActorID: 1, Title: "task", Status: "bogus",
	})

	require.ErrorIs(t, err, apperr.ErrValidation)
}

func TestCreateTask_SuccessInvalidatesCache(t *testing.T) {
	const teamID = int64(7)

	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		createTask: func(_ context.Context, task domain.Task) (domain.Task, error) {
			task.ID = 100
			return task, nil
		},
	}
	rd := &mockRedis{}
	uc := newUseCase(my, rd, nil, nil)

	out, err := uc.CreateTask(context.Background(), dto.CreateTaskIn{
		TeamID: teamID, ActorID: 1, Title: "task",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(100), out.Task.ID)
	assert.Equal(t, domain.StatusTodo, out.Task.Status, "пустой статус должен становиться todo")
	assert.Equal(t, []int64{teamID}, rd.invalidatedTeams)
}

func TestCreateTask_CacheFailureDoesNotFailRequest(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		createTask: func(_ context.Context, task domain.Task) (domain.Task, error) {
			return task, nil
		},
	}
	rd := &mockRedis{
		invalidate: func(context.Context, int64) error {
			return errors.New("redis is down")
		},
	}
	uc := newUseCase(my, rd, nil, nil)

	_, err := uc.CreateTask(context.Background(), dto.CreateTaskIn{
		TeamID: 1, ActorID: 1, Title: "task",
	})

	require.NoError(t, err, "недоступный кеш не должен ронять создание задачи")
}
