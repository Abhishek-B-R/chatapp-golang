package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword),12)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plainTextPassword string) (bool,error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
				return false, nil
			default:
				return false, err
		}
	}
	return true, nil
}

type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash password `json:"-"`
	AvatarURL string `json:"avatar_url"`
	Bio string `json:"bio"`
	LastSeenAt time.Time `json:"last_seen_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore{
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(ctx context.Context, user *User) error
    GetUserByID(ctx context.Context, id int64) (*User, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    GetUserByUsername(ctx context.Context, username string) (*User, error)
    UpdateLastSeen(ctx context.Context, userID int64) error
	UpdateUser(ctx context.Context, user *User) error
	UpdateUserPassword(ctx context.Context, password string,userID int64) error
	GetUserToken(ctx context.Context, plainTextPassword string) (*User, error) 
	GetCurrentUser(ctx context.Context, userID int64) (*User, error)
}

func (pg *PostgresUserStore) CreateUser(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash, avatar_url, bio)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at
	`

	err := pg.db.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash.hash, user.AvatarURL, user.Bio).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresUserStore) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User;
	query := `
		SELECT username, email, avatar_url, bio, created_at FROM users
		WHERE id = $1
	`

	err := pg.db.QueryRowContext(ctx, query, id).Scan(&user.Username, &user.Email, &user.AvatarURL, &user.Bio, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pg *PostgresUserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User;
	query := `
		SELECT username, email, avatar_url, bio, created_at FROM users
		WHERE email = $1
	`

	err := pg.db.QueryRowContext(ctx, query, email).Scan(&user.Username, &user.Email, &user.AvatarURL, &user.Bio, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pg *PostgresUserStore) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User;
	query := `
		SELECT id, username, email, password_hash, avatar_url, bio, created_at FROM users
		WHERE username = $1
	`

	err := pg.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.AvatarURL, &user.Bio, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pg *PostgresUserStore) UpdateLastSeen(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET last_seen_at = $1
		WHERE id = $2
		AND (last_seen_at IS NULL OR last_seen_at < $1);
	`

	_, err := pg.db.ExecContext(ctx, query, time.Now().UTC(), userID)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresUserStore) UpdateUser(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, avatar_url = $3, bio = $4, last_seen_at = NOW(), updated_at = NOW()
		WHERE id = $5;
	`

	_, err := pg.db.ExecContext(ctx, query, user.Username, user.Email, user.AvatarURL, user.Bio, user.ID)
	return err
}

func (pg *PostgresUserStore) UpdateUserPassword(ctx context.Context, password string, userID int64) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	q1 := `
		UPDATE users
		SET password_hash = $1
		WHERE id = $2;
	`

	q2 := `
		DELETE FROM tokens
		WHERE user_id = $1
	`
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password),12)
	if err != nil {
		return err
	}
	
	res, err := tx.ExecContext(ctx, q1, string(hash), userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.ExecContext(ctx, q2, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (pg *PostgresUserStore) GetUserToken(ctx context.Context, plainTextPassword string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(plainTextPassword))

	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.avatar_url, u.bio, u.last_seen_at, u.created_at, u.updated_at
		FROM users u
		INNER JOIN tokens t ON t.user_id = u.id
		WHERE t.token_hash = $1 AND t.expires_at > $2
	`

	user := &User{
		PasswordHash: password{},
	}

	err := pg.db.QueryRowContext(ctx, query, tokenHash[:], time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.AvatarURL,
		&user.Bio,
		&user.LastSeenAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil 
}

func (pg *PostgresUserStore) GetCurrentUser(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT username, email, avatar_url, bio, created_at FROM users
		WHERE id = $1
	`

	var user User
	err := pg.db.QueryRowContext(ctx, query, userID).Scan(&user.Username, &user.Email, &user.AvatarURL, &user.Bio, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, err
}