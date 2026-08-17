package usecase

import (
	"context"
	"strings"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/dto"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/passwordhash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_HashesPassword(t *testing.T) {
	const password = "password123"

	var stored domain.User

	my := &mockMySQL{
		createUser: func(_ context.Context, u domain.User) (int64, error) {
			stored = u
			return 42, nil
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	out, err := uc.CreateUser(context.Background(), dto.CreateUserIn{
		Name: "alice", Email: "alice@example.com", Password: password,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(42), out.ID)
	assert.NotEqual(t, password, stored.PasswordHash, "пароль не должен храниться в открытом виде")

	ok, err := passwordhash.Verify(password, stored.PasswordHash)
	require.NoError(t, err)
	assert.True(t, ok, "сохранённый хеш должен проверяться исходным паролем")
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	my := &mockMySQL{
		createUser: func(context.Context, domain.User) (int64, error) {
			return 0, apperr.ErrDuplicate
		},
	}
	uc := newUseCase(my, nil, nil, nil)

	_, err := uc.CreateUser(context.Background(), dto.CreateUserIn{
		Name: "alice", Email: "alice@example.com", Password: "password123",
	})

	require.ErrorIs(t, err, apperr.ErrEmailTaken)
}

func TestCreateUser_PasswordLongerThanBcryptLimit(t *testing.T) {
	uc := newUseCase(&mockMySQL{
		createUser: func(context.Context, domain.User) (int64, error) {
			t.Fatal("до сохранения дойти не должно")
			return 0, nil
		},
	}, nil, nil, nil)

	_, err := uc.CreateUser(context.Background(), dto.CreateUserIn{
		Name: "alice", Email: "alice@example.com", Password: strings.Repeat("п", 72),
	})

	require.ErrorIs(t, err, passwordhash.ErrPasswordTooLong,
		"72 кириллических символа — это 144 байта, bcrypt их не принимает")
}

func TestCreateUser_InvalidName(t *testing.T) {
	uc := newUseCase(&mockMySQL{
		createUser: func(context.Context, domain.User) (int64, error) {
			t.Fatal("до сохранения дойти не должно")
			return 0, nil
		},
	}, nil, nil, nil)

	_, err := uc.CreateUser(context.Background(), dto.CreateUserIn{
		Name: "ab", Email: "alice@example.com", Password: "password123",
	})

	require.ErrorIs(t, err, apperr.ErrValidation)
}
