package database

import (
	"database/sql"
	"time"
)

// SessionStorage represents a userbot session entry
type SessionStorage struct {
	Key   string
	Value []byte
}

// MonitoredChat represents a Telegram chat being monitored
type MonitoredChat struct {
	ChatID              int64
	Title               sql.NullString
	Username            sql.NullString
	LastProcessedMsgID  int
	LastPts             int
	IsActive            bool
	AddedAt             time.Time
}

// RawMessage represents a collected Telegram message
type RawMessage struct {
	ID                int
	ChatID            int64
	TelegramMsgID     int
	SenderID          sql.NullInt64
	SenderName        sql.NullString
	MessageText       sql.NullString
	IsTranscribed     bool
	IsForward         bool
	ForwardSourceName sql.NullString
	CreatedAt         time.Time
	SavedAt           time.Time
}

// Trigger represents a keyword alert trigger
type Trigger struct {
	ID         int
	Phrase     string
	IsRegex    bool
	AlertLevel string
}

// DailyReport represents an AI-generated intelligence report
type DailyReport struct {
	ID         int
	ChatID     int64
	ReportDate time.Time
	Summary    sql.NullString
	FullJSON   []byte // JSONB stored as bytes
	CreatedAt  time.Time
}
