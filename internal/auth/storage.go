package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id        string    `database:"id"`
	Username  string    `database:"username"`
	Email     string    `database:"email"`
	FirstName string    `database:"first_name"`
	LastName  string    `database:"last_name"`
	Password  string    `database:"password"`
	CreatedAt time.Time `database:"created_at"`
	UpdatedAt time.Time `database:"updated_at"`
}

type UserModel interface {
	GetByUUID(string) (*User, error)
	GetByUsername(string) (*User, error)
	HasUsername(string) error
	Insert(*User) error

	FindByUsernameOrEmail(string) (*User, error)

	//Truncate() error
}

// TODO: Use some scan tools to automatically scan '*' instead of writing all fields names by hand
type userModel struct {
	db *pgxpool.Pool
}

func (u *userModel) GetByUUID(uuid string) (*User, error) {
	var user User

	query :=
		`SELECT (id, username, email, first_name, last_name, password) FROM "user" WHERE id = $1`

	if err := u.db.QueryRow(context.Background(), query, uuid).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userModel) GetByUsername(username string) (*User, error) {
	var user User

	query :=
		`SELECT * FROM "user" WHERE username = $1`

	if err := u.db.QueryRow(context.Background(), query, username).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userModel) HasUsername(username string) error {
	var exists bool

	query :=
		`SELECT EXISTS(SELECT 1 FROM "user" WHERE username = $1)`

	if err := u.db.QueryRow(context.Background(), query, pgx.QueryResultFormats{pgx.TextFormatCode}, username).Scan(&exists); err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("username already exists: %s", username)
	} else {
		return nil
	}
}

func (u *userModel) FindByUsernameOrEmail(usernameOrEmail string) (*User, error) {
	var user User

	query :=
		`SELECT (id, username, email, first_name, last_name, password) FROM "user" WHERE username = $1 OR email = $1`

	if err := u.db.QueryRow(context.Background(), query, usernameOrEmail).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userModel) Insert(user *User) error {
	query :=
		`INSERT INTO "user" (username, email, first_name, last_name, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	if err := u.db.QueryRow(context.Background(), query, user.Username, user.Email, user.FirstName, user.LastName, user.Password).Scan(&user.Id); err != nil {
		return err
	}
	return nil
}

//func (u *userModel) Truncate() error {
//	_, err := u.db.Exec(context.Background(), `TRUNCATE TABLE "user"`)
//	return err
//}
