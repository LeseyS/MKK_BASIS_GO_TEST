package dto

import "time"

type TaskHistoryIn struct {
	TaskID  int64
	ActorID int64
	Limit   int
	Offset  int
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

type TaskHistoryOut struct {
	Entries []TaskHistoryEntry
	Total   int64
	Limit   int
	Offset  int
}
