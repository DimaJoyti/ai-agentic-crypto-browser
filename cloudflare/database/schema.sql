-- Cloudflare D1 Database Schema for AI Agentic Crypto Browser
-- SQLite-compatible schema migrated from PostgreSQL

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT DEFAULT 'user',
    created_at TEXT NOT NULL,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    is_active INTEGER DEFAULT 1,
    email_verified INTEGER DEFAULT 0
);

-- User sessions
CREATE TABLE IF NOT EXISTS user_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL,
    last_accessed TEXT,
    ip_address TEXT,
    user_agent TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- User wallets
CREATE TABLE IF NOT EXISTS user_wallets (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    user_id TEXT NOT NULL,
    address TEXT NOT NULL,
    signature TEXT,
    chain TEXT DEFAULT 'ethereum',
    connected_at TEXT NOT NULL,
    is_primary INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, address, chain)
);

-- Conversations
CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    is_archived INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Messages
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    created_at TEXT NOT NULL,
    tokens_used INTEGER DEFAULT 0,
    model_used TEXT,
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);

-- Trading orders
CREATE TABLE IF NOT EXISTS trading_orders (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    symbol TEXT NOT NULL,
    side TEXT NOT NULL CHECK (side IN ('buy', 'sell')),
    type TEXT NOT NULL CHECK (type IN ('market', 'limit', 'stop', 'stop_limit')),
    quantity REAL NOT NULL,
    price REAL,
    stop_price REAL,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'filled', 'cancelled', 'rejected')),
    filled_quantity REAL DEFAULT 0,
    average_price REAL,
    created_at TEXT NOT NULL,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    exchange TEXT DEFAULT 'binance',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Transactions
CREATE TABLE IF NOT EXISTS transactions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    tx_hash TEXT,
    from_address TEXT,
    to_address TEXT NOT NULL,
    amount REAL NOT NULL,
    token TEXT DEFAULT 'ETH',
    chain TEXT DEFAULT 'ethereum',
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'failed')),
    gas_used REAL,
    gas_price REAL,
    block_number INTEGER,
    created_at TEXT NOT NULL,
    confirmed_at TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Portfolio holdings
CREATE TABLE IF NOT EXISTS portfolio_holdings (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    user_id TEXT NOT NULL,
    symbol TEXT NOT NULL,
    quantity REAL NOT NULL,
    average_cost REAL NOT NULL,
    current_price REAL,
    last_updated TEXT DEFAULT CURRENT_TIMESTAMP,
    exchange TEXT,
    wallet_address TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, symbol, exchange, wallet_address)
);

-- AI analysis results
CREATE TABLE IF NOT EXISTS ai_analysis (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    user_id TEXT,
    symbol TEXT NOT NULL,
    analysis_type TEXT NOT NULL,
    timeframe TEXT DEFAULT '1h',
    result TEXT NOT NULL, -- JSON string
    confidence REAL,
    created_at TEXT NOT NULL,
    expires_at TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Price predictions
CREATE TABLE IF NOT EXISTS price_predictions (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    symbol TEXT NOT NULL,
    predicted_price REAL NOT NULL,
    current_price REAL NOT NULL,
    timeframe TEXT NOT NULL,
    confidence REAL NOT NULL,
    model_used TEXT,
    factors TEXT, -- JSON string
    created_at TEXT NOT NULL,
    target_time TEXT NOT NULL,
    actual_price REAL,
    accuracy REAL
);

-- Risk assessments
CREATE TABLE IF NOT EXISTS risk_assessments (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    user_id TEXT NOT NULL,
    risk_score REAL NOT NULL,
    risk_level TEXT NOT NULL,
    factors TEXT, -- JSON string
    recommendations TEXT, -- JSON string
    portfolio_value REAL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Market data cache
CREATE TABLE IF NOT EXISTS market_data (
    symbol TEXT PRIMARY KEY,
    price REAL NOT NULL,
    volume_24h REAL,
    change_24h REAL,
    market_cap REAL,
    last_updated TEXT NOT NULL
);

-- User preferences
CREATE TABLE IF NOT EXISTS user_preferences (
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
CREATE TABLE IF NOT EXISTS user_api_keys (
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

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token_hash ON user_sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_user_wallets_user_id ON user_wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_trading_orders_user_id ON trading_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_trading_orders_symbol ON trading_orders(symbol);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_tx_hash ON transactions(tx_hash);
CREATE INDEX IF NOT EXISTS idx_portfolio_holdings_user_id ON portfolio_holdings(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_analysis_symbol ON ai_analysis(symbol);
CREATE INDEX IF NOT EXISTS idx_ai_analysis_created_at ON ai_analysis(created_at);
CREATE INDEX IF NOT EXISTS idx_price_predictions_symbol ON price_predictions(symbol);
CREATE INDEX IF NOT EXISTS idx_risk_assessments_user_id ON risk_assessments(user_id);
CREATE INDEX IF NOT EXISTS idx_market_data_last_updated ON market_data(last_updated);
CREATE INDEX IF NOT EXISTS idx_user_api_keys_user_id ON user_api_keys(user_id);
