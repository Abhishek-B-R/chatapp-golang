package store

import (
	"database/sql"
	"fmt"
	"strings"
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
	LastReadMessageID int64 `json:"last_read_message_id"`
	JoinedAt time.Time `json:"joined_at"`
	Muted bool `json:"muted"`
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
	role = strings.TrimSpace(strings.ToLower(role))

	switch role {
	case "owner", "admin", "member":
	default:
		fmt.Printf("WARN: invalid role: %q", role)
		role = "member"
	}

	query := `
		INSERT INTO chat_members (user_id, chat_id, role, muted)
		VALUES ($1, $2, $3, $4)
	`

	_, err := pg.db.Exec(query, userID, chatID, role, false)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresChatMemberStore) RemoveMember(chatID, userID int64) error {
	query := `
		DELETE FROM chat_members 
		WHERE chat_id = $1 AND user_id = $2
	`

	results, err := pg.db.Exec(query, chatID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

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