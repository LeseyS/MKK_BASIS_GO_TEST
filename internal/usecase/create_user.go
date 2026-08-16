package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/passwordhash"
)

func (uc *UseCase) CreateUser(ctx context.Context, in dto.CreateUserIn) (dto.CreateUserOut, error) {
	var output dto.CreateUserOut

	hash, err := passwordhash.Hash(in.Password)
	if err != nil {
		return output, fmt.Errorf("hashing password: %w", err)
	}

	u, err := domain.NewUser(in.Email, in.Name, hash)
	if err != nil {
		return output, fmt.Errorf("domain.NewUser: %w", err)
	}

	id, err := uc.mysql.CreateUser(ctx, u)
	if err != nil {
		if errors.Is(err, apperr.ErrDuplicate) {
			return output, apperr.ErrEmailTaken
		}
		return output, fmt.Errorf("mysql.CreateUser: %w", err)
	}

	output.ID = id

	return output, nil
}
