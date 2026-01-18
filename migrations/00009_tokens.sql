-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    token_hash BYTEA UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for token lookup
CREATE INDEX idx_tokens_hash ON tokens(token_hash);

-- Index for finding user's tokens
CREATE INDEX idx_tokens_user_id ON tokens(user_id);

-- Index for cleaning up expired tokens
CREATE INDEX idx_tokens_expires_at ON tokens(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tokens_expires_at;
DROP INDEX IF EXISTS idx_tokens_user_id;
DROP INDEX IF EXISTS idx_tokens_hash;
DROP TABLE IF EXISTS tokens;
-- +goose StatementEnd