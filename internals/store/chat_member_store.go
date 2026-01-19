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

// New struct - for API responses with user details
type ChatMemberWithUser struct {
	ChatID        int64     `json:"chatId"`
	UserID        int64     `json:"userId"`
	Role          string    `json:"role"`
	JoinedAt      time.Time `json:"joinedAt"`
	Username      string    `json:"username"`
	Email         string    `json:"email,omitempty"`
	ProfilePicURL string    `json:"profilePicUrl,omitempty"`
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
    GetChatMembers(chatID int64) ([]*ChatMemberWithUser, error)
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

func (pg *PostgresChatMemberStore) GetChatMembers(chatID int64) ([]*ChatMemberWithUser, error) {
	query := `
		SELECT 
			cm.chat_id,
			cm.user_id,
			cm.role,
			cm.joined_at,
			u.username,
			u.email,
			u.avatar_url  -- if you have this
		FROM chat_members cm
		JOIN users u ON cm.user_id = u.id
		WHERE cm.chat_id = $1
		ORDER BY cm.joined_at DESC
	`
	
	rows, err := pg.db.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var members []*ChatMemberWithUser
	for rows.Next() {
		member := &ChatMemberWithUser{}
		err := rows.Scan(
			&member.ChatID,
			&member.UserID,
			&member.Role,
			&member.JoinedAt,
			&member.Username,
			&member.Email,   
			&member.ProfilePicURL,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	
	return members, rows.Err()
}

func (pg *PostgresChatMemberStore) GetUserRole(chatID, userID int64) (string, error) {
	role := ""
	query := `
		SELECT role FROM chat_members
		WHERE chat_id = $1 AND user_id = $2
	`

	err := pg.db.QueryRow(query, chatID, userID).Scan(&role)
	if err != nil {
		return "",err
	}
	fmt.Println(role)
	return role, nil
}

func (pg *PostgresChatMemberStore) IsMember(chatID, userID int64) (bool, error) {
	role := ""
	query := `
		SELECT role FROM chat_members
		WHERE chat_id = $1 AND user_id = $2
	`

	err := pg.db.QueryRow(query, chatID, userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (pg *PostgresChatMemberStore) UpdateLastRead(chatID, userID, messageID int64) error {
	query := `
		UPDATE chat_members
		SET last_read_message_id = $1
		WHERE chat_id = $2 AND user_id = $3
		AND (last_read_message_id IS NULL OR last_read_message_id < $1);
	`

	_, err := pg.db.Exec(query, messageID, chatID, userID)
	if err != nil {
		return err
	}

	return nil
}