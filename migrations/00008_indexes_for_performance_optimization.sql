-- +goose Up
-- +goose StatementBegin
-- Composite index for finding user's chats sorted by last activity
CREATE INDEX idx_chat_members_user_last_message 
ON chat_members(user_id, chat_id);

-- Partial index for active (non-deleted) messages only
CREATE INDEX idx_messages_active 
ON messages(chat_id, created_at DESC) 
WHERE deleted_at IS NULL;

-- Index for finding DM chats between two users
-- This helps with GetOrCreateDM queries
CREATE INDEX idx_dm_chats 
ON chat_members(chat_id, user_id) 
WHERE (SELECT is_group FROM chats WHERE id = chat_id) = false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_dm_chats;
DROP INDEX IF EXISTS idx_messages_active;
DROP INDEX IF EXISTS idx_chat_members_user_last_message;
-- +goose StatementEnd