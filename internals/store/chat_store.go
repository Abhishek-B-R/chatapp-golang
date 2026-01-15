package store

import (
	"database/sql"
	"fmt"
	"time"
)

type ChatType string
const (
	GROUP ChatType = "group"
	DM ChatType = "dm"
)

type Chat struct {
	ChatID int64 `json:"chat_id"`
	Type ChatType `json:"type"`
	Name *string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatMembers struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
	Role string `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type PostgresChatStore struct {
	db *sql.DB
}

func NewPostgresChatStore(db *sql.DB) *PostgresChatStore {
	return &PostgresChatStore{db:db}
}

type ChatStore interface {
	CreateChat(chat *Chat) error
	GetChatByID(chatID int64) (*Chat, error)
	UpdateChat(chat *Chat) error
	DeleteChat(chatID int64) error
}

func (pg *PostgresChatStore) CreateChat(chat *Chat) error {
	fmt.Println("Chat created")
	return nil
}

func (pg *PostgresChatStore) DeleteChat(chatID int64) error {
	fmt.Println("Chat deleted")
	return nil
}

func (pg *PostgresChatStore) GetChatByID(chatID int64) (*Chat, error) {
	fmt.Println("have to implement it still")
	return nil, nil
}

func (pg *PostgresChatStore) UpdateChat(chatID int64) error {
	fmt.Println("Chat updated")
	return nil
}