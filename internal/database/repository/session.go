package repository

import (
	"database/sql"
	"fmt"

	"telemonitor/internal/database"
)

// SessionRepository handles session_storage operations
type SessionRepository struct {
	db *database.DB
}

// NewSessionRepository creates a new SessionRepository
func NewSessionRepository(db *database.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Save stores or updates a session
func (r *SessionRepository) Save(key string, value []byte) error {
	query := `
		INSERT INTO session_storage (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value
	`
	_, err := r.db.Exec(query, key, value)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

// Load retrieves a session by key
func (r *SessionRepository) Load(key string) ([]byte, error) {
	var value []byte
	query := `SELECT value FROM session_storage WHERE key = $1`
	
	err := r.db.QueryRow(query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return nil, nil // No session found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load session: %w", err)
	}
	
	return value, nil
}

// Delete removes a session
func (r *SessionRepository) Delete(key string) error {
	query := `DELETE FROM session_storage WHERE key = $1`
	_, err := r.db.Exec(query, key)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
