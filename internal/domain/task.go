package domain

import (
	"fmt"
	"time"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/validate"
)

type TaskStatus string

type Task struct {
	ID          int64      `json:"id"`
	TeamID      int64      `json:"team_id"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status" validate:"required,oneof=todo in_progress done"`
	AssigneeID  *int64     `json:"assignee_id"`
	CreatedBy   int64      `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

func NewTask(teamID int64, title, description string, status TaskStatus, createdBy int64, assigneeID *int64) (Task, error) {
	if status == "" {
		status = StatusTodo
	}

	t := Task{
		TeamID:      teamID,
		Title:       title,
		Description: description,
		Status:      status,
		AssigneeID:  assigneeID,
		CreatedBy:   createdBy,
	}

	if err := t.Validate(); err != nil {
		return Task{}, fmt.Errorf("validate task: %w", err)
	}

	return t, nil
}

func (u Task) Validate() error {
	if err := validate.V.Struct(u); err != nil {
		return apperr.ErrValidation
	}
	return nil
}
