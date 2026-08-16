package dto

import (
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
)

type CreateTeamIn struct {
	Name    string `json:"name" validate:"required"`
	OwnerID int64  `json:"owner_id"`
}

type CreateTeamOut struct {
	Team domain.Team `json:"team"`
}
