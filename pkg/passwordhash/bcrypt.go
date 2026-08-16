package passwordhash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const cost = 12

var ErrPasswordTooLong = errors.New("password exceeds maximum length of 72 bytes")

func Hash(password string) (string, error) {
	if len(password) > 72 {
		return "", ErrPasswordTooLong
	}
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Verify(password, encoded string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(encoded), []byte(password))
	if err == nil {
		return true, nil
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, err
}
