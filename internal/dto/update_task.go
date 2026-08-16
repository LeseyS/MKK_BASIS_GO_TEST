package dto

import "github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"

type UpdateTaskIn struct {
	ActorID     int64   `json:"-"`
	TaskID      int64   `json:"-"`
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty" validate:"omitempty,oneof=todo in_progress done"`
	AssigneeID  *int64  `json:"assignee_id,omitempty"`
}

type UpdateTaskOut struct {
	Task domain.Task `json:"task"`
}
