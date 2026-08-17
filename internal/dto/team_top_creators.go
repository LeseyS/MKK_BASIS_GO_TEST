package dto

type TeamTopCreator struct {
	TeamID       int64  `json:"team_id"`
	TeamName     string `json:"team_name"`
	UserID       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	UserEmail    string `json:"user_email"`
	TasksCreated int64  `json:"tasks_created"`
	Position     int    `json:"position"`
}
