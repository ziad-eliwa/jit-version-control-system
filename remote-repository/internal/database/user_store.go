package database

import (
	"database/sql"
	"log/slog"

	"github.com/ziad-eliwa/jit-version-control-system/internal/pkg/hashing"
)

type User struct {
	Username     string           `json:"username"`
	Password     hashing.Password `json:"omit"`
	FullName     string           `json:"full_name"`
	Bio          string           `json:"bio"`
	EmailAddress string           `json:"Email Address"`
}

type UserStore interface {
	CreateUser(user *User) (*User,error)
	GetUserbyUsername(username string) (*User, error)  
	GetUserbyEmailAddress(email string) (*User,error)
}

type PostgresUserStore struct {
	DB *sql.DB
	Logger *slog.Logger
}

func (pg *PostgresUserStore) CreateUser(user *User) (*User, error) {
	tx,err := pg.DB.Begin()

	if err != nil {
		return nil, err
	}

	query := 
	`INSERT INTO Users
	(username,fullname, password_hash, bio, email_address)
	VALUES ($1,$2,$3,$4,$5);`

	_, err = tx.Exec(query,user.Username,user.FullName,user.Password.Hash,user.Bio,user.EmailAddress)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return user, nil 		
}

func (pg *PostgresUserStore) GetUserbyUsername(username string) (*User,error) {
	user := &User{}

	query := 
	`SELECT * FROM Users WHERE username = $1`

	err := pg.DB.QueryRow(query,username).Scan(&user)

	if err != nil {
		return nil, err
	}

	return user,nil
}

func (pg *PostgresUserStore) GetUserbyEmailAddress(email string) (*User,error) {
	user := &User{}

	query := 
	`SELECT * FROM Users WHERE email = $1`

	err := pg.DB.QueryRow(query,email).Scan(&user)

	if err != nil {
		return nil, err
	}

	return user,nil
}
