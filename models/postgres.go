package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// will open a sql connection to the database
// caller is responsible for closing the db connection
func NewPostgresDB(cfg PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.String())

	if err != nil {
		return nil, fmt.Errorf("models.NewPostgresDB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("models.NewPostgresDB: %w", err)
	}

	return db, nil
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "lenslocked",
		Password: "lenslocked123",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}
