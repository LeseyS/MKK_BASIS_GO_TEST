package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/require"
)

func inviteMocks(role domain.Role) *mockMySQL {
	return &mockMySQL{
		getMemberRole: func(context.Context, int64, int64) (domain.Role, error) {
			return role, nil
		},
		getUserByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{ID: 2, Email: "invitee@example.com"}, nil
		},
		addMemberTeam: func(context.Context, dto.AddMemberTeam) error {
			return nil
		},
		getTeamByID: func(context.Context, int64) (domain.Team, error) {
			return domain.Team{ID: 1, Name: "Alpha"}, nil
		},
	}
}

func TestInviteUser_RoleCheck(t *testing.T) {
	cases := []struct {
		role      domain.Role
		wantAllow bool
	}{
		{domain.RoleOwner, true},
		{domain.RoleAdmin, true},
		{domain.RoleMember, false},
	}

	for _, c := range cases {
		t.Run(string(c.role), func(t *testing.T) {
			uc := newUseCase(inviteMocks(c.role), nil, nil, nil)

			_, err := uc.InviteUser(context.Background(), dto.InviteUserIn{
				TeamID: 1, UserID: 1, Email: "invitee@example.com",
			})

			if c.wantAllow {
				require.NoError(t, err)
				return
			}
			require.ErrorIs(t, err, apperr.ErrForbidden)
		})
	}
}

func TestInviteUser_ActorNotInTeam(t *testing.T) {
	my := inviteMocks(domain.RoleOwner)
	my.getMemberRole = func(context.Context, int64, int64) (domain.Role, error) {
		return "", apperr.ErrNotFound
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.InviteUser(context.Background(), dto.InviteUserIn{
		TeamID: 1, UserID: 42, Email: "invitee@example.com",
	})

	require.ErrorIs(t, err, apperr.ErrForbidden)
}

func TestInviteUser_UnknownEmail(t *testing.T) {
	my := inviteMocks(domain.RoleOwner)
	my.getUserByEmail = func(context.Context, string) (domain.User, error) {
		return domain.User{}, apperr.ErrNotFound
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.InviteUser(context.Background(), dto.InviteUserIn{
		TeamID: 1, UserID: 1, Email: "nobody@example.com",
	})

	require.ErrorIs(t, err, apperr.ErrNotFound)
}

func TestInviteUser_AlreadyMember(t *testing.T) {
	my := inviteMocks(domain.RoleOwner)
	my.addMemberTeam = func(context.Context, dto.AddMemberTeam) error {
		return apperr.ErrDuplicate
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.InviteUser(context.Background(), dto.InviteUserIn{
		TeamID: 1, UserID: 1, Email: "invitee@example.com",
	})

	require.ErrorIs(t, err, apperr.ErrAlreadyMember)
}

func TestInviteUser_InviteeAddedAsMember(t *testing.T) {
	var added dto.AddMemberTeam

	my := inviteMocks(domain.RoleOwner)
	my.addMemberTeam = func(_ context.Context, in dto.AddMemberTeam) error {
		added = in
		return nil
	}
	nt := &mockNotifier{}
	uc := newUseCase(my, nil, nil, nt)

	emailSent, err := uc.InviteUser(context.Background(), dto.InviteUserIn{
		TeamID: 1, UserID: 1, Email: "invitee@example.com",
	})

	require.NoError(t, err)
	require.True(t, emailSent)
	require.Equal(t, int64(2), added.UserID, "в команду добавляется приглашённый, а не приглашающий")
	require.Equal(t, domain.RoleMember, added.Role)
	require.Equal(t, []string{"invitee@example.com"}, nt.sentTo)
}

func TestInviteUser_EmailFailureStillSucceeds(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"провайдер вернул ошибку", errors.New("smtp timeout")},
		{"разомкнутый circuit breaker", gobreaker.ErrOpenState},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			nt := &mockNotifier{
				sendInvite: func(context.Context, string, string) error {
					return c.err
				},
			}
			uc := newUseCase(inviteMocks(domain.RoleOwner), nil, nil, nt)

			emailSent, err := uc.InviteUser(context.Background(), dto.InviteUserIn{
				TeamID: 1, UserID: 1, Email: "invitee@example.com",
			})

			require.NoError(t, err, "участник уже добавлен — отказ почты не повод ронять запрос")
			require.False(t, emailSent, "флаг должен честно говорить, что письмо не ушло")
		})
	}
}
