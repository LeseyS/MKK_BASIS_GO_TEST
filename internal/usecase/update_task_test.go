package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/stretchr/testify/require"
)

func TestUpdateTask_TaskNotFound(t *testing.T) {
	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{}, apperr.ErrNotFound
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.UpdateTask(context.Background(), dto.UpdateTaskIn{TaskID: 1, ActorID: 1})

	require.ErrorIs(t, err, apperr.ErrNotFound)
}

func TestUpdateTask_NotAMember(t *testing.T) {
	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{ID: 1, TeamID: 7}, nil
		},
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return "", apperr.ErrNotFound
		},
		updateTask: func(context.Context, dto.UpdateTaskIn) error {
			t.Fatal("до записи дойти не должно")
			return nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.UpdateTask(context.Background(), dto.UpdateTaskIn{TaskID: 1, ActorID: 42})

	require.ErrorIs(t, err, apperr.ErrForbidden)
}

func TestUpdateTask_AssigneeOutsideTeam(t *testing.T) {
	const actorID, assigneeID = int64(1), int64(2)

	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{ID: 1, TeamID: 7}, nil
		},
		getMemberRole: func(_ context.Context, _ int64, userID int64) (domain.Role, error) {
			if userID == assigneeID {
				return "", apperr.ErrNotFound
			}
			return domain.RoleMember, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	assignee := assigneeID
	_, err := uc.UpdateTask(context.Background(), dto.UpdateTaskIn{
		TaskID: 1, ActorID: actorID, AssigneeID: &assignee,
	})

	require.ErrorIs(t, err, apperr.ErrInvalidAssignee)
}

func TestUpdateTask_SuccessInvalidatesCache(t *testing.T) {
	const teamID = int64(7)

	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{ID: 1, TeamID: teamID, Status: domain.StatusDone}, nil
		},
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		updateTask: func(context.Context, dto.UpdateTaskIn) error {
			return nil
		},
	}
	rd := &mockRedis{}
	uc := newUseCase(my, rd, nil, nil)

	status := "done"
	out, err := uc.UpdateTask(context.Background(), dto.UpdateTaskIn{
		TaskID: 1, ActorID: 1, Status: &status,
	})

	require.NoError(t, err)
	require.Equal(t, domain.StatusDone, out.Task.Status)
	require.Equal(t, []int64{teamID}, rd.invalidatedTeams)
	require.Equal(t, 2, my.getTaskByIDCalls, "задача читается до правки и перечитывается после")
}

func TestUpdateTask_FailedWriteKeepsCache(t *testing.T) {
	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{ID: 1, TeamID: 7}, nil
		},
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		updateTask: func(context.Context, dto.UpdateTaskIn) error {
			return errors.New("write failed")
		},
	}
	rd := &mockRedis{}
	uc := newUseCase(my, rd, nil, nil)

	status := "done"
	_, err := uc.UpdateTask(context.Background(), dto.UpdateTaskIn{
		TaskID: 1, ActorID: 1, Status: &status,
	})

	require.Error(t, err)
	require.Empty(t, rd.invalidatedTeams, "при откате правки кеш сбрасывать не нужно")
}
