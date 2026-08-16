package dto

type CreateUserIn struct {
	Name     string `json:"name" validate:"required,min=3,max=64"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}
type CreateUserOut struct {
	ID int64 `json:"id"`
}
