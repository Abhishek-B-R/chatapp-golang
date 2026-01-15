package store

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	UserID int64 `json:"user_id"`
	Name string `json:"name"`
	Email string `json:"email"`
	AvatarURL string `json:"avatar_url"`
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
	fmt.Println("user created")
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