package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/passwordhash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserLogin_UnknownEmail(t *testing.T) {
	my := &mockMySQL{
		getUserByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{}, apperr.ErrNotFound
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.UserLogin(context.Background(), dto.UserLoginIn{
		Email: "nobody@example.com", Password: "password123",
	})

	require.ErrorIs(t, err, apperr.ErrInvalidCredentials)
	assert.NotErrorIs(t, err, apperr.ErrNotFound,
		"несуществующий email и неверный пароль должны быть неотличимы снаружи")
}

func TestUserLogin_WrongPassword(t *testing.T) {
	hash, err := passwordhash.Hash("correct-password")
	require.NoError(t, err)

	my := &mockMySQL{
		getUserByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{ID: 1, Email: "user@example.com", PasswordHash: hash}, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err = uc.UserLogin(context.Background(), dto.UserLoginIn{
		Email: "user@example.com", Password: "wrong-password",
	})

	require.ErrorIs(t, err, apperr.ErrInvalidCredentials)
}

func TestUserLogin_Success(t *testing.T) {
	hash, err := passwordhash.Hash("correct-password")
	require.NoError(t, err)

	my := &mockMySQL{
		getUserByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{ID: 42, Email: "user@example.com", PasswordHash: hash}, nil
		},
	}
	ti := &mockTokenIssuer{
		generate: func(userID int64, email string) (string, error) {
			assert.Equal(t, int64(42), userID)
			assert.Equal(t, "user@example.com", email)
			return "signed-token", nil
		},
	}
	uc := newUseCase(my, nil, ti, nil)

	out, err := uc.UserLogin(context.Background(), dto.UserLoginIn{
		Email: "user@example.com", Password: "correct-password",
	})

	require.NoError(t, err)
	assert.Equal(t, "signed-token", out.Token)
	assert.Equal(t, int64(42), out.User.ID)

	body, err := json.Marshal(out)
	require.NoError(t, err)
	assert.NotContains(t, string(body), hash, "хеш пароля не должен попадать в ответ")
}

func TestUserLogin_TokenIssuerFailure(t *testing.T) {
	hash, err := passwordhash.Hash("correct-password")
	require.NoError(t, err)

	my := &mockMySQL{
		getUserByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{ID: 1, PasswordHash: hash}, nil
		},
	}
	ti := &mockTokenIssuer{
		generate: func(int64, string) (string, error) {
			return "", errors.New("signing key missing")
		},
	}
	uc := newUseCase(my, nil, ti, nil)

	out, err := uc.UserLogin(context.Background(), dto.UserLoginIn{
		Email: "user@example.com", Password: "correct-password",
	})

	require.Error(t, err)
	assert.Empty(t, out.Token)
}
