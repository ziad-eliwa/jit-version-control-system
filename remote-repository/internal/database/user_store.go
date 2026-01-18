package database

import (
	"database/sql"
	"log/slog"
)

type User struct {
	Username     string `json:"username"`
	PasswordHash []byte `json:"omit"`
	FullName     string `json:"full_name"`
	Bio          string `json:"bio,omitempty"`
	EmailAddress string `json:"email"`
}

type UserProfile struct {
	Username      string   `json:"username"`               // Both
	FullName      string   `json:"full_name"`              // Both
	Bio           string   `json:"bio,omitempty"`          // Both
	EmailAddress  string   `json:"email"`                  // Both
	Repos         []string `json:"repositories,omitempty"` // All repos if self, only public if visitor
	ReposCount    int      `json:"repo_count,omitempty"`   // Self
	TotalCommits  int      `json:"commit_count,omitempty"` // Self
	TopRepository string   `json:"top_repo,omitempty"`     // Self
}

type UserStore interface {
	CreateUser(user *User) (*User, error)
	GetUserbyUsername(username string) (*User, error)
	GetUserbyEmailAddress(email string) (*User, error)
	DeleteUser(username string) error
	GetUserByToken(token string) (string, error)

	GetUserSelfProfile(username string) (*UserProfile, error)
	GetUserProfile(username string) (*UserProfile, error)
}

type PostgresUserStore struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func (pg *PostgresUserStore) DeleteUser(username string) error {
	tx, err := pg.DB.Begin()

	if err != nil {
		return err
	}

	query := `DELETE FROM Users WHERE username = $1`

	_, err = tx.Exec(query, username)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresUserStore) CreateUser(user *User) (*User, error) {
	tx, err := pg.DB.Begin()

	if err != nil {
		return nil, err
	}

	query :=
		`INSERT INTO Users
	(username,fullname, password_hash, bio, email_address)
	VALUES ($1,$2,$3,$4,$5);`

	_, err = tx.Exec(query, user.Username, user.FullName, user.PasswordHash, user.Bio, user.EmailAddress)

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

func (pg *PostgresUserStore) GetUserbyUsername(username string) (*User, error) {
	user := &User{}

	query :=
		`SELECT * FROM Users WHERE username = $1`

	err := pg.DB.QueryRow(query, username).Scan(&user.Username, &user.FullName, &user.PasswordHash, &user.Bio, &user.EmailAddress)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pg *PostgresUserStore) GetUserbyEmailAddress(email string) (*User, error) {
	user := &User{}

	query :=
		`SELECT * FROM Users WHERE email = $1`

	err := pg.DB.QueryRow(query, email).Scan(&user.Username, &user.FullName, &user.PasswordHash, &user.Bio, &user.EmailAddress)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pg *PostgresUserStore) GetUserByToken(token string) (string, error) {
	query :=
		`SELECT username FROM RefreshTokens WHERE refreshtoken = $1`
	var username string
	err := pg.DB.QueryRow(query, token).Scan(&username)

	if err != nil {
		return "", err
	}

	return username, nil
}

func (pg *PostgresUserStore) GetUserSelfProfile(username string) (*UserProfile, error) {
	return nil,nil
}

func (pg *PostgresUserStore) GetUserProfile(username string) (*UserProfile, error) {
	return nil,nil
}
