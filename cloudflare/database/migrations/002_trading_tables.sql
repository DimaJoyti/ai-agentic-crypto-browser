-- Migration 002: Trading and transaction tables

-- Trading orders
CREATE TABLE trading_orders (
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

-- Blockchain transactions
CREATE TABLE transactions (
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
CREATE TABLE portfolio_holdings (
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

-- Create indexes
CREATE INDEX idx_trading_orders_user_id ON trading_orders(user_id);
CREATE INDEX idx_trading_orders_symbol ON trading_orders(symbol);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_tx_hash ON transactions(tx_hash);
CREATE INDEX idx_portfolio_holdings_user_id ON portfolio_holdings(user_id);
