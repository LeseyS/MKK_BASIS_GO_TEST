package usecase

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
)

func (uc *UseCase) TeamTopCreators(ctx context.Context, userID int64) ([]dto.TeamTopCreator, error) {
	top, err := uc.mysql.TeamTopCreators(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.mysql.TeamTopCreators: %w", err)
	}

	return top, nil
}
