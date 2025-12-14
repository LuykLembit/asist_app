-- Migration: Create daily_reports table
-- Purpose: Store AI-generated intelligence reports with indefinite retention

CREATE TABLE IF NOT EXISTS daily_reports (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES monitored_chats(chat_id) ON DELETE CASCADE,
    report_date DATE NOT NULL,
    summary TEXT,
    full_json JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(chat_id, report_date)
);

-- Index for date-based queries
CREATE INDEX IF NOT EXISTS idx_daily_reports_report_date ON daily_reports(report_date DESC);

-- Index for chat-based queries
CREATE INDEX IF NOT EXISTS idx_daily_reports_chat_id ON daily_reports(chat_id);

-- GIN index for JSONB search capabilities
CREATE INDEX IF NOT EXISTS idx_daily_reports_full_json ON daily_reports USING GIN(full_json);

-- Index for full-text search on summary
CREATE INDEX IF NOT EXISTS idx_daily_reports_summary_fts ON daily_reports USING GIN(to_tsvector('english', summary));
