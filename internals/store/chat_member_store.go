package store

import "time"


type ChatGroupRole string
const (
	OWNER ChatGroupRole = "owner"
	ADMIN ChatGroupRole = "admin"
	MEMBER ChatGroupRole = "member"
)

type ChatMembers struct {
	ChatID int `json:"chat_id"`
	UserID int `json:"user_id"`
	Role ChatGroupRole `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

