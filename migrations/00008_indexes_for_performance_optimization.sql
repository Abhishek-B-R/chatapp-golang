-- +goose Up
-- +goose StatementBegin
-- Partial index for active (non-deleted) messages only
CREATE INDEX idx_messages_active 
ON messages(chat_id, created_at DESC) 
WHERE deleted_at IS NULL;

-- Note: Removed the other two indexes as they're redundant
-- idx_chat_members_user_last_message is covered by existing idx_chat_members_user_id
-- idx_dm_chats can't use subquery in WHERE clause - we'll handle DM lookups differently
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_messages_active;
-- +goose StatementEnd