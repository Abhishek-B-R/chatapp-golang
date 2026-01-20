-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT REFERENCES chats(id) ON DELETE CASCADE NOT NULL,
    sender_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    
    content TEXT,
    type VARCHAR(20) DEFAULT 'text' NOT NULL,
    
    reply_to_message_id BIGINT REFERENCES messages(id) ON DELETE SET NULL,
    
    edited_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraint: type must be valid
    CONSTRAINT valid_message_type CHECK (type IN ('text', 'system'))
);

-- Most important index: get messages for a chat, sorted by time
CREATE INDEX idx_messages_chat_created ON messages(chat_id, created_at DESC);

-- Index for finding replies to a message
CREATE INDEX idx_messages_reply_to ON messages(reply_to_message_id) 
    WHERE reply_to_message_id IS NOT NULL;

-- Index for counting unread messages (messages after a certain ID)
CREATE INDEX idx_messages_chat_id_filter ON messages(chat_id, id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_messages_chat_id_filter;
DROP INDEX IF EXISTS idx_messages_reply_to;
DROP INDEX IF EXISTS idx_messages_chat_created;
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd