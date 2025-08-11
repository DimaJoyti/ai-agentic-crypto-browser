-- Trading Bots Database Schema
-- Migration 006: Create tables for 7 trading bots system

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Trading Bot Configurations Table
CREATE TABLE IF NOT EXISTS trading_bot_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    strategy VARCHAR(50) NOT NULL CHECK (strategy IN ('dca', 'grid', 'momentum', 'mean_reversion', 'arbitrage', 'scalping', 'swing')),
    trading_pairs TEXT[] NOT NULL,
    exchange VARCHAR(50) NOT NULL,
    base_currency VARCHAR(10) NOT NULL DEFAULT 'USDT',
    strategy_params JSONB NOT NULL DEFAULT '{}',
    risk_profile JSONB NOT NULL DEFAULT '{}',
    capital_config JSONB NOT NULL DEFAULT '{}',
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    
    CONSTRAINT unique_bot_name_per_user UNIQUE (name, created_by)
);

-- Trading Bot States Table
CREATE TABLE IF NOT EXISTS trading_bot_states (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bot_id UUID NOT NULL REFERENCES trading_bot_configs(id) ON DELETE CASCADE,
    state VARCHAR(20) NOT NULL CHECK (state IN ('idle', 'running', 'paused', 'stopped', 'error')),
    last_execution TIMESTAMP WITH TIME ZONE,
    error_count INTEGER DEFAULT 0,
    error_message TEXT,
    runtime_data JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_bot_state UNIQUE (bot_id)
);

-- Trading Bot Performance Metrics Table
CREATE TABLE IF NOT EXISTS trading_bot_performance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bot_id UUID NOT NULL REFERENCES trading_bot_configs(id) ON DELETE CASCADE,
    total_trades INTEGER DEFAULT 0,
    winning_trades INTEGER DEFAULT 0,
    losing_trades INTEGER DEFAULT 0,
    win_rate DECIMAL(5,4) DEFAULT 0.0000,
    total_profit DECIMAL(20,8) DEFAULT 0.00000000,
    total_loss DECIMAL(20,8) DEFAULT 0.00000000,
    net_profit DECIMAL(20,8) DEFAULT 0.00000000,
    max_drawdown DECIMAL(5,4) DEFAULT 0.0000,
    sharpe_ratio DECIMAL(8,4) DEFAULT 0.0000,
    roi DECIMAL(8,4) DEFAULT 0.0000,
    volatility DECIMAL(8,4) DEFAULT 0.0000,
    avg_trade_duration INTERVAL,
    best_trade DECIMAL(20,8) DEFAULT 0.00000000,
    worst_trade DECIMAL(20,8) DEFAULT 0.00000000,
    consecutive_wins INTEGER DEFAULT 0,
    consecutive_losses INTEGER DEFAULT 0,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_bot_performance UNIQUE (bot_id)
);

-- Trading Bot Trade History Table
CREATE TABLE IF NOT EXISTS trading_bot_trades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bot_id UUID NOT NULL REFERENCES trading_bot_configs(id) ON DELETE CASCADE,
    trade_id VARCHAR(255) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) NOT NULL CHECK (side IN ('buy', 'sell')),
    order_type VARCHAR(20) NOT NULL CHECK (order_type IN ('market', 'limit', 'stop', 'stop_limit')),
    amount DECIMAL(20,8) NOT NULL,
    price DECIMAL(20,8) NOT NULL,
    executed_amount DECIMAL(20,8) DEFAULT 0.00000000,
    executed_price DECIMAL(20,8) DEFAULT 0.00000000,
    fee DECIMAL(20,8) DEFAULT 0.00000000,
    fee_currency VARCHAR(10),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'executed', 'cancelled', 'failed', 'partially_filled')),
    strategy_metadata JSONB DEFAULT '{}',
    exchange_order_id VARCHAR(255),
    executed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_trade_id_per_bot UNIQUE (bot_id, trade_id)
);

-- Trading Bot Positions Table
CREATE TABLE IF NOT EXISTS trading_bot_positions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bot_id UUID NOT NULL REFERENCES trading_bot_configs(id) ON DELETE CASCADE,
    position_id VARCHAR(255) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) NOT NULL CHECK (side IN ('long', 'short')),
    size DECIMAL(20,8) NOT NULL,
    entry_price DECIMAL(20,8) NOT NULL,
    current_price DECIMAL(20,8),
    unrealized_pnl DECIMAL(20,8) DEFAULT 0.00000000,
    realized_pnl DECIMAL(20,8) DEFAULT 0.00000000,
    stop_loss DECIMAL(20,8),
    take_profit DECIMAL(20,8),
    margin_used DECIMAL(20,8) DEFAULT 0.00000000,
    leverage DECIMAL(8,2) DEFAULT 1.00,
    status VARCHAR(20) NOT NULL CHECK (status IN ('open', 'closed', 'liquidated')),
    opened_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    closed_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT unique_position_id_per_bot UNIQUE (bot_id, position_id)
);

-- Trading Bot Alerts Table
CREATE TABLE IF NOT EXISTS trading_bot_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bot_id UUID NOT NULL REFERENCES trading_bot_configs(id) ON DELETE CASCADE,
    alert_type VARCHAR(50) NOT NULL CHECK (alert_type IN ('error', 'warning', 'info', 'performance', 'risk')),
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    acknowledged BOOLEAN DEFAULT false,
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    acknowledged_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Trading Bot Backtests Table
CREATE TABLE IF NOT EXISTS trading_bot_backtests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bot_id UUID NOT NULL REFERENCES trading_bot_configs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    initial_balance DECIMAL(20,8) NOT NULL,
    final_balance DECIMAL(20,8),
    total_return DECIMAL(8,4),
    max_drawdown DECIMAL(5,4),
    sharpe_ratio DECIMAL(8,4),
    total_trades INTEGER DEFAULT 0,
    win_rate DECIMAL(5,4),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'running', 'completed', 'failed')),
    results JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Exchange API Keys Table (Encrypted)
CREATE TABLE IF NOT EXISTS exchange_api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exchange VARCHAR(50) NOT NULL,
    api_key_encrypted TEXT NOT NULL,
    api_secret_encrypted TEXT NOT NULL,
    passphrase_encrypted TEXT,
    sandbox_mode BOOLEAN DEFAULT true,
    permissions TEXT[] DEFAULT ARRAY['read'],
    is_active BOOLEAN DEFAULT true,
    last_used TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_exchange_per_user UNIQUE (user_id, exchange)
);

-- Market Data Cache Table
CREATE TABLE IF NOT EXISTS market_data_cache (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol VARCHAR(20) NOT NULL,
    exchange VARCHAR(50) NOT NULL,
    price DECIMAL(20,8) NOT NULL,
    volume DECIMAL(20,8) NOT NULL,
    high_24h DECIMAL(20,8),
    low_24h DECIMAL(20,8),
    change_24h DECIMAL(8,4),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_symbol_exchange_timestamp UNIQUE (symbol, exchange, timestamp)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_trading_bot_configs_strategy ON trading_bot_configs(strategy);
CREATE INDEX IF NOT EXISTS idx_trading_bot_configs_enabled ON trading_bot_configs(enabled);
CREATE INDEX IF NOT EXISTS idx_trading_bot_configs_created_by ON trading_bot_configs(created_by);

CREATE INDEX IF NOT EXISTS idx_trading_bot_states_bot_id ON trading_bot_states(bot_id);
CREATE INDEX IF NOT EXISTS idx_trading_bot_states_state ON trading_bot_states(state);
CREATE INDEX IF NOT EXISTS idx_trading_bot_states_last_execution ON trading_bot_states(last_execution);

CREATE INDEX IF NOT EXISTS idx_trading_bot_performance_bot_id ON trading_bot_performance(bot_id);
CREATE INDEX IF NOT EXISTS idx_trading_bot_performance_net_profit ON trading_bot_performance(net_profit);
CREATE INDEX IF NOT EXISTS idx_trading_bot_performance_win_rate ON trading_bot_performance(win_rate);

CREATE INDEX IF NOT EXISTS idx_trading_bot_trades_bot_id ON trading_bot_trades(bot_id);
CREATE INDEX IF NOT EXISTS idx_trading_bot_trades_symbol ON trading_bot_trades(symbol);
CREATE INDEX IF NOT EXISTS idx_trading_bot_trades_status ON trading_bot_trades(status);
CREATE INDEX IF NOT EXISTS idx_trading_bot_trades_created_at ON trading_bot_trades(created_at);
CREATE INDEX IF NOT EXISTS idx_trading_bot_trades_executed_at ON trading_bot_trades(executed_at);

CREATE INDEX IF NOT EXISTS idx_trading_bot_positions_bot_id ON trading_bot_positions(bot_id);
CREATE INDEX IF NOT EXISTS idx_trading_bot_positions_symbol ON trading_bot_positions(symbol);
CREATE INDEX IF NOT EXISTS idx_trading_bot_positions_status ON trading_bot_positions(status);

CREATE INDEX IF NOT EXISTS idx_trading_bot_alerts_bot_id ON trading_bot_alerts(bot_id);
CREATE INDEX IF NOT EXISTS idx_trading_bot_alerts_severity ON trading_bot_alerts(severity);
CREATE INDEX IF NOT EXISTS idx_trading_bot_alerts_acknowledged ON trading_bot_alerts(acknowledged);
CREATE INDEX IF NOT EXISTS idx_trading_bot_alerts_created_at ON trading_bot_alerts(created_at);

CREATE INDEX IF NOT EXISTS idx_trading_bot_backtests_bot_id ON trading_bot_backtests(bot_id);
CREATE INDEX IF NOT EXISTS idx_trading_bot_backtests_status ON trading_bot_backtests(status);

CREATE INDEX IF NOT EXISTS idx_exchange_api_keys_user_id ON exchange_api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_exchange_api_keys_exchange ON exchange_api_keys(exchange);
CREATE INDEX IF NOT EXISTS idx_exchange_api_keys_is_active ON exchange_api_keys(is_active);

CREATE INDEX IF NOT EXISTS idx_market_data_cache_symbol_exchange ON market_data_cache(symbol, exchange);
CREATE INDEX IF NOT EXISTS idx_market_data_cache_timestamp ON market_data_cache(timestamp);

-- Create triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_trading_bot_configs_updated_at 
    BEFORE UPDATE ON trading_bot_configs 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_trading_bot_states_updated_at 
    BEFORE UPDATE ON trading_bot_states 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_exchange_api_keys_updated_at 
    BEFORE UPDATE ON exchange_api_keys 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create views for easier querying
CREATE OR REPLACE VIEW trading_bots_overview AS
SELECT 
    c.id,
    c.name,
    c.strategy,
    c.trading_pairs,
    c.exchange,
    c.enabled,
    s.state,
    s.last_execution,
    s.error_count,
    p.total_trades,
    p.win_rate,
    p.net_profit,
    p.max_drawdown,
    c.created_at
FROM trading_bot_configs c
LEFT JOIN trading_bot_states s ON c.id = s.bot_id
LEFT JOIN trading_bot_performance p ON c.id = p.bot_id;

CREATE OR REPLACE VIEW active_trading_bots AS
SELECT * FROM trading_bots_overview 
WHERE enabled = true AND state IN ('running', 'paused');

CREATE OR REPLACE VIEW bot_performance_summary AS
SELECT 
    c.strategy,
    COUNT(*) as bot_count,
    AVG(p.win_rate) as avg_win_rate,
    SUM(p.net_profit) as total_profit,
    AVG(p.max_drawdown) as avg_max_drawdown,
    SUM(p.total_trades) as total_trades
FROM trading_bot_configs c
JOIN trading_bot_performance p ON c.id = p.bot_id
WHERE c.enabled = true
GROUP BY c.strategy;

-- Insert default data for testing
INSERT INTO trading_bot_configs (name, strategy, trading_pairs, exchange, strategy_params, risk_profile, capital_config) VALUES
('DCA Bitcoin Bot', 'dca', ARRAY['BTC/USDT'], 'binance', 
 '{"investment_amount": 100, "interval": "1h", "max_deviation": 0.05}',
 '{"max_position_size": 0.20, "stop_loss": 0.15, "take_profit": 0.30, "max_drawdown": 0.10}',
 '{"initial_balance": 10000, "allocation_percentage": 0.20}'),
 
('Grid BNB Bot', 'grid', ARRAY['BNB/USDT'], 'binance',
 '{"grid_levels": 20, "grid_spacing": 0.02, "upper_bound": 1.20, "lower_bound": 0.80, "order_amount": 50}',
 '{"max_position_size": 0.15, "stop_loss": 0.20, "take_profit": 0.25, "max_drawdown": 0.12}',
 '{"initial_balance": 7500, "allocation_percentage": 0.15}'),
 
('Momentum SOL Bot', 'momentum', ARRAY['SOL/USDT'], 'coinbase',
 '{"momentum_period": 14, "rsi_threshold_buy": 30, "rsi_threshold_sell": 70, "volume_threshold": 1.5}',
 '{"max_position_size": 0.10, "stop_loss": 0.08, "take_profit": 0.15, "max_drawdown": 0.15}',
 '{"initial_balance": 5000, "allocation_percentage": 0.10}');

COMMENT ON TABLE trading_bot_configs IS 'Configuration settings for trading bots';
COMMENT ON TABLE trading_bot_states IS 'Current state and runtime information for trading bots';
COMMENT ON TABLE trading_bot_performance IS 'Performance metrics and statistics for trading bots';
COMMENT ON TABLE trading_bot_trades IS 'Trade execution history for trading bots';
COMMENT ON TABLE trading_bot_positions IS 'Current and historical positions for trading bots';
COMMENT ON TABLE trading_bot_alerts IS 'Alerts and notifications generated by trading bots';
COMMENT ON TABLE trading_bot_backtests IS 'Backtest results and configurations for trading bots';
COMMENT ON TABLE exchange_api_keys IS 'Encrypted API keys for cryptocurrency exchanges';
COMMENT ON TABLE market_data_cache IS 'Cached market data for trading pairs';
