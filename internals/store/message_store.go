package store

import "time"

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

type MessageStore struct {
	MessageID int64 `json:"message_id"`
	ChatID int64 `json:"chat_id"`
	SenderID int64 `json:"sender_id"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Attachments MessageAttachments `json:"attachments"`
}

type MessageAttachments struct {
	ID int64 `json:"id"`
	MessageID int64 `json:"message_id"`
	Type AttachmentType `json:"type"`
	URL string `json:"url"`
	MetaData *string `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
}