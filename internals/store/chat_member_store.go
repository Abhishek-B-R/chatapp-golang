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
	ChatID int `json:"chat_id"`
	UserID int `json:"user_id"`
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
	AddMember(chatID, userID int, role string) error
    RemoveMember(chatID, userID int) error
    GetChatMembers(chatID int) ([]*ChatMember, error)
    GetUserRole(chatID, userID int) (string, error)
    IsMember(chatID, userID int) (bool, error)
    UpdateLastRead(chatID, userID, messageID int) error
}

func (pg *PostgresChatMemberStore) AddMember(chatID, userID int, role string) error {
	fmt.Println("Added new member")
	return nil
}

func (pg *PostgresChatMemberStore) RemoveMember(chatID, userID int) error {
	fmt.Println("removed member")
	return nil
}

func (pg *PostgresChatMemberStore) GetChatMembers(chatID int) ([]*ChatMember, error){
	fmt.Println("Fetching chat members")
	return nil, nil
}

func (pg *PostgresChatMemberStore) GetUserRole(chatID, userID int) (string, error) {
	fmt.Println("Fetching user role")
	return "member", nil
}

func (pg *PostgresChatMemberStore) IsMember(chatID, userID int) (bool, error) {
	fmt.Println("checking if this user is part of this chat")
	return false, nil
}

func (pg *PostgresChatMemberStore) UpdateLastRead(chatID, userID, messageID int) error {
	fmt.Println("updating latest info")
	return nil
}