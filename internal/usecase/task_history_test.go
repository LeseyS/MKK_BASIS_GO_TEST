package usecase

import (
	"context"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskHistory_TaskNotFound(t *testing.T) {
	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{}, apperr.ErrNotFound
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.TaskHistory(context.Background(), dto.TaskHistoryIn{TaskID: 1, ActorID: 1})

	require.ErrorIs(t, err, apperr.ErrNotFound)
}

func TestTaskHistory_NotAMember(t *testing.T) {
	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{ID: 1, TeamID: 7}, nil
		},
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return "", apperr.ErrNotFound
		},
		taskHistory: func(context.Context, dto.TaskHistoryIn) ([]dto.TaskHistoryEntry, int64, error) {
			t.Fatal("до выборки истории дойти не должно")
			return nil, 0, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.TaskHistory(context.Background(), dto.TaskHistoryIn{TaskID: 1, ActorID: 42})

	require.ErrorIs(t, err, apperr.ErrForbidden)
}

func TestTaskHistory_LimitNormalization(t *testing.T) {
	cases := []struct {
		name      string
		limit     int
		wantLimit int
	}{
		{"без лимита — дефолт", 0, defaultTasksLimit},
		{"выше максимума — кламп", 1000, maxTasksLimit},
		{"в допустимых пределах", 5, 5},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got dto.TaskHistoryIn

			my := &mockMySQL{
				getTaskByID: func(context.Context, int64) (domain.Task, error) {
					return domain.Task{ID: 1, TeamID: 7}, nil
				},
				getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
					return domain.RoleMember, nil
				},
				taskHistory: func(_ context.Context, in dto.TaskHistoryIn) ([]dto.TaskHistoryEntry, int64, error) {
					got = in
					return []dto.TaskHistoryEntry{}, 0, nil
				},
			}
			uc := newUseCase(my, nil, nil, nil)

			out, err := uc.TaskHistory(context.Background(), dto.TaskHistoryIn{
				TaskID: 1, ActorID: 1, Limit: c.limit,
			})

			require.NoError(t, err)
			assert.Equal(t, c.wantLimit, got.Limit)
			assert.Equal(t, c.wantLimit, out.Limit)
		})
	}
}

func TestTaskHistory_ReturnsEntries(t *testing.T) {
	oldValue, newValue := "todo", "done"

	my := &mockMySQL{
		getTaskByID: func(context.Context, int64) (domain.Task, error) {
			return domain.Task{ID: 1, TeamID: 7}, nil
		},
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		taskHistory: func(context.Context, dto.TaskHistoryIn) ([]dto.TaskHistoryEntry, int64, error) {
			return []dto.TaskHistoryEntry{
				{ID: 2, Field: "status", OldValue: &oldValue, NewValue: &newValue},
			}, 1, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	out, err := uc.TaskHistory(context.Background(), dto.TaskHistoryIn{TaskID: 1, ActorID: 1})

	require.NoError(t, err)
	require.Len(t, out.Entries, 1)
	assert.Equal(t, int64(1), out.Total)
	assert.Equal(t, "status", out.Entries[0].Field)
	assert.Equal(t, "done", *out.Entries[0].NewValue)
}
