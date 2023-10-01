package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/trenchesdeveloper/lenslocked/rand"
)

const (
	// The minimum number of bytes we want for our session ID
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when the session is created
	// This will be sent to the user in a cookie
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	row := ss.DB.QueryRow(`
	INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2)
	ON CONFLICT (user_id) DO UPDATE SET token_hash = $2
	RETURNING id`,
	session.UserID, session.TokenHash)

	if err := row.Scan(&session.ID); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// hash the token
	tokenHash := ss.hash(token)
	// query the db for the session
	row := ss.DB.QueryRow(`
	SELECT users.id, users.email, users.password_hash
	FROM sessions
	JOIN users ON users.id = sessions.user_id
	WHERE sessions.token_hash = $1`, tokenHash,
)

	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash ); err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)

	_, err := ss.DB.Exec("DELETE FROM sessions WHERE token_hash = $1", tokenHash)

	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))

	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
