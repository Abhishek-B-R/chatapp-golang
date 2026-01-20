package store

import (
	"database/sql"
	"encoding/json"
	"errors"
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
    GetChatMessages(chatID, limit, offset int64) (*[]Message, error)
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

func (pg *PostgresMessageStore) GetMessage(msgID int64) (*Message, error) {
	var msg Message;
	q1 := `
		SELECT
			id,
			chat_id,
			sender_id,
			type,
			content,
			reply_to_message_id,
			created_at,
			edited_at,
			deleted_at
		FROM messages
		WHERE id = $1
	`
	err := pg.db.QueryRow(q1, msgID).Scan(
		&msg.ID, 
		&msg.ChatID, 
		&msg.SenderID, 
		&msg.Type, 
		&msg.Content, 
		&msg.ReplyToMessageID, 
		&msg.CreatedAt, 
		&msg.EditedAt,
		&msg.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	attachments, err := pg.getAttachmentsForMessages([]int64{msg.ID})
	if err != nil {
		return nil, err
	}

	msg.Attachments = attachments
	return &msg, nil
}

func (pg *PostgresMessageStore) GetChatMessages(chatID, limit, offset int64) (*[]Message, error) {
	q1 := `
		SELECT
			id,
			chat_id,
			sender_id,
			type,
			content,
			reply_to_message_id,
			created_at,
			edited_at,
			deleted_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := pg.db.Query(q1, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgIDs []int64
	var msgs []Message

	for rows.Next() {
		var m Message
		err := rows.Scan(
			&m.ID,
			&m.ChatID,
			&m.SenderID,
			&m.Type,
			&m.Content,
			&m.ReplyToMessageID,
			&m.CreatedAt,
			&m.EditedAt,
			&m.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, m)
		msgIDs = append(msgIDs, m.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(msgIDs) == 0 {
		return &msgs, nil
	}

	//fetch attachments
	attachments, err := pg.getAttachmentsForMessages(msgIDs)
	if err != nil {
		return nil, err
	}

	// stich them with messages
	attMap := make(map[int64][]MessageAttachment)
	for _, a := range attachments {
		attMap[a.MessageID] = append(attMap[a.MessageID], a)
	}

	for i := range msgs {
		msgs[i].Attachments = attMap[msgs[i].ID]

		if msgs[i].DeletedAt != nil {
			msgs[i].Content = nil
			msgs[i].Type = "deleted"
			msgs[i].Attachments = nil
		}
	}
	
	return &msgs, nil
}

func (pg *PostgresMessageStore) UpdateMessage(msg *Message) error {
	query := `
		UPDATE messages
        SET content = $1, edited_at = NOW()
        WHERE id = $2
          AND deleted_at IS NULL
        RETURNING id
	`

	err := pg.db.QueryRow(query, msg.Content, msg.ID).Scan(&msg.ID)
	if err == sql.ErrNoRows {	
        return errors.New("message not found or already deleted")
    }

	return err
}

func (pg *PostgresMessageStore) DeleteMessage(msgID int64) error {
	query := `
		UPDATE messages
		SET deleted_at = NOW()
		WHERE id = $2
	` // not actually deleted, just marked as deleted

	_, err := pg.db.Exec(query, msgID)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresMessageStore) GetUnreadCount(chatID, userID int64) (int64, error) {
	fmt.Println("getting unread chat count")
	return 0, nil
}

func (pg *PostgresMessageStore) getAttachmentsForMessages(messageIDs []int64) ([]MessageAttachment, error) {
	const q = `
		SELECT
			id,
			message_id,
			type,
			url,
			filename,
			size_bytes,
			metadata,
			created_at
		FROM message_attachments
		WHERE message_id = ANY($1)
		ORDER BY created_at ASC;
`

	rows, err := pg.db.Query(q, messageIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var atts []MessageAttachment
	for rows.Next() {
		var a MessageAttachment
		err := rows.Scan(
			&a.ID,
			&a.MessageID,
			&a.Type,
			&a.URL,
			&a.Filename,
			&a.SizeBytes,
			&a.Metadata,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		atts = append(atts, a)
	}

	return atts, rows.Err()
}
