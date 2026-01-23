package store

import (
	"context"
	"database/sql"
	"time"
)

type Chat struct {
	ChatID int64 `json:"id"`
	IsGroup bool `json:"is_group"`
	Name *string `json:"name"`
	CreatedBy int64 `json:"created_by"`
	LastMessageAt *int64 `json:"last_message_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostgresChatStore struct {
	db *sql.DB
}

func NewPostgresChatStore(db *sql.DB) *PostgresChatStore {
	return &PostgresChatStore{db:db}
}

type ChatStore interface {
	CreateChat(ctx context.Context, chat *Chat, userID int64) (*Chat, error)
	GetUserChats(ctx context.Context, userID int64) (*[]Chat, error)
	GetChatByID(ctx context.Context, chatID int64) (*Chat, error)
	UpdateChat(ctx context.Context, chat *Chat, chatID int64) error
	DeleteChat(ctx context.Context, chatID int64) error
}

func (pg *PostgresChatStore) CreateChat(ctx context.Context, chat *Chat, userID int64) (*Chat, error) {
	tx, err := pg.db.Begin()
	defer func() {
		_ = tx.Rollback()
	}()

	if err != nil {
		return nil, err
	}

	q1 := `
		INSERT INTO chats (is_group, name, created_by)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, q1, chat.IsGroup, chat.Name, userID).Scan(&chat.ChatID)
	if err != nil {
		return nil, err
	}

	q2 := `
		INSERT INTO chat_members (user_id, chat_id, role, muted)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.ExecContext(ctx, q2, userID, chat.ChatID, "admin", false)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return chat, nil
}

func (pg *PostgresChatStore) UpdateChat(ctx context.Context, chat *Chat, chatID int64) error {
	query := `
		UPDATE chats
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id
	`

	err := pg.db.QueryRowContext(ctx, query, chat.Name, chatID).Scan(&chat.ChatID)
	return err
}

func (pg *PostgresChatStore) DeleteChat(ctx context.Context, chatID int64) error {
	query := `
		DELETE from chats
		WHERE id = $1
	`

	results, err := pg.db.ExecContext(ctx, query, chatID)
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

func (pg *PostgresChatStore) GetChatByID(ctx context.Context, chatID int64) (*Chat, error) {
	var chat Chat
	query := `
		SELECT id, is_group, name, created_by, last_message_at, created_at, updated_at
		FROM chats
		WHERE id = $1
	`

	err := pg.db.QueryRowContext(ctx, query, chatID).Scan(
		&chat.ChatID,  
		&chat.IsGroup, 
		&chat.Name, 
		&chat.CreatedBy, 
		&chat.LastMessageAt, 
		&chat.CreatedAt, 
		&chat.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return &chat, nil
}

func (pg *PostgresChatStore) GetUserChats(ctx context.Context, userID int64) (*[]Chat, error) {
	var chats []Chat
	query := `
		SELECT 
			c.id, 
			c.is_group, 
			c.name, 
			c.created_by, 
			c.last_message_at, 
			c.created_at, 
			c.updated_at
		FROM chats c
		INNER JOIN chat_members cm ON c.id = cm.chat_id
		WHERE cm.user_id = $1
		ORDER BY c.last_message_at DESC NULLS LAST
	`

	rows, err := pg.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat Chat
		err = rows.Scan(
			&chat.ChatID, 
			&chat.IsGroup, 
			&chat.Name, 
			&chat.CreatedBy, 
			&chat.LastMessageAt, 
			&chat.CreatedAt, 
			&chat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &chats, err
}