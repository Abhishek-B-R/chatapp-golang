-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chats (
    id BIGSERIAL PRIMARY KEY,
    is_group BOOLEAN DEFAULT false NOT NULL,
    name TEXT,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    last_message_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraint: group chats must have a name
    CONSTRAINT group_chat_must_have_name CHECK (
        (is_group = false) OR (is_group = true AND name IS NOT NULL)
    )
);

-- Index for faster lookups by creator
CREATE INDEX idx_chats_created_by ON chats(created_by);

-- Index for sorting chats by last message time
CREATE INDEX idx_chats_last_message_at ON chats(last_message_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chats_last_message_at;
DROP INDEX IF EXISTS idx_chats_created_by;
DROP TABLE IF EXISTS chats;
-- +goose StatementEnd