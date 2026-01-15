package database

import (
	"database/sql"
	"fmt"
	"time"
)

type TokenStore interface {
	StoreRefreshToken(username, token string) error
	GetRefreshToken(username string) (*RefreshToken, error)
	RevokeToken(username string) error
}

type RefreshToken struct {
	Token     string    `json:"token"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
	RevokedAt time.Time `json:"revoked_at,omitempty"`
}

type PostgresTokenStore struct {
	DB *sql.DB
}

func (pg *PostgresTokenStore) StoreRefreshToken(username, token string) error {
	query :=
		`INSERT INTO RefreshTokens (username,refreshtoken,created_at,revoked) VALUES ($1,$2,$3,$4)`

	tx, err := pg.DB.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec(query, username, token, time.Now(), 0)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return tx.Rollback()
	}
	return nil
}

func (pg *PostgresTokenStore) GetRefreshToken(username, token string) (*RefreshToken, error) {
	query :=
		`SELECT refreshtoken, created_at,revoked, revoked_at FROM RefreshTokens WHERE username = $1 AND refreshtoken = $2`
	refreshtoken := &RefreshToken{}
	err := pg.DB.QueryRow(query, username, token).Scan(&refreshtoken)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Invalid Refresh Token")
		}
		return nil, err
	}

	return refreshtoken, nil
}

func (pg *PostgresTokenStore) RevokeAllTokens(username string) error {
	query := 
	`UPDATE RefreshTokens WHERE username = $1 SET revoked = true`

	tx,err := pg.DB.Begin()

	if err != nil {
		return err
	}

	_,err = tx.Exec(query,username)

	if err != nil {
		return err
	}

	return nil 
}
