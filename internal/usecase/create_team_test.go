package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateWithOwner_CreatorBecomesOwner(t *testing.T) {
	const ownerID = int64(42)

	var (
		gotTeamID int64
		gotUserID int64
		gotRole   domain.Role
	)

	my := &mockMySQL{
		createTeam: func(context.Context, dto.CreateTeamIn) (int64, error) {
			return 7, nil
		},
		createTeamMembers: func(_ context.Context, teamID, userID int64, role domain.Role) error {
			gotTeamID, gotUserID, gotRole = teamID, userID, role
			return nil
		},
		getTeamByID: func(_ context.Context, teamID int64) (domain.Team, error) {
			return domain.Team{ID: teamID, Name: "Alpha", CreatedBy: ownerID}, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	out, err := uc.CreateWithOwner(context.Background(), dto.CreateTeamIn{
		Name: "Alpha", OwnerID: ownerID,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(7), out.Team.ID)
	assert.Equal(t, int64(7), gotTeamID)
	assert.Equal(t, ownerID, gotUserID)
	assert.Equal(t, domain.RoleOwner, gotRole, "создатель команды получает роль owner")
}

func TestCreateWithOwner_MembershipFailurePropagates(t *testing.T) {
	my := &mockMySQL{
		createTeam: func(context.Context, dto.CreateTeamIn) (int64, error) {
			return 7, nil
		},
		createTeamMembers: func(context.Context, int64, int64, domain.Role) error {
			return errors.New("insert failed")
		},
		getTeamByID: func(context.Context, int64) (domain.Team, error) {
			t.Fatal("после сбоя вставки участника читать команду не нужно")
			return domain.Team{}, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.CreateWithOwner(context.Background(), dto.CreateTeamIn{
		Name: "Alpha", OwnerID: 1,
	})

	require.Error(t, err)
}

func TestTeamListForUser(t *testing.T) {
	my := &mockMySQL{
		teamListForUser: func(_ context.Context, userID int64) ([]domain.Team, error) {
			assert.Equal(t, int64(42), userID)
			return []domain.Team{{ID: 1, Name: "Alpha"}, {ID: 2, Name: "Beta"}}, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	out, err := uc.TeamListForUser(context.Background(), 42)

	require.NoError(t, err)
	require.Len(t, out, 2)
	assert.Equal(t, "Alpha", out[0].Team.Name)
	assert.Equal(t, "Beta", out[1].Team.Name)
}

func TestTeamListForUser_Empty(t *testing.T) {
	my := &mockMySQL{
		teamListForUser: func(context.Context, int64) ([]domain.Team, error) {
			return nil, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	out, err := uc.TeamListForUser(context.Background(), 1)

	require.NoError(t, err)
	assert.Empty(t, out)
}
