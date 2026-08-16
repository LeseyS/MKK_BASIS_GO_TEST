package apperr

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrDuplicate          = errors.New("already exists")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidAssignee    = errors.New("invalid assignee")
	ErrAlreadyExists      = errors.New("already exists")
	ErrForbidden          = errors.New("forbidden: insufficient role")
	ErrAlreadyMember      = errors.New("user is already a member of this team")
	ErrValidation         = errors.New("validation error")
)
