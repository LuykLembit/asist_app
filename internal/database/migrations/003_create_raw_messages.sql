-- Migration: Create raw_messages table
-- Purpose: Archive of collected messages with 7-day retention

CREATE TABLE IF NOT EXISTS raw_messages (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES monitored_chats(chat_id) ON DELETE CASCADE,
    telegram_msg_id INTEGER NOT NULL,
    sender_id BIGINT,
    sender_name TEXT,
    message_text TEXT,
    is_transcribed BOOLEAN DEFAULT FALSE,
    is_forward BOOLEAN DEFAULT FALSE,
    forward_source_name TEXT,
    created_at TIMESTAMP NOT NULL,
    saved_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(chat_id, telegram_msg_id)
);

-- Index for time-based queries (24h reports, retention cleanup)
CREATE INDEX IF NOT EXISTS idx_raw_messages_created_at ON raw_messages(created_at DESC);

-- Index for chat-based queries
CREATE INDEX IF NOT EXISTS idx_raw_messages_chat_id ON raw_messages(chat_id);

-- Index for sender lookups
CREATE INDEX IF NOT EXISTS idx_raw_messages_sender_id ON raw_messages(sender_id) WHERE sender_id IS NOT NULL;

-- Composite index for chat + time range queries (most common pattern)
CREATE INDEX IF NOT EXISTS idx_raw_messages_chat_created ON raw_messages(chat_id, created_at DESC);
