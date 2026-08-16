package dto

type InviteUserIn struct {
	TeamID int64  `json:"-"`
	UserID int64  `json:"-"`
	Email  string `json:"email" validate:"required,email"`
}

type InviteUserOut struct {
	Invited   bool `json:"invited"`
	EmailSent bool `json:"email_sent"`
}
