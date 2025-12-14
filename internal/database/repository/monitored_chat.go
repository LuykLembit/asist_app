package repository

import (
	"database/sql"
	"fmt"
	"time"

	"telemonitor/internal/database"
)

// MonitoredChatRepository handles monitored_chats operations
type MonitoredChatRepository struct {
	db *database.DB
}

// NewMonitoredChatRepository creates a new MonitoredChatRepository
func NewMonitoredChatRepository(db *database.DB) *MonitoredChatRepository {
	return &MonitoredChatRepository{db: db}
}

// Create adds a new monitored chat
func (r *MonitoredChatRepository) Create(chat *database.MonitoredChat) error {
	query := `
		INSERT INTO monitored_chats (chat_id, title, username, is_active, added_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (chat_id) DO NOTHING
	`
	
	if chat.AddedAt.IsZero() {
		chat.AddedAt = time.Now()
	}
	
	_, err := r.db.Exec(query, chat.ChatID, chat.Title, chat.Username, chat.IsActive, chat.AddedAt)
	if err != nil {
		return fmt.Errorf("failed to create monitored chat: %w", err)
	}
	return nil
}

// GetByChatID retrieves a chat by ID
func (r *MonitoredChatRepository) GetByChatID(chatID int64) (*database.MonitoredChat, error) {
	query := `
		SELECT chat_id, title, username, last_processed_msg_id, last_pts, is_active, added_at
		FROM monitored_chats
		WHERE chat_id = $1
	`
	
	chat := &database.MonitoredChat{}
	err := r.db.QueryRow(query, chatID).Scan(
		&chat.ChatID,
		&chat.Title,
		&chat.Username,
		&chat.LastProcessedMsgID,
		&chat.LastPts,
		&chat.IsActive,
		&chat.AddedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get monitored chat: %w", err)
	}
	
	return chat, nil
}

// GetAll retrieves all monitored chats
func (r *MonitoredChatRepository) GetAll() ([]*database.MonitoredChat, error) {
	query := `
		SELECT chat_id, title, username, last_processed_msg_id, last_pts, is_active, added_at
		FROM monitored_chats
		ORDER BY added_at DESC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all monitored chats: %w", err)
	}
	defer rows.Close()
	
	var chats []*database.MonitoredChat
	for rows.Next() {
		chat := &database.MonitoredChat{}
		if err := rows.Scan(
			&chat.ChatID,
			&chat.Title,
			&chat.Username,
			&chat.LastProcessedMsgID,
			&chat.LastPts,
			&chat.IsActive,
			&chat.AddedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan monitored chat: %w", err)
		}
		chats = append(chats, chat)
	}
	
	return chats, nil
}

// GetActive retrieves all active monitored chats
func (r *MonitoredChatRepository) GetActive() ([]*database.MonitoredChat, error) {
	query := `
		SELECT chat_id, title, username, last_processed_msg_id, last_pts, is_active, added_at
		FROM monitored_chats
		WHERE is_active = TRUE
		ORDER BY added_at DESC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active monitored chats: %w", err)
	}
	defer rows.Close()
	
	var chats []*database.MonitoredChat{}
	for rows.Next() {
		chat := &database.MonitoredChat{}
		if err := rows.Scan(
			&chat.ChatID,
			&chat.Title,
			&chat.Username,
			&chat.LastProcessedMsgID,
			&chat.LastPts,
			&chat.IsActive,
			&chat.AddedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan monitored chat: %w", err)
		}
		chats = append(chats, chat)
	}
	
	return chats, nil
}

// Update updates a monitored chat
func (r *MonitoredChatRepository) Update(chat *database.MonitoredChat) error {
	query := `
		UPDATE monitored_chats
		SET title = $2, username = $3, last_processed_msg_id = $4, 
		    last_pts = $5, is_active = $6
		WHERE chat_id = $1
	`
	
	_, err := r.db.Exec(query,
		chat.ChatID,
		chat.Title,
		chat.Username,
		chat.LastProcessedMsgID,
		chat.LastPts,
		chat.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to update monitored chat: %w", err)
	}
	return nil
}

// UpdateState updates only the message tracking state
func (r *MonitoredChatRepository) UpdateState(chatID int64, lastMsgID, lastPts int) error {
	query := `
		UPDATE monitored_chats
		SET last_processed_msg_id = $2, last_pts = $3
		WHERE chat_id = $1
	`
	
	_, err := r.db.Exec(query, chatID, lastMsgID, lastPts)
	if err != nil {
		return fmt.Errorf("failed to update chat state: %w", err)
	}
	return nil
}

// Delete removes a monitored chat
func (r *MonitoredChatRepository) Delete(chatID int64) error {
	query := `DELETE FROM monitored_chats WHERE chat_id = $1`
	_, err := r.db.Exec(query, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete monitored chat: %w", err)
	}
	return nil
}

// SetActive sets the active status of a chat
func (r *MonitoredChatRepository) SetActive(chatID int64, active bool) error {
	query := `UPDATE monitored_chats SET is_active = $2 WHERE chat_id = $1`
	_, err := r.db.Exec(query, chatID, active)
	if err != nil {
		return fmt.Errorf("failed to set chat active status: %w", err)
	}
	return nil
}
