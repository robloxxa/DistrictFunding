package auth

import (
	"context"
	"fmt"
	"github.com/robloxxa/DistrictFunding/pkg/db"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// FindByUsernameOrEmailQuery string = `SELECT (id, username, email, first_name, last_name, password) FROM account WHERE username = $1 OR email = $1`

	HasUsernameSql = `SELECT EXISTS(SELECT 1 FROM account WHERE username = $1)`
)

type Account struct {
	Id        string    `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AccountModel interface {
	GetByUUID(string) (*Account, error)
	GetByUsername(string) (*Account, error)
	HasUsername(string) error
	Create(*Account) error
	FindByUsernameOrEmail(string) (*Account, error)

	//Truncate() error
}

type accountModel struct {
	db *pgxpool.Pool
}

func (u *accountModel) GetByUUID(uuid string) (*Account, error) {
	query := `SELECT * FROM account WHERE id = $1`

	return db.QueryOneRowToAddrStruct[Account](context.Background(), u.db, query, uuid)
}

func (u *accountModel) GetByUsername(username string) (*Account, error) {
	query := `SELECT * FROM account WHERE username = $1`

	return db.QueryOneRowToAddrStruct[Account](context.Background(), u.db, query, username)
}

func (u *accountModel) HasUsername(username string) error {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM account WHERE username = $1)`

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
	query := `SELECT * FROM account WHERE username = $1 OR email = $1`

	return db.QueryOneRowToAddrStruct[Account](context.Background(), u.db, query, usernameOrEmail)
}

func (u *accountModel) Create(account *Account) error {
	sql :=
		`INSERT INTO account (username, email, first_name, last_name, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	if _, err := u.db.Exec(context.Background(), sql, account.Username, account.Email, account.FirstName, account.LastName, account.Password); err != nil {
		return err
	}
	return nil
}

//func (u *accountModel) Truncate() error {
//	_, err := u.db.Exec(context.Background(), `TRUNCATE TABLE "user"`)
//	return err
//}
