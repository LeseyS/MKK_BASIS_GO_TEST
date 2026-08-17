package domain

import (
	"strings"
	"testing"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/stretchr/testify/require"
)

func TestNewUser_NameValidation(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"слишком короткое", "ab", true},
		{"минимально допустимое", "abc", false},
		{"обычное", "alice", false},
		{"максимально допустимое", strings.Repeat("a", 64), false},
		{"слишком длинное", strings.Repeat("a", 65), true},
		{"пустое", "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewUser("user@example.com", c.value, "hash")

			if c.wantErr {
				require.ErrorIs(t, err, apperr.ErrValidation)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewUser_KeepsFields(t *testing.T) {
	u, err := NewUser("user@example.com", "alice", "stored-hash")

	require.NoError(t, err)
	require.Equal(t, "user@example.com", u.Email)
	require.Equal(t, "alice", u.Name)
	require.Equal(t, "stored-hash", u.PasswordHash)
}
