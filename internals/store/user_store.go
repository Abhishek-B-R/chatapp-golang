package store

import (
	"database/sql"
	"time"
)

type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
	AvatarURL string `json:"avatar_url"`
	Bio string `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
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
	CreateUser(user *User) error
    GetUserByID(id int64) (*User, error)
    GetUserByEmail(email string) (*User, error)
    GetUserByUsername(username string) (*User, error)
    UpdateLastSeen(userID int64) error
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash, avatar_url, bio)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at
	`

	err := pg.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.AvatarURL, user.Bio).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresUserStore) GetUserByID(id int64) (*User, error) {
	var user User;
	query := `
		SELECT username, email, avatar_url, bio, created_at FROM users
		WHERE id = $1
	`

	err := pg.db.QueryRow(query, id).Scan(&user.Username, &user.Email, &user.AvatarURL, &user.Bio, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pg *PostgresUserStore) GetUserByEmail(email string) (*User, error) {
	var user User;
	query := `
		SELECT username, email, avatar_url, bio, created_at FROM users
		WHERE email = $1
	`

	err := pg.db.QueryRow(query, email).Scan(&user.Username, &user.Email, &user.AvatarURL, &user.Bio, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pg *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	var user User;
	query := `
		SELECT username, email, avatar_url, bio, created_at FROM users
		WHERE username = $1
	`

	err := pg.db.QueryRow(query, username).Scan(&user.Username, &user.Email, &user.AvatarURL, &user.Bio, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pg *PostgresUserStore) UpdateLastSeen(userID int64) error {
	query := `
		UPDATE users
		SET last_seen_at = $1
		WHERE id = $2
		AND (last_seen_at IS NULL OR last_seen_at < $1);
	`

	_, err := pg.db.Exec(query, time.Now().UTC(), userID)
	if err != nil {
		return err
	}
	return nil
}