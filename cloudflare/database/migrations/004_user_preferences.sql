-- Migration 004: User preferences and API keys

-- User preferences
CREATE TABLE user_preferences (
    user_id TEXT PRIMARY KEY,
    theme TEXT DEFAULT 'dark',
    language TEXT DEFAULT 'en',
    timezone TEXT DEFAULT 'UTC',
    notifications TEXT, -- JSON string
    trading_preferences TEXT, -- JSON string
    created_at TEXT NOT NULL,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- API keys (encrypted)
CREATE TABLE user_api_keys (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    user_id TEXT NOT NULL,
    exchange TEXT NOT NULL,
    api_key_encrypted TEXT NOT NULL,
    api_secret_encrypted TEXT NOT NULL,
    passphrase_encrypted TEXT,
    is_active INTEGER DEFAULT 1,
    created_at TEXT NOT NULL,
    last_used TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_user_api_keys_user_id ON user_api_keys(user_id);
