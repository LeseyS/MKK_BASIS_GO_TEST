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

func TestListTasks_NotAMember(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return "", apperr.ErrNotFound
		},
		listTasks: func(context.Context, dto.ListTasks) ([]domain.Task, int64, error) {
			t.Fatal("до выборки задач дойти не должно")
			return nil, 0, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.ListTasks(context.Background(), dto.ListTasks{TeamID: 1, UserID: 42})

	require.ErrorIs(t, err, apperr.ErrForbidden)
}

func TestListTasks_LimitNormalization(t *testing.T) {
	cases := []struct {
		name      string
		limit     int
		offset    int
		wantLimit int
		wantOff   int
	}{
		{"без лимита — дефолт", 0, 0, defaultTasksLimit, 0},
		{"отрицательный лимит — дефолт", -5, 0, defaultTasksLimit, 0},
		{"в допустимых пределах", 50, 10, 50, 10},
		{"выше максимума — кламп", 500, 0, maxTasksLimit, 0},
		{"отрицательный offset", 10, -3, 10, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got dto.ListTasks

			my := &mockMySQL{
				getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
					return domain.RoleMember, nil
				},
				listTasks: func(_ context.Context, f dto.ListTasks) ([]domain.Task, int64, error) {
					got = f
					return []domain.Task{}, 0, nil
				},
			}
			rd := &mockRedis{
				getTaskList: func(context.Context, dto.ListTasks) (dto.ListTasksOut, error) {
					return dto.ListTasksOut{}, apperr.ErrNotFound
				},
			}
			uc := newUseCase(my, rd, nil, nil)

			out, err := uc.ListTasks(context.Background(), dto.ListTasks{
				TeamID: 1, UserID: 1, Limit: c.limit, Offset: c.offset,
			})

			require.NoError(t, err)
			require.Equal(t, c.wantLimit, got.Limit, "лимит, ушедший в БД")
			require.Equal(t, c.wantOff, got.Offset, "offset, ушедший в БД")
			require.Equal(t, c.wantLimit, out.Limit, "лимит в ответе")
			require.Equal(t, c.wantOff, out.Offset, "offset в ответе")
		})
	}
}

func TestListTasks_CacheHitSkipsDB(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		listTasks: func(context.Context, dto.ListTasks) ([]domain.Task, int64, error) {
			t.Fatal("при попадании в кеш поход в БД не нужен")
			return nil, 0, nil
		},
	}
	rd := &mockRedis{
		getTaskList: func(context.Context, dto.ListTasks) (dto.ListTasksOut, error) {
			return dto.ListTasksOut{Tasks: []domain.Task{{ID: 1}}, Total: 1}, nil
		},
	}
	uc := newUseCase(my, rd, nil, nil)

	out, err := uc.ListTasks(context.Background(), dto.ListTasks{TeamID: 1, UserID: 1})

	require.NoError(t, err)
	require.Len(t, out.Tasks, 1)
	require.Equal(t, defaultTasksLimit, out.Limit, "limit проставляется и на закешированном ответе")
	require.Empty(t, rd.setFilters, "повторно класть в кеш не нужно")
}

func TestListTasks_CacheMissStoresResult(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		listTasks: func(context.Context, dto.ListTasks) ([]domain.Task, int64, error) {
			return []domain.Task{{ID: 1}, {ID: 2}}, 2, nil
		},
	}
	rd := &mockRedis{
		getTaskList: func(context.Context, dto.ListTasks) (dto.ListTasksOut, error) {
			return dto.ListTasksOut{}, apperr.ErrNotFound
		},
	}
	uc := newUseCase(my, rd, nil, nil)

	out, err := uc.ListTasks(context.Background(), dto.ListTasks{TeamID: 1, UserID: 1})

	require.NoError(t, err)
	require.Equal(t, int64(2), out.Total)
	require.Len(t, rd.setFilters, 1, "результат должен уехать в кеш")
}

func TestListTasks_BrokenCacheFallsBackToDB(t *testing.T) {
	my := &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return domain.RoleMember, nil
		},
		listTasks: func(context.Context, dto.ListTasks) ([]domain.Task, int64, error) {
			return []domain.Task{{ID: 1}}, 1, nil
		},
	}
	rd := &mockRedis{
		getTaskList: func(context.Context, dto.ListTasks) (dto.ListTasksOut, error) {
			return dto.ListTasksOut{}, errors.New("redis is down")
		},
	}
	uc := newUseCase(my, rd, nil, nil)

	out, err := uc.ListTasks(context.Background(), dto.ListTasks{TeamID: 1, UserID: 1})

	require.NoError(t, err, "недоступный кеш не должен ломать выдачу")
	require.Equal(t, int64(1), out.Total)
}
