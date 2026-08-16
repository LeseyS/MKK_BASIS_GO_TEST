package dto

import "github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"

type ListTasks struct {
	TeamID     int64
	UserID     int64
	Status     string
	AssigneeID int64
	Limit      int
	Offset     int
}

type ListTasksOut struct {
	Tasks  []domain.Task
	Total  int64
	Limit  int
	Offset int
}
