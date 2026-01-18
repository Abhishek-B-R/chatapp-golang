-- +goose Up
-- +goose StatementBegin
-- Add foreign key constraint for last_read_message_id
-- This is separate because messages table must exist first
ALTER TABLE chat_members 
ADD CONSTRAINT fk_last_read_message 
FOREIGN KEY (last_read_message_id) 
REFERENCES messages(id) 
ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chat_members 
DROP CONSTRAINT IF EXISTS fk_last_read_message;
-- +goose StatementEnd