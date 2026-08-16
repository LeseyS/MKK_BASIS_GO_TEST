package usecase

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (uc *UseCase) CreateWithOwner(ctx context.Context, in dto.CreateTeamIn) (dto.CreateTeamOut, error) {
	var out dto.CreateTeamOut

	err := transaction.Wrap(ctx, func(ctx context.Context) error {
		teamID, err := uc.mysql.CreateTeam(ctx, in)
		if err != nil {
			return fmt.Errorf("mysql.CreateTeam: %w", err)
		}

		err = uc.mysql.CreateTeamMembers(ctx, teamID, in.OwnerID, domain.RoleOwner)
		if err != nil {
			return fmt.Errorf("mysql.CreateTeamMembers: %w", err)
		}

		team, err := uc.mysql.GetTeamByID(ctx, teamID)
		if err != nil {
			return fmt.Errorf("mysql.GetTeamByID: %w", err)
		}

		out.Team = team

		return nil
	})

	if err != nil {
		return out, fmt.Errorf("transaction.Wrap: %w", err)
	}

	return out, nil
}
