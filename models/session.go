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
		UserID: userID,
		Token:  token,
		TokenHash: ss.hash(token),
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	var user User

	row := ss.DB.QueryRow("SELECT users.id, users.email FROM users JOIN sessions ON users.id = sessions.user_id WHERE sessions.token_hash = $1", token)

	if err := row.Scan(&user.ID, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))

	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
