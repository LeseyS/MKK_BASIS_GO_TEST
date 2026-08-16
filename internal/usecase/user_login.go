package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/passwordhash"
)

func (uc *UseCase) UserLogin(ctx context.Context, in dto.UserLoginIn) (dto.UserLoginOut, error) {
	var out dto.UserLoginOut

	u, err := uc.mysql.GetUserByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return out, apperr.ErrInvalidCredentials
		}

		return out, fmt.Errorf("domain.GetUserByEmail: %w", err)
	}

	ok, err := passwordhash.Verify(in.Password, u.PasswordHash)
	if err != nil || !ok {
		return out, apperr.ErrInvalidCredentials
	}

	out.Token, err = uc.tokenIssuer.Generate(u.ID, u.Email)
	if err != nil {
		return out, fmt.Errorf("tokenIssuer.Generate: %w", err)
	}

	out.User = u

	return out, nil
}
