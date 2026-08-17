package usecase

import (
	"context"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type analyticsMySQL struct {
	MySQL

	teamStats            func(context.Context, int64) ([]dto.TeamStats, error)
	teamTopCreators      func(context.Context, int64) ([]dto.TeamTopCreator, error)
	tasksInvalidAssignee func(context.Context, int64) ([]dto.TaskInvalidAssignee, error)
}

func (m *analyticsMySQL) TeamStats(ctx context.Context, userID int64) ([]dto.TeamStats, error) {
	return m.teamStats(ctx, userID)
}

func (m *analyticsMySQL) TeamTopCreators(ctx context.Context, userID int64) ([]dto.TeamTopCreator, error) {
	return m.teamTopCreators(ctx, userID)
}

func (m *analyticsMySQL) TasksInvalidAssignee(ctx context.Context, userID int64) ([]dto.TaskInvalidAssignee, error) {
	return m.tasksInvalidAssignee(ctx, userID)
}

func TestTeamStats(t *testing.T) {
	my := &analyticsMySQL{
		teamStats: func(_ context.Context, userID int64) ([]dto.TeamStats, error) {
			assert.Equal(t, int64(42), userID, "выборка ограничена командами вызывающего")
			return []dto.TeamStats{{TeamID: 1, Name: "Alpha", MembersCount: 3, DoneLast7Days: 2}}, nil
		},
	}
	uc := New(my, &mockRedis{}, &mockTokenIssuer{}, &mockNotifier{})

	out, err := uc.TeamStats(context.Background(), 42)

	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, int64(2), out[0].DoneLast7Days)
}

func TestTeamStats_DBFailure(t *testing.T) {
	my := &analyticsMySQL{
		teamStats: func(context.Context, int64) ([]dto.TeamStats, error) {
			return nil, errDB
		},
	}
	uc := New(my, &mockRedis{}, &mockTokenIssuer{}, &mockNotifier{})

	_, err := uc.TeamStats(context.Background(), 1)

	require.ErrorIs(t, err, errDB)
}

func TestTeamTopCreators(t *testing.T) {
	my := &analyticsMySQL{
		teamTopCreators: func(_ context.Context, userID int64) ([]dto.TeamTopCreator, error) {
			assert.Equal(t, int64(42), userID)
			return []dto.TeamTopCreator{
				{TeamID: 1, UserID: 1, TasksCreated: 5, Position: 1},
				{TeamID: 1, UserID: 2, TasksCreated: 3, Position: 2},
			}, nil
		},
	}
	uc := New(my, &mockRedis{}, &mockTokenIssuer{}, &mockNotifier{})

	out, err := uc.TeamTopCreators(context.Background(), 42)

	require.NoError(t, err)
	require.Len(t, out, 2)
	assert.Equal(t, 1, out[0].Position)
}

func TestTeamTopCreators_DBFailure(t *testing.T) {
	my := &analyticsMySQL{
		teamTopCreators: func(context.Context, int64) ([]dto.TeamTopCreator, error) {
			return nil, errDB
		},
	}
	uc := New(my, &mockRedis{}, &mockTokenIssuer{}, &mockNotifier{})

	_, err := uc.TeamTopCreators(context.Background(), 1)

	require.ErrorIs(t, err, errDB)
}

func TestTasksInvalidAssignee(t *testing.T) {
	my := &analyticsMySQL{
		tasksInvalidAssignee: func(_ context.Context, userID int64) ([]dto.TaskInvalidAssignee, error) {
			assert.Equal(t, int64(42), userID)
			return []dto.TaskInvalidAssignee{{TaskID: 9, TeamID: 1, AssigneeID: 7}}, nil
		},
	}
	uc := New(my, &mockRedis{}, &mockTokenIssuer{}, &mockNotifier{})

	out, err := uc.TasksInvalidAssignee(context.Background(), 42)

	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, int64(9), out[0].TaskID)
}

func TestTasksInvalidAssignee_DBFailure(t *testing.T) {
	my := &analyticsMySQL{
		tasksInvalidAssignee: func(context.Context, int64) ([]dto.TaskInvalidAssignee, error) {
			return nil, errDB
		},
	}
	uc := New(my, &mockRedis{}, &mockTokenIssuer{}, &mockNotifier{})

	_, err := uc.TasksInvalidAssignee(context.Background(), 1)

	require.ErrorIs(t, err, errDB)
}
