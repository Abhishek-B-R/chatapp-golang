package store

import (
	"database/sql"
	"time"

	"github.com/Abhishek-B-R/chat-app-golang/internals/tokens"
)

type Token struct {
	ID        int64		`json:"id"`
	UserID   int64		`json:"user_id"`
	TokenHash string	`json:"token_hash"`
	ExpiresAt time.Time	`json:"expires_at"`
	CreatedAt time.Time	`json:"created_at"`
}

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db:db}
}

type TokenStore interface {
    Insert(token *tokens.Token) error
    CreateNewToken(
        userId int64,
        ttl time.Duration,
    ) (*tokens.Token, error)
    DeleteAllTokensForUser(userID int) error
}

func (pg *PostgresTokenStore) CreateNewToken(userId int64, ttl time.Duration) (*tokens.Token, error) {
	tok, err := tokens.GenerateToken(userId, ttl)
	    if err != nil {
        return nil, err
    }

    err = pg.Insert(tok)
    return tok, err
}

func (pg *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
		INSERT INTO tokens (user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
	`

	_, err := pg.db.Exec(query, token.UserID, token.Hash, token.ExpiresAt)
	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userId int) error {
	query := `
	DELETE FROM tokens
	WHERE user_id = $1
	`

	_, err := t.db.Exec(query, userId)
	return err
}