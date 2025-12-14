-- Migration: Create triggers table
-- Purpose: Keyword definitions for real-time alerts

CREATE TABLE IF NOT EXISTS triggers (
    id SERIAL PRIMARY KEY,
    phrase TEXT NOT NULL,
    is_regex BOOLEAN DEFAULT FALSE,
    alert_level VARCHAR(20) DEFAULT 'info',
    CHECK (alert_level IN ('info', 'warning', 'critical'))
);

-- Index for phrase lookups
CREATE INDEX IF NOT EXISTS idx_triggers_phrase ON triggers(phrase);

-- Index for filtering by alert level
CREATE INDEX IF NOT EXISTS idx_triggers_alert_level ON triggers(alert_level);
