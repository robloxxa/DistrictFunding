package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// FindByUsernameOrEmailQuery string = `SELECT (id, username, email, first_name, last_name, password) FROM account WHERE username = $1 OR email = $1`

	HasUsernameSql = `SELECT EXISTS(SELECT 1 FROM account WHERE username = $1)`
)

type Account struct {
	Id        string    `database:"id"`
	Username  string    `database:"username"`
	Email     string    `database:"email"`
	FirstName string    `database:"first_name"`
	LastName  string    `database:"last_name"`
	Password  string    `database:"password"`
	CreatedAt time.Time `database:"created_at"`
	UpdatedAt time.Time `database:"updated_at"`
}

type AccountModel interface {
	GetByUUID(string) (*Account, error)
	GetByUsername(string) (*Account, error)
	HasUsername(string) error
	Create(*Account) error
	FindByUsernameOrEmail(string) (*Account, error)

	//Truncate() error
}

// TODO: Use some scan tools to automatically scan '*' instead of writing all fields names by hand
type accountModel struct {
	db *pgxpool.Pool
}

func (u *accountModel) GetByUUID(uuid string) (*Account, error) {
	var user Account

	query :=
		`SELECT (id, username, email, first_name, last_name, password) FROM account WHERE id = $1`

	if err := u.db.QueryRow(context.Background(), query, uuid).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *accountModel) GetByUsername(username string) (*Account, error) {
	var user Account

	query :=
		`SELECT * FROM account WHERE username = $1`

	if err := u.db.QueryRow(context.Background(), query, username).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *accountModel) HasUsername(username string) error {
	var exists bool

	query :=
		`SELECT EXISTS(SELECT 1 FROM account WHERE username = $1)`

	if err := u.db.QueryRow(context.Background(), query, pgx.QueryResultFormats{pgx.TextFormatCode}, username).Scan(&exists); err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("username already exists: %s", username)
	} else {
		return nil
	}
}

func (u *accountModel) FindByUsernameOrEmail(usernameOrEmail string) (*Account, error) {
	query :=
		`SELECT (id, username, email, first_name, last_name, password) FROM account WHERE username = $1 OR email = $1`
	return QueryOneRowToAddrStruct[Account](context.Background(), u.db, query, usernameOrEmail)
}

func (u *accountModel) Create(account *Account) error {
	query :=
		`INSERT INTO account (username, email, first_name, last_name, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	if _, err := u.db.Exec(context.Background(), query, account.Username, account.Email, account.FirstName, account.LastName, account.Password); err != nil {
		return err
	}
	return nil
}

func QueryOneRowToAddrStruct[T any](ctx context.Context, db *pgxpool.Pool, query string, arguments ...any) (*T, error) {
	rows, err := db.Query(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}

	t, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[T])
	if err != nil {
		return nil, err
	}

	return t, nil
}

//func (u *accountModel) Truncate() error {
//	_, err := u.db.Exec(context.Background(), `TRUNCATE TABLE "user"`)
//	return err
//}
