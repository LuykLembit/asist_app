-- Migration: Create session_storage table
-- Purpose: Store userbot session encryption keys for persistence across restarts

CREATE TABLE IF NOT EXISTS session_storage (
    key TEXT PRIMARY KEY,
    value BYTEA NOT NULL
);

-- Index for faster lookups (though typically only one session)
CREATE INDEX IF NOT EXISTS idx_session_storage_key ON session_storage(key);
