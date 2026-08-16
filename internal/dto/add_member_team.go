package dto

import "github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"

type AddMemberTeam struct {
	TeamID int64
	UserID int64
	Role   domain.Role
}
