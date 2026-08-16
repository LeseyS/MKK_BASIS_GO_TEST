package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
	"github.com/sony/gobreaker"
)

func (uc *UseCase) InviteUser(ctx context.Context, in dto.InviteUserIn) (bool, error) {
	var team domain.Team

	err := transaction.Wrap(ctx, func(ctx context.Context) error {
		role, err := uc.mysql.GetMemberRole(ctx, in.TeamID, in.UserID)
		if err != nil {
			if errors.Is(err, apperr.ErrNotFound) {
				return apperr.ErrForbidden
			}
			return fmt.Errorf("uc.mysql.GetMemberRole: %w", err)
		}

		if role != domain.RoleOwner && role != domain.RoleAdmin {
			return apperr.ErrForbidden
		}

		invite, err := uc.mysql.GetUserByEmail(ctx, in.Email)
		if err != nil {
			if errors.Is(err, apperr.ErrNotFound) {
				return fmt.Errorf("%w: no user registered with that email", apperr.ErrNotFound)
			}
			return fmt.Errorf("uc.mysql.GetUserByEmail: %w", err)
		}

		if err := uc.mysql.AddMemberTeam(ctx, dto.AddMemberTeam{
			TeamID: in.TeamID,
			UserID: invite.ID,
			Role:   domain.RoleMember,
		}); err != nil {
			if errors.Is(err, apperr.ErrDuplicate) {
				return apperr.ErrAlreadyMember
			}
			return fmt.Errorf("uc.mysql.AddMemberTeam: %w", err)
		}

		team, err = uc.mysql.GetTeamByID(ctx, in.TeamID)
		if err != nil {
			return fmt.Errorf("uc.mysql.GetTeamByID: %w", err)
		}

		return nil
	})
	if err != nil {
		return false, fmt.Errorf("transaction.Wrap: %w", err)
	}

	if err := uc.notifier.SendInvite(ctx, in.Email, team.Name); err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return true, nil
		}
		return true, fmt.Errorf("uc.notifier.SendInvite: %w", err)
	}

	return true, nil
}
