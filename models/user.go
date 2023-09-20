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

	user := User{
		Email:    email,
		Password: passwordHash,
	}

	rows := us.DB.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", email, passwordHash)

	if err := rows.Scan(&user.ID); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)

	// get the user from the db by email
	user := User{
		Email: email,
	}
	row := us.DB.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", email)

	if err := row.Scan(&user.ID, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("authenticate: %w", err)
		}
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	// check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("authenticate: %w", err)
		}
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return &user, nil
}
