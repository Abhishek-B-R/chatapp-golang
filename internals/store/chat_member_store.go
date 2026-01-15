package store

import (
	"database/sql"
	"fmt"
	"time"
)


type ChatGroupRole string
const (
	OWNER ChatGroupRole = "owner"
	ADMIN ChatGroupRole = "admin"
	MEMBER ChatGroupRole = "member"
)

type ChatMember struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
	Role ChatGroupRole `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type PostgresChatMemberStore struct{
	db *sql.DB
}

func NewPostgresChatMemberStore(db *sql.DB) *PostgresChatMemberStore{
	return &PostgresChatMemberStore{
		db: db,
	}
}

type ChatMemberStore interface {
	AddMember(chatID, userID int64, role string) error
    RemoveMember(chatID, userID int64) error
    GetChatMembers(chatID int64) ([]*ChatMember, error)
    GetUserRole(chatID, userID int64) (string, error)
    IsMember(chatID, userID int64) (bool, error)
    UpdateLastRead(chatID, userID, messageID int64) error
}

func (pg *PostgresChatMemberStore) AddMember(chatID, userID int64, role string) error {
	fmt.Println("Added new member")
	return nil
}

func (pg *PostgresChatMemberStore) RemoveMember(chatID, userID int64) error {
	fmt.Println("removed member")
	return nil
}

func (pg *PostgresChatMemberStore) GetChatMembers(chatID int64) ([]*ChatMember, error){
	fmt.Println("Fetching chat members")
	return nil, nil
}

func (pg *PostgresChatMemberStore) GetUserRole(chatID, userID int64) (string, error) {
	fmt.Println("Fetching user role")
	return "member", nil
}

func (pg *PostgresChatMemberStore) IsMember(chatID, userID int64) (bool, error) {
	fmt.Println("checking if this user is part of this chat")
	return false, nil
}

func (pg *PostgresChatMemberStore) UpdateLastRead(chatID, userID, messageID int64) error {
	fmt.Println("updating latest info")
	return nil
}