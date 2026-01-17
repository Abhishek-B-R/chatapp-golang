-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat_members (
    chat_id BIGINT REFERENCES chats(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member' NOT NULL,
    last_read_message_id BIGINT,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    muted BOOLEAN DEFAULT false,
    pinned BOOLEAN DEFAULT false,
    
    PRIMARY KEY (chat_id, user_id),
    
    -- Constraint: role must be one of: owner, admin, member
    CONSTRAINT valid_role CHECK (role IN ('owner', 'admin', 'member'))
);

-- Index for getting all chats for a user (most common query)
CREATE INDEX idx_chat_members_user_id ON chat_members(user_id);

-- Index for getting all members of a chat
CREATE INDEX idx_chat_members_chat_id ON chat_members(chat_id);

-- Index for finding chat owner quickly
CREATE INDEX idx_chat_members_owner ON chat_members(chat_id, role) WHERE role = 'owner';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chat_members_owner;
DROP INDEX IF EXISTS idx_chat_members_chat_id;
DROP INDEX IF EXISTS idx_chat_members_user_id;
DROP TABLE IF EXISTS chat_members;
-- +goose StatementEnd