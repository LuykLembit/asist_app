-- Migration: Create monitored_chats table
-- Purpose: Registry of Telegram chats being monitored

CREATE TABLE IF NOT EXISTS monitored_chats (
    chat_id BIGINT PRIMARY KEY,
    title TEXT,
    username TEXT,
    last_processed_msg_id INTEGER DEFAULT 0,
    last_pts INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    added_at TIMESTAMP DEFAULT NOW()
);

-- Index for filtering active chats
CREATE INDEX IF NOT EXISTS idx_monitored_chats_is_active ON monitored_chats(is_active);

-- Index for username lookups
CREATE INDEX IF NOT EXISTS idx_monitored_chats_username ON monitored_chats(username) WHERE username IS NOT NULL;
