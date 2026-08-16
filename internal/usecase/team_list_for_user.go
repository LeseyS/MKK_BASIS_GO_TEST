package usecase

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
)

func (uc *UseCase) TeamListForUser(ctx context.Context, userID int64) ([]dto.TeamList, error) {
	var out []dto.TeamList
	teams, err := uc.mysql.TeamListForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("mysql.GetTeamListForUser: %w", err)
	}

	for _, t := range teams {
		out = append(out, dto.TeamList{
			Team: t,
		})
	}

	return out, nil
}
