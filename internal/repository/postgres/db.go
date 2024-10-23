package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)

type Database struct {
	Db *sql.DB
}

func New(dsn string) (*Database, error) {
	const op = "repository.db.New"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Database{Db: db}, nil
}

func (db *Database) Stop() error {
	return db.Db.Close()
}
