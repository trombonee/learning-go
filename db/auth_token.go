package db

import "fmt"

type AuthTokenDB struct {
	db *Database
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
	CreatedAt    string
	Email        string
}

func (t AuthTokenDB) InsertAuthToken(token *AuthToken) error {
	if t.db == nil {
		return fmt.Errorf("no reference to db on this auth_token instance")
	}

	_, err := t.db.Exec("INSERT INTO auth_tokens (access_token, refresh_token, created_at, email) VALUES (?, ?, ?, ?)", &token.AccessToken, &token.RefreshToken, &token.CreatedAt, &token.Email)

	return err
}

func (t AuthTokenDB) FetchAuthToken(email string) (AuthToken, error) {
	var token AuthToken

	if t.db == nil {
		return token, fmt.Errorf("no reference to db on this auth_token instance")
	}

	err := t.db.QueryRow("SELECT access_token, refresh_token, created_at, email FROM auth_tokens WHERE email = ?", email).Scan(&token.AccessToken, &token.RefreshToken, &token.CreatedAt, &token.Email)

	return token, err
}

func NewAuthTokenDB(db *Database) *AuthTokenDB {
	return &AuthTokenDB{db: db}
}
