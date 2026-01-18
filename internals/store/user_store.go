package store

import (
	"database/sql"
	"fmt"
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
	fmt.Println("getting user by their id")
	return nil, nil
}

func (pg *PostgresUserStore) GetUserByEmail(email string) (*User, error) {
	fmt.Println("getting user by their email")
	return nil, nil
}

func (pg *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	fmt.Println("getting user by their username")
	return nil, nil
}

func (pg *PostgresUserStore) UpdateLastSeen(userID int64) error {
	fmt.Println("updating last seen for user with id: ", userID)
	return nil
}