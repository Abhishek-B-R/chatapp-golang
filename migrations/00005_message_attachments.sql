-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS message_attachments (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,
    url TEXT NOT NULL,
    filename TEXT,
    size_bytes BIGINT,
    metadata JSONB DEFAULT '{}'::jsonb NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Constraint: type must be valid
    CONSTRAINT valid_attachment_type CHECK (type IN ('image', 'video', 'pdf', 'file'))
);

-- Index for getting all attachments for a message
CREATE INDEX idx_attachments_message_id ON message_attachments(message_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_attachments_type;
DROP INDEX IF EXISTS idx_attachments_message_id;
DROP TABLE IF EXISTS message_attachments;
-- +goose StatementEnd