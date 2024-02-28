package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id        string `database:"id"`
	Username  string `database:"username"`
	Email     string `database:"email"`
	FirstName string `database:"first_name"`
	LastName  string `database:"last_name"`
	Password  string `database:"password"`
}

type UserModel interface {
	GetByUUID(string) (*User, error)
	GetByUsername(string) (*User, error)
	HasUsername(string) (bool, error)
	Insert(*User) (bool, error)
}

type userModel struct {
	db *pgxpool.Pool
}

func (u *userModel) GetByUUID(uuid string) (*User, error) {
	var user User

	query := "SELECT * FROM users WHERE id = $1"
	if err := u.db.QueryRow(context.Background(), query, uuid).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userModel) GetByUsername(username string) (*User, error) {
	var user User
	query := "SELECT * FROM users WHERE username = $1"
	if err := u.db.QueryRow(context.Background(), query, username).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userModel) HasUsername(username string) (bool, error) {
	var b bool
	query := "SELECT EXISTS(SELECT true FROM users WHERE username = $1)"

	if err := u.db.QueryRow(context.Background(), query, username).Scan(&b); err != nil {
		return false, err
	}

	return b, nil
}

func (u *userModel) Insert(user *User) (bool, error) {
	query := "INSERT INTO users (username, email, first_name, last_name, password) VALUES ($1, $2, $3, $4, $5)"
	tag, err := u.db.Exec(context.Background(), query, user.Username, user.Email, user.FirstName, user.LastName, user.Password)
	return tag.Insert(), err
}
