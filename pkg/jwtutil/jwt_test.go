package jwtutil

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndParse(t *testing.T) {
	issuer := NewIssuer("secret", time.Hour)

	token, err := issuer.Generate(42, "user@example.com")
	require.NoError(t, err)

	claims, err := issuer.Parse(token)
	require.NoError(t, err)
	require.Equal(t, int64(42), claims.UserID)
	require.Equal(t, "user@example.com", claims.Email)
	require.Equal(t, "42", claims.Subject)
}

func TestParse_ExpiredToken(t *testing.T) {
	issuer := NewIssuer("secret", -time.Minute)

	token, err := issuer.Generate(1, "user@example.com")
	require.NoError(t, err)

	_, err = issuer.Parse(token)

	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestParse_WrongSecret(t *testing.T) {
	token, err := NewIssuer("secret", time.Hour).Generate(1, "user@example.com")
	require.NoError(t, err)

	_, err = NewIssuer("another-secret", time.Hour).Parse(token)

	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestParse_Garbage(t *testing.T) {
	issuer := NewIssuer("secret", time.Hour)

	for _, token := range []string{"", "not-a-token", "a.b.c"} {
		_, err := issuer.Parse(token)
		require.ErrorIs(t, err, ErrInvalidToken, "token=%q", token)
	}
}

func TestParse_RejectsNoneAlgorithm(t *testing.T) {
	claims := Claims{
		UserID: 42,
		Email:  "attacker@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	unsigned, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).
		SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = NewIssuer("secret", time.Hour).Parse(unsigned)

	require.ErrorIs(t, err, ErrInvalidToken, "токен с alg=none принимать нельзя")
}
