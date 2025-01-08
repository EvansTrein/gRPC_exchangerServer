package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db  *sql.DB
	log *slog.Logger
}
// database connection
func New(storagePath string, log *slog.Logger) (*SQLiteDB, error) {

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("connect to DB successfully")

	return &SQLiteDB{db: db, log: log}, nil
}

// database disconnection
func (s *SQLiteDB) Close() error {
	if s.db == nil {
		return fmt.Errorf("database connection is already closed")
	}

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	s.db = nil

	return nil
}
