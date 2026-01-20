package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type MessageType string
type AttachmentType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeSystem MessageType = "system"
)

const (
	AttachmentImage AttachmentType = "image"
	AttachmentVideo AttachmentType = "video"
	AttachmentPDF   AttachmentType = "pdf"
	AttachmentFile  AttachmentType = "file"
)

type Message struct {
	ID               int64       `json:"id"`
	ChatID            int64       `json:"chat_id"`
	SenderID          *int64      `json:"sender_id,omitempty"`

	Type              MessageType `json:"type"`
	Content           *string     `json:"content,omitempty"`

	ReplyToMessageID  *int64      `json:"reply_to_message_id,omitempty"`

	EditedAt          *time.Time  `json:"edited_at,omitempty"`
	DeletedAt         *time.Time  `json:"deleted_at,omitempty"`
	CreatedAt         time.Time   `json:"created_at"`

	Attachments       []MessageAttachment `json:"attachments,omitempty"`
}

type MessageAttachment struct {
	ID        int64          `json:"id"`
	MessageID int64          `json:"message_id"`
	Type      AttachmentType `json:"type"`
	URL       string         `json:"url"`
	Filename  *string        `json:"filename,omitempty"`
	SizeBytes *int64         `json:"size_bytes,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

type PostgresMessageStore struct {
	db *sql.DB
}

func NewPostgresMessageStore(db *sql.DB) *PostgresMessageStore {
	return &PostgresMessageStore{
		db: db,
	}
}

type MessageStore interface{
	CreateMessage(msg *Message) error
    GetMessage(id int64) (*Message, error)
    GetChatMessages(chatID, limit, offset int64) ([]*Message, error)
    UpdateMessage(msg *Message) error
    DeleteMessage(id int64) error // soft delete
    GetUnreadCount(chatID, userID int64) (int64, error)
}

func (pg *PostgresMessageStore) CreateMessage(msg *Message) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q1 := `
		INSERT INTO messages (chat_id, sender_id, type, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err = tx.QueryRow(q1, msg.ChatID, msg.SenderID, msg.Type, msg.Content).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return err
	}

	if len(msg.Attachments) > 0 {
		q2 := `
			INSERT INTO message_attachments (message_id, type, url, filename, size_bytes, metadata)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at
		`

		for i := range msg.Attachments {
			a := &msg.Attachments[i]

			if a.Metadata == nil {
				a.Metadata = json.RawMessage(`{}`)
			}

			err := tx.QueryRow(
				q2,
				msg.ID,
				a.Type,
				a.URL,
				a.Filename,
				a.SizeBytes,
				a.Metadata,
			).Scan(&a.ID, &a.CreatedAt)

			if err != nil {
				return err
			}

			a.MessageID = msg.ID
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (pg *PostgresMessageStore) GetMessage(id int64) (*Message, error) {
	fmt.Println("getting message with id: ",id)
	return nil, nil
}

func (pg *PostgresMessageStore) GetChatMessages(chatID, limit, offset int64) ([]*Message, error) {
	fmt.Println("getting chat messages")
	return nil, nil
}

func (pg *PostgresMessageStore) UpdateMessage(msg *Message) error {
	fmt.Println("updating message")
	return nil
}

func (pg *PostgresMessageStore) DeleteMessage(id int64) error {
	fmt.Println("deleting message") // not actually deleted, just marked as deleted
	return nil
}

func (pg *PostgresMessageStore) GetUnreadCount(chatID, userID int64) (int64, error) {
	fmt.Println("getting unread chat count")
	return 0, nil
}