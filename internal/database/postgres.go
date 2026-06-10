package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresDB(host, port, user, password, dbname, sslmode string) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
