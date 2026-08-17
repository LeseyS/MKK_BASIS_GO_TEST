package dto

type TaskInvalidAssignee struct {
	TaskID        int64  `json:"task_id"`
	TeamID        int64  `json:"team_id"`
	TeamName      string `json:"team_name"`
	Title         string `json:"title"`
	Status        string `json:"status"`
	AssigneeID    int64  `json:"assignee_id"`
	AssigneeEmail string `json:"assignee_email"`
}
