package store

import (
	"time"
)

type UserStore struct {
	UserID int64 `json:"user_id"`
	Name string `json:"name"`
	Email string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}