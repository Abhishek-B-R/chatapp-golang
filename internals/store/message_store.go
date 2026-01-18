package store

import (
	"database/sql"
	"fmt"
	"time"
)

type AttachmentType string

const (
	JPG  AttachmentType = "jpg"
	PNG  AttachmentType = "png"
	JPEG AttachmentType = "jpeg"
	WEBP AttachmentType = "webp"
	GIF  AttachmentType = "gif"
	VID  AttachmentType = "mp4"
	PDF  AttachmentType = "pdf"
	XLSX AttachmentType = "xlsx"
	TXT  AttachmentType = "txt"
)

type Message struct {
	MessageID int64 `json:"message_id"`
	ChatID int64 `json:"chat_id"`
	SenderID int64 `json:"sender_id"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Attachments []MessageAttachments `json:"attachments"`
}

type MessageAttachments struct {
	ID int64 `json:"id"`
	MessageID int64 `json:"message_id"`
	Type AttachmentType `json:"type"`
	URL string `json:"url"`
	MetaData *string `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
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
	fmt.Println("creating message")
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