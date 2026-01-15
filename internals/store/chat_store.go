package store

import (
	"database/sql"
	"fmt"
	"time"
)

type Chat struct {
	ChatID int64 `json:"chat_id"`
	IsGroup bool `json:"is_group"`
	Name *string `json:"name"`
	CreatedBy int64 `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type PostgresChatStore struct {
	db *sql.DB
}

func NewPostgresChatStore(db *sql.DB) *PostgresChatStore {
	return &PostgresChatStore{db:db}
}

type ChatStore interface {
	CreateChat(chat *Chat) (*Chat, error)
	GetUserChats(userID int64) ([]*Chat, error)
	GetChatByID(chatID int64) (*Chat, error)
	UpdateChat(chat *Chat) error
	DeleteChat(chatID int64) error
}

func (pg *PostgresChatStore) CreateChat(chat *Chat) (*Chat, error) {
	fmt.Println("Chat created")
	return nil,nil
}

func (pg *PostgresChatStore) DeleteChat(chatID int64) error {
	fmt.Println("Chat deleted")
	return nil
}

func (pg *PostgresChatStore) GetChatByID(chatID int64) (*Chat, error) {
	fmt.Println("have to implement it still")
	return nil, nil
}

func (pg *PostgresChatStore) UpdateChat(chat *Chat) error {
	fmt.Println("Chat updated")
	return nil
}

func (pg *PostgresChatStore) GetUserChats(userID int64) ([]*Chat, error) {
	fmt.Println("user chats searching")
	return nil, nil
}