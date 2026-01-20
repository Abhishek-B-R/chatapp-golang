package store

import (
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
	UpdateUser(user *User) error
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash, avatar_url, bio)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at
	`

	err := pg.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.AvatarURL, user.Bio).Scan(&user.ID, &user.CreatedAt)
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

func (pg *PostgresUserStore) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, avatar_url = $3, bio = $4, last_seen_at = NOW()
		WHERE id = $5;
	`

	_, err := pg.db.Exec(query, user.Username, user.Email, user.AvatarURL, user.Bio, user.ID)
	return err
}

func (pg *PostgresUserStore) UpdateUserPassword(password string) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, avatar_url = $3, bio = $4, last_seen_at = NOW()
		WHERE id = $5;
	`

	
	hash, err := bcrypt.GenerateFromPassword([]byte(password),12)
	if err != nil {
		return err
	}
	
	_, err = pg.db.Exec(query, hash)
	return err
}