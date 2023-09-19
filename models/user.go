package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService struct {
	DB *sql.DB
}


func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	passwordHash := string(hashedBytes)

	user := User {
		Email: email,
		Password: passwordHash,
	}

	rows := us.DB.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", email, passwordHash)

	if err := rows.Scan(&user.ID); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}
