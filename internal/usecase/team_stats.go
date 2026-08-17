package usecase

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
)

func (uc *UseCase) TeamStats(ctx context.Context, userID int64) ([]dto.TeamStats, error) {
	stats, err := uc.mysql.TeamStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.mysql.TeamStats: %w", err)
	}

	return stats, nil
}
