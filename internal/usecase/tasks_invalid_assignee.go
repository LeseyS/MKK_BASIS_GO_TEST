package usecase

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
)

func (uc *UseCase) TasksInvalidAssignee(ctx context.Context, userID int64) ([]dto.TaskInvalidAssignee, error) {
	tasks, err := uc.mysql.TasksInvalidAssignee(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.mysql.TasksInvalidAssignee: %w", err)
	}

	return tasks, nil
}
