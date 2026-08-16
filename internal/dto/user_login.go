package dto

import "github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"

type UserLoginIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLoginOut struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}
