package dto

type TeamStats struct {
	TeamID        int64  `json:"team_id"`
	Name          string `json:"name"`
	MembersCount  int64  `json:"members_count"`
	DoneLast7Days int64  `json:"done_last_7_days"`
}
