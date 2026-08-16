package dto

import "github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"

type CreateTaskOut struct {
	Task domain.Task `json:"task"`
}

type CreateTaskIn struct {
	ActorID     int64  `json:"-"`
	TeamID      int64  `json:"team_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	AssigneeID  *int64 `json:"assignee_id"`
}
