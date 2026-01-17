-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS message_attachments (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT REFERENCES messages(id) ON DELETE CASCADE NOT NULL,
    type VARCHAR(20) NOT NULL,
    url TEXT NOT NULL,
    filename TEXT,
    size_bytes BIGINT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraint: type must be valid
    CONSTRAINT valid_attachment_type CHECK (type IN ('image', 'video', 'pdf', 'file'))
);

-- Index for getting all attachments for a message
CREATE INDEX idx_attachments_message_id ON message_attachments(message_id);

-- Index for filtering by attachment type
CREATE INDEX idx_attachments_type ON message_attachments(type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_attachments_type;
DROP INDEX IF EXISTS idx_attachments_message_id;
DROP TABLE IF EXISTS message_attachments;
-- +goose StatementEnd