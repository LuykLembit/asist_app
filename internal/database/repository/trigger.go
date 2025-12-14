package repository

import (
	"database/sql"
	"fmt"

	"telemonitor/internal/database"
)

// TriggerRepository handles triggers operations
type TriggerRepository struct {
	db *database.DB
}

// NewTriggerRepository creates a new TriggerRepository
func NewTriggerRepository(db *database.DB) *TriggerRepository {
	return &TriggerRepository{db: db}
}

// Create inserts a new trigger
func (r *TriggerRepository) Create(trigger *database.Trigger) error {
	query := `
		INSERT INTO triggers (phrase, is_regex, alert_level)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	
	err := r.db.QueryRow(query, trigger.Phrase, trigger.IsRegex, trigger.AlertLevel).Scan(&trigger.ID)
	if err != nil {
		return fmt.Errorf("failed to create trigger: %w", err)
	}
	
	return nil
}

// GetByID retrieves a trigger by ID
func (r *TriggerRepository) GetByID(id int) (*database.Trigger, error) {
	query := `SELECT id, phrase, is_regex, alert_level FROM triggers WHERE id = $1`
	
	trigger := &database.Trigger{}
	err := r.db.QueryRow(query, id).Scan(&trigger.ID, &trigger.Phrase, &trigger.IsRegex, &trigger.AlertLevel)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get trigger: %w", err)
	}
	
	return trigger, nil
}

// GetAll retrieves all triggers
func (r *TriggerRepository) GetAll() ([]*database.Trigger, error) {
	query := `SELECT id, phrase, is_regex, alert_level FROM triggers ORDER BY id`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all triggers: %w", err)
	}
	defer rows.Close()
	
	var triggers []*database.Trigger
	for rows.Next() {
		trigger := &database.Trigger{}
		if err := rows.Scan(&trigger.ID, &trigger.Phrase, &trigger.IsRegex, &trigger.AlertLevel); err != nil {
			return nil, fmt.Errorf("failed to scan trigger: %w", err)
		}
		triggers = append(triggers, trigger)
	}
	
	return triggers, nil
}

// Update updates a trigger
func (r *TriggerRepository) Update(trigger *database.Trigger) error {
	query := `
		UPDATE triggers
		SET phrase = $2, is_regex = $3, alert_level = $4
		WHERE id = $1
	`
	
	_, err := r.db.Exec(query, trigger.ID, trigger.Phrase, trigger.IsRegex, trigger.AlertLevel)
	if err != nil {
		return fmt.Errorf("failed to update trigger: %w", err)
	}
	
	return nil
}

// Delete removes a trigger
func (r *TriggerRepository) Delete(id int) error {
	query := `DELETE FROM triggers WHERE id = $1`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trigger: %w", err)
	}
	
	return nil
}
