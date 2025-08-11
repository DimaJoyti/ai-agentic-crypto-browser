-- Migration 003: AI and analytics tables

-- AI analysis results
CREATE TABLE ai_analysis (
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
CREATE TABLE price_predictions (
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
CREATE TABLE risk_assessments (
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
CREATE TABLE market_data (
    symbol TEXT PRIMARY KEY,
    price REAL NOT NULL,
    volume_24h REAL,
    change_24h REAL,
    market_cap REAL,
    last_updated TEXT NOT NULL
);

-- Create indexes
CREATE INDEX idx_ai_analysis_symbol ON ai_analysis(symbol);
CREATE INDEX idx_ai_analysis_created_at ON ai_analysis(created_at);
CREATE INDEX idx_price_predictions_symbol ON price_predictions(symbol);
CREATE INDEX idx_risk_assessments_user_id ON risk_assessments(user_id);
CREATE INDEX idx_market_data_last_updated ON market_data(last_updated);
