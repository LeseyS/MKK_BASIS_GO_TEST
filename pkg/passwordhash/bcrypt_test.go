package passwordhash

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashAndVerify(t *testing.T) {
	const password = "correct-horse-battery"

	hash, err := Hash(password)
	require.NoError(t, err)
	assert.NotEqual(t, password, hash)

	ok, err := Verify(password, hash)
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestVerify_WrongPassword(t *testing.T) {
	hash, err := Hash("correct-horse-battery")
	require.NoError(t, err)

	ok, err := Verify("wrong-password", hash)

	require.NoError(t, err, "несовпадение пароля — не ошибка, а отрицательный ответ")
	assert.False(t, ok)
}

func TestVerify_MalformedHash(t *testing.T) {
	ok, err := Verify("password", "not-a-bcrypt-hash")

	require.Error(t, err, "битый хеш в базе должен быть отличим от неверного пароля")
	assert.False(t, ok)
}

func TestHash_SaltMakesHashesUnique(t *testing.T) {
	first, err := Hash("same-password")
	require.NoError(t, err)

	second, err := Hash("same-password")
	require.NoError(t, err)

	assert.NotEqual(t, first, second, "одинаковые пароли должны давать разные хеши")
}

func TestHash_LengthLimit(t *testing.T) {
	cases := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"72 байта", strings.Repeat("a", 72), false},
		{"73 байта", strings.Repeat("a", 73), true},
		{"72 кириллических символа = 144 байта", strings.Repeat("п", 72), true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Hash(c.password)

			if c.wantErr {
				require.ErrorIs(t, err, ErrPasswordTooLong)
				return
			}
			require.NoError(t, err)
		})
	}
}
