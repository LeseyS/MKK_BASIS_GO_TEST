package domain

import (
	"fmt"
	"time"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/validate"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name" validate:"required,min=3,max=64"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewUser(email string, name string, password string) (User, error) {
	u := User{
		Email:        email,
		Name:         name,
		PasswordHash: password,
	}

	if err := u.Validate(); err != nil {
		return u, fmt.Errorf("validate user: %w", err)
	}

	return u, nil
}

func (u User) Validate() error {
	if err := validate.V.Struct(u); err != nil {
		return apperr.ErrValidation
	}
	return nil
}
