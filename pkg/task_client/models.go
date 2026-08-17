package task_client

import "time"

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Team struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID          int64     `json:"id"`
	TeamID      int64     `json:"team_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssigneeID  *int64    `json:"assignee_id"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID int64 `json:"id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type CreateTeamResponse struct {
	Team Team `json:"team"`
}

type TeamListItem struct {
	Team Team `json:"team"`
}

type InviteRequest struct {
	Email string `json:"email"`
}

type InviteResponse struct {
	Invited   bool `json:"invited"`
	EmailSent bool `json:"email_sent"`
}

type TeamStats struct {
	TeamID        int64  `json:"team_id"`
	Name          string `json:"name"`
	MembersCount  int64  `json:"members_count"`
	DoneLast7Days int64  `json:"done_last_7_days"`
}

type teamStatsResponse struct {
	Teams []TeamStats `json:"teams"`
}

type TeamTopCreator struct {
	TeamID       int64  `json:"team_id"`
	TeamName     string `json:"team_name"`
	UserID       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	UserEmail    string `json:"user_email"`
	TasksCreated int64  `json:"tasks_created"`
	Position     int    `json:"position"`
}

type topCreatorsResponse struct {
	TopCreators []TeamTopCreator `json:"top_creators"`
}

type CreateTaskRequest struct {
	TeamID      int64  `json:"team_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	AssigneeID  *int64 `json:"assignee_id,omitempty"`
}

type CreateTaskResponse struct {
	Task Task `json:"task"`
}

type ListTasksRequest struct {
	TeamID     int64
	Status     string
	AssigneeID int64
	Limit      int
	Offset     int
}

type ListTasksResponse struct {
	Tasks  []Task `json:"tasks"`
	Total  int64  `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	AssigneeID  *int64  `json:"assignee_id,omitempty"`
}

type UpdateTaskResponse struct {
	Task Task `json:"task"`
}

type TaskHistoryEntry struct {
	ID             int64     `json:"id"`
	Field          string    `json:"field"`
	OldValue       *string   `json:"old_value"`
	NewValue       *string   `json:"new_value"`
	ChangedBy      int64     `json:"changed_by"`
	ChangedByName  string    `json:"changed_by_name"`
	ChangedByEmail string    `json:"changed_by_email"`
	ChangedAt      time.Time `json:"changed_at"`
}

type TaskHistoryResponse struct {
	History []TaskHistoryEntry `json:"history"`
	Total   int64              `json:"total"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
}

type TaskInvalidAssignee struct {
	TaskID        int64  `json:"task_id"`
	TeamID        int64  `json:"team_id"`
	TeamName      string `json:"team_name"`
	Title         string `json:"title"`
	Status        string `json:"status"`
	AssigneeID    int64  `json:"assignee_id"`
	AssigneeEmail string `json:"assignee_email"`
}

type invalidAssigneesResponse struct {
	Tasks []TaskInvalidAssignee `json:"tasks"`
}
