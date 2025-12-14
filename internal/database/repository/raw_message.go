package repository

import (
	"database/sql"
	"fmt"
	"time"

	"telemonitor/internal/database"
)

// RawMessageRepository handles raw_messages operations
type RawMessageRepository struct {
	db *database.DB
}

// NewRawMessageRepository creates a new RawMessageRepository
func NewRawMessageRepository(db *database.DB) *RawMessageRepository {
	return &RawMessageRepository{db: db}
}

// Create inserts a new raw message
func (r *RawMessageRepository) Create(msg *database.RawMessage) error {
	query := `
		INSERT INTO raw_messages (
			chat_id, telegram_msg_id, sender_id, sender_name, message_text,
			is_transcribed, is_forward, forward_source_name, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (chat_id, telegram_msg_id) DO NOTHING
		RETURNING id
	`
	
	err := r.db.QueryRow(query,
		msg.ChatID,
		msg.TelegramMsgID,
		msg.SenderID,
		msg.SenderName,
		msg.MessageText,
		msg.IsTranscribed,
		msg.IsForward,
		msg.ForwardSourceName,
		msg.CreatedAt,
	).Scan(&msg.ID)
	
	if err == sql.ErrNoRows {
		// Duplicate message, ignore
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to create raw message: %w", err)
	}
	
	return nil
}

// GetByChatIDAndTimeRange retrieves messages for a chat within a time range
func (r *RawMessageRepository) GetByChatIDAndTimeRange(chatID int64, start, end time.Time) ([]*database.RawMessage, error) {
	query := `
		SELECT id, chat_id, telegram_msg_id, sender_id, sender_name, message_text,
		       is_transcribed, is_forward, forward_source_name, created_at, saved_at
		FROM raw_messages
		WHERE chat_id = $1 AND created_at >= $2 AND created_at < $3
		ORDER BY created_at ASC
	`
	
	rows, err := r.db.Query(query, chatID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()
	
	var messages []*database.RawMessage
	for rows.Next() {
		msg := &database.RawMessage{}
		if err := rows.Scan(
			&msg.ID,
			&msg.ChatID,
			&msg.TelegramMsgID,
			&msg.SenderID,
			&msg.SenderName,
			&msg.MessageText,
			&msg.IsTranscribed,
			&msg.IsForward,
			&msg.ForwardSourceName,
			&msg.CreatedAt,
			&msg.SavedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}
	
	return messages, nil
}

// GetLast24Hours retrieves messages from the last 24 hours for a chat
func (r *RawMessageRepository) GetLast24Hours(chatID int64) ([]*database.RawMessage, error) {
	now := time.Now()
	start := now.Add(-24 * time.Hour)
	return r.GetByChatIDAndTimeRange(chatID, start, now)
}

// DeleteOlderThan removes messages older than the specified time
func (r *RawMessageRepository) DeleteOlderThan(olderThan time.Time) (int64, error) {
	query := `DELETE FROM raw_messages WHERE created_at < $1`
	
	result, err := r.db.Exec(query, olderThan)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old messages: %w", err)
	}
	
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}
	
	return count, nil
}

// DeleteOlderThanDays removes messages older than specified days
func (r *RawMessageRepository) DeleteOlderThanDays(days int) (int64, error) {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	return r.DeleteOlderThan(cutoffTime)
}

// CountByChat returns the number of messages for a chat
func (r *RawMessageRepository) CountByChat(chatID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM raw_messages WHERE chat_id = $1`
	
	err := r.db.QueryRow(query, chatID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}
	
	return count, nil
}

// CountTotal returns the total number of messages
func (r *RawMessageRepository) CountTotal() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM raw_messages`
	
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count total messages: %w", err)
	}
	
	return count, nil
}
