package models

import "database/sql"

type Session struct {
	ID        int
	UserID    int
	// Token is only set when the session is created
	// This will be sent to the user in a cookie
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	var user User

	row := ss.DB.QueryRow("SELECT users.id, users.email FROM users JOIN sessions ON users.id = sessions.user_id WHERE sessions.token_hash = $1", token)

	if err := row.Scan(&user.ID, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}