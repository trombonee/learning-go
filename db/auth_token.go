package db

import (
	"fmt"
	"strconv"
	"time"
)

const expirationTime = 3600

type AuthToken struct {
	AccessToken  string
	RefreshToken string
	CreatedAt    string
	Email        string
}

func (t *AuthToken) IsTokenExpired() bool {
	if t.CreatedAt == "" {
		return true
	}

	currTime := time.Now()

	epochInt, err := strconv.ParseInt(t.CreatedAt, 10, 64)
	if err != nil {
		return true
	}

	createdAt := time.Unix(epochInt+expirationTime, 0)

	return currTime.After(createdAt)
}

type AuthTokenDB struct {
	db *Database
}

func (t AuthTokenDB) InsertAuthToken(token *AuthToken) error {
	if t.db == nil {
		return fmt.Errorf("no reference to db on this auth_token instance")
	}

	_, err := t.db.Exec("INSERT INTO auth_tokens (access_token, refresh_token, created_at, email) VALUES (?, ?, ?, ?)", &token.AccessToken, &token.RefreshToken, &token.CreatedAt, &token.Email)

	return err
}

func (t AuthTokenDB) FetchAuthToken(email string) (*AuthToken, error) {
	var token AuthToken

	if t.db == nil {
		return nil, fmt.Errorf("no reference to db on this auth_token instance")
	}

	err := t.db.QueryRow("SELECT access_token, refresh_token, created_at, email FROM auth_tokens WHERE email = ?", email).Scan(&token.AccessToken, &token.RefreshToken, &token.CreatedAt, &token.Email)

	return &token, err
}

func NewAuthTokenDB(db *Database) *AuthTokenDB {
	return &AuthTokenDB{db: db}
}
