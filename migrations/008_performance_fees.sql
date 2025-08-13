-- Performance-Based Fee System Schema
-- Migration 008: Performance fee tracking and high-water mark system

-- Performance Fee Configuration table
CREATE TABLE IF NOT EXISTS performance_fee_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL UNIQUE,
    fee_percentage DECIMAL(5,2) NOT NULL DEFAULT 20.00, -- 20% default
    high_water_mark DECIMAL(15,2) NOT NULL DEFAULT 10000.00, -- Starting portfolio value
    minimum_fee DECIMAL(10,2) NOT NULL DEFAULT 1.00, -- Minimum fee per trade
    maximum_fee DECIMAL(10,2) NOT NULL DEFAULT 1000.00, -- Maximum fee per trade
    fee_frequency VARCHAR(20) NOT NULL DEFAULT 'per_trade', -- per_trade, monthly, quarterly
    only_profitable_trades BOOLEAN NOT NULL DEFAULT true,
    hurdle_rate DECIMAL(5,2) DEFAULT 0.00, -- Minimum return before fees apply
    fee_tier VARCHAR(20) DEFAULT 'standard', -- standard, premium, enterprise
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Trade Records table for performance tracking
CREATE TABLE IF NOT EXISTS trade_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) NOT NULL, -- buy, sell, long, short
    quantity DECIMAL(20,8) NOT NULL,
    entry_price DECIMAL(20,8) NOT NULL,
    exit_price DECIMAL(20,8),
    entry_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    exit_timestamp TIMESTAMP WITH TIME ZONE,
    pnl DECIMAL(15,2) DEFAULT 0.00, -- Profit and Loss
    pnl_percentage DECIMAL(8,4) DEFAULT 0.00, -- PnL as percentage
    fees DECIMAL(10,2) DEFAULT 0.00, -- Trading fees paid to exchange
    performance_fee DECIMAL(10,2) DEFAULT 0.00, -- Our performance fee
    status VARCHAR(20) DEFAULT 'open', -- open, closed, cancelled
    strategy_id VARCHAR(255),
    exchange VARCHAR(50),
    order_type VARCHAR(20), -- market, limit, stop
    leverage DECIMAL(5,2) DEFAULT 1.00,
    risk_score DECIMAL(3,2), -- 0.00 to 1.00
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Performance Summaries table for analytics
CREATE TABLE IF NOT EXISTS performance_summaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    period VARCHAR(20) NOT NULL, -- daily, weekly, monthly, quarterly, yearly, all_time
    period_start DATE,
    period_end DATE,
    total_trades BIGINT DEFAULT 0,
    profitable_trades BIGINT DEFAULT 0,
    losing_trades BIGINT DEFAULT 0,
    total_pnl DECIMAL(15,2) DEFAULT 0.00,
    total_performance_fees DECIMAL(15,2) DEFAULT 0.00,
    win_rate DECIMAL(5,4) DEFAULT 0.0000, -- 0.0000 to 1.0000
    average_return DECIMAL(8,4) DEFAULT 0.0000,
    max_drawdown DECIMAL(8,4) DEFAULT 0.0000,
    sharpe_ratio DECIMAL(6,4) DEFAULT 0.0000,
    sortino_ratio DECIMAL(6,4) DEFAULT 0.0000,
    calmar_ratio DECIMAL(6,4) DEFAULT 0.0000,
    current_high_water_mark DECIMAL(15,2) DEFAULT 0.00,
    max_consecutive_wins INTEGER DEFAULT 0,
    max_consecutive_losses INTEGER DEFAULT 0,
    largest_win DECIMAL(15,2) DEFAULT 0.00,
    largest_loss DECIMAL(15,2) DEFAULT 0.00,
    average_win DECIMAL(15,2) DEFAULT 0.00,
    average_loss DECIMAL(15,2) DEFAULT 0.00,
    profit_factor DECIMAL(6,4) DEFAULT 0.0000, -- Gross profit / Gross loss
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, period, period_start)
);

-- Performance Fee Billing table for invoicing
CREATE TABLE IF NOT EXISTS performance_fee_billing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    billing_period DATE NOT NULL,
    total_trades INTEGER DEFAULT 0,
    profitable_trades INTEGER DEFAULT 0,
    total_pnl DECIMAL(15,2) DEFAULT 0.00,
    total_performance_fees DECIMAL(15,2) DEFAULT 0.00,
    fee_rate DECIMAL(5,2) NOT NULL,
    high_water_mark_start DECIMAL(15,2) NOT NULL,
    high_water_mark_end DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, invoiced, paid, disputed
    stripe_invoice_id VARCHAR(255),
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    paid_at TIMESTAMP WITH TIME ZONE,
    due_date TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}'
);

-- Portfolio Snapshots for high-water mark tracking
CREATE TABLE IF NOT EXISTS portfolio_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    snapshot_date DATE NOT NULL,
    portfolio_value DECIMAL(15,2) NOT NULL,
    cash_balance DECIMAL(15,2) DEFAULT 0.00,
    positions_value DECIMAL(15,2) DEFAULT 0.00,
    unrealized_pnl DECIMAL(15,2) DEFAULT 0.00,
    realized_pnl DECIMAL(15,2) DEFAULT 0.00,
    total_deposits DECIMAL(15,2) DEFAULT 0.00,
    total_withdrawals DECIMAL(15,2) DEFAULT 0.00,
    is_high_water_mark BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, snapshot_date)
);

-- Fee Tier Definitions table
CREATE TABLE IF NOT EXISTS performance_fee_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tier_name VARCHAR(20) NOT NULL UNIQUE,
    fee_percentage DECIMAL(5,2) NOT NULL,
    minimum_portfolio_value DECIMAL(15,2) DEFAULT 0.00,
    maximum_portfolio_value DECIMAL(15,2),
    hurdle_rate DECIMAL(5,2) DEFAULT 0.00,
    high_water_mark_required BOOLEAN DEFAULT true,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_performance_fee_configs_user_id ON performance_fee_configs(user_id);

CREATE INDEX IF NOT EXISTS idx_trade_records_user_id ON trade_records(user_id);
CREATE INDEX IF NOT EXISTS idx_trade_records_status ON trade_records(status);
CREATE INDEX IF NOT EXISTS idx_trade_records_timestamp ON trade_records(entry_timestamp);
CREATE INDEX IF NOT EXISTS idx_trade_records_user_status ON trade_records(user_id, status);
CREATE INDEX IF NOT EXISTS idx_trade_records_user_date ON trade_records(user_id, DATE(entry_timestamp));
CREATE INDEX IF NOT EXISTS idx_trade_records_symbol ON trade_records(symbol);
CREATE INDEX IF NOT EXISTS idx_trade_records_strategy ON trade_records(strategy_id);

CREATE INDEX IF NOT EXISTS idx_performance_summaries_user_period ON performance_summaries(user_id, period);
CREATE INDEX IF NOT EXISTS idx_performance_summaries_period_dates ON performance_summaries(period, period_start, period_end);

CREATE INDEX IF NOT EXISTS idx_performance_fee_billing_user_period ON performance_fee_billing(user_id, billing_period);
CREATE INDEX IF NOT EXISTS idx_performance_fee_billing_status ON performance_fee_billing(status);

CREATE INDEX IF NOT EXISTS idx_portfolio_snapshots_user_date ON portfolio_snapshots(user_id, snapshot_date);
CREATE INDEX IF NOT EXISTS idx_portfolio_snapshots_high_water_mark ON portfolio_snapshots(user_id, is_high_water_mark) WHERE is_high_water_mark = true;

-- Functions for automatic calculations
CREATE OR REPLACE FUNCTION calculate_trade_pnl()
RETURNS TRIGGER AS $$
BEGIN
    -- Calculate PnL when trade is closed
    IF NEW.status = 'closed' AND NEW.exit_price IS NOT NULL THEN
        IF NEW.side = 'buy' OR NEW.side = 'long' THEN
            NEW.pnl = (NEW.exit_price - NEW.entry_price) * NEW.quantity;
        ELSE -- sell or short
            NEW.pnl = (NEW.entry_price - NEW.exit_price) * NEW.quantity;
        END IF;
        
        -- Calculate PnL percentage
        IF NEW.entry_price > 0 THEN
            NEW.pnl_percentage = (NEW.pnl / (NEW.entry_price * NEW.quantity)) * 100;
        END IF;
        
        NEW.updated_at = NOW();
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically calculate PnL
CREATE TRIGGER trigger_calculate_trade_pnl
    BEFORE UPDATE ON trade_records
    FOR EACH ROW
    EXECUTE FUNCTION calculate_trade_pnl();

-- Function to update performance summary
CREATE OR REPLACE FUNCTION update_performance_summary_on_trade()
RETURNS TRIGGER AS $$
BEGIN
    -- Update all-time summary when trade is completed
    IF NEW.status = 'closed' AND (OLD.status IS NULL OR OLD.status != 'closed') THEN
        INSERT INTO performance_summaries (
            user_id, period, total_trades, profitable_trades, losing_trades,
            total_pnl, win_rate, largest_win, largest_loss
        )
        SELECT 
            NEW.user_id,
            'all_time',
            COUNT(*),
            COUNT(CASE WHEN pnl > 0 THEN 1 END),
            COUNT(CASE WHEN pnl < 0 THEN 1 END),
            COALESCE(SUM(pnl), 0),
            CASE WHEN COUNT(*) > 0 THEN 
                COUNT(CASE WHEN pnl > 0 THEN 1 END)::decimal / COUNT(*)::decimal 
            ELSE 0 END,
            COALESCE(MAX(pnl), 0),
            COALESCE(MIN(pnl), 0)
        FROM trade_records 
        WHERE user_id = NEW.user_id AND status = 'closed'
        ON CONFLICT (user_id, period, period_start)
        DO UPDATE SET
            total_trades = EXCLUDED.total_trades,
            profitable_trades = EXCLUDED.profitable_trades,
            losing_trades = EXCLUDED.losing_trades,
            total_pnl = EXCLUDED.total_pnl,
            win_rate = EXCLUDED.win_rate,
            largest_win = EXCLUDED.largest_win,
            largest_loss = EXCLUDED.largest_loss,
            last_updated = NOW();
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update performance summary
CREATE TRIGGER trigger_update_performance_summary
    AFTER UPDATE ON trade_records
    FOR EACH ROW
    EXECUTE FUNCTION update_performance_summary_on_trade();

-- Function to update high water mark
CREATE OR REPLACE FUNCTION update_high_water_mark()
RETURNS TRIGGER AS $$
DECLARE
    current_portfolio_value DECIMAL(15,2);
    current_hwm DECIMAL(15,2);
BEGIN
    -- Calculate current portfolio value
    SELECT COALESCE(SUM(pnl), 0) + 10000 INTO current_portfolio_value
    FROM trade_records 
    WHERE user_id = NEW.user_id AND status = 'closed';
    
    -- Get current high water mark
    SELECT high_water_mark INTO current_hwm
    FROM performance_fee_configs
    WHERE user_id = NEW.user_id;
    
    -- Update high water mark if current value is higher
    IF current_portfolio_value > current_hwm THEN
        UPDATE performance_fee_configs 
        SET high_water_mark = current_portfolio_value, updated_at = NOW()
        WHERE user_id = NEW.user_id;
        
        -- Create portfolio snapshot
        INSERT INTO portfolio_snapshots (
            user_id, snapshot_date, portfolio_value, is_high_water_mark
        ) VALUES (
            NEW.user_id, CURRENT_DATE, current_portfolio_value, true
        ) ON CONFLICT (user_id, snapshot_date) 
        DO UPDATE SET 
            portfolio_value = EXCLUDED.portfolio_value,
            is_high_water_mark = EXCLUDED.is_high_water_mark;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update high water mark
CREATE TRIGGER trigger_update_high_water_mark
    AFTER UPDATE ON trade_records
    FOR EACH ROW
    WHEN (NEW.status = 'closed' AND (OLD.status IS NULL OR OLD.status != 'closed'))
    EXECUTE FUNCTION update_high_water_mark();

-- Insert default fee tiers
INSERT INTO performance_fee_tiers (tier_name, fee_percentage, minimum_portfolio_value, description) VALUES
('starter', 15.00, 0, 'Starter tier - 15% performance fee'),
('standard', 20.00, 10000, 'Standard tier - 20% performance fee'),
('premium', 25.00, 100000, 'Premium tier - 25% performance fee'),
('enterprise', 30.00, 1000000, 'Enterprise tier - 30% performance fee')
ON CONFLICT (tier_name) DO NOTHING;

-- Create views for analytics
CREATE OR REPLACE VIEW performance_analytics AS
SELECT 
    ps.user_id,
    ps.period,
    ps.total_trades,
    ps.profitable_trades,
    ps.win_rate,
    ps.total_pnl,
    ps.total_performance_fees,
    ps.sharpe_ratio,
    ps.max_drawdown,
    pfc.fee_percentage,
    pfc.high_water_mark,
    ps.last_updated
FROM performance_summaries ps
JOIN performance_fee_configs pfc ON ps.user_id = pfc.user_id;

CREATE OR REPLACE VIEW top_performers AS
SELECT 
    user_id,
    total_pnl,
    total_performance_fees,
    win_rate,
    sharpe_ratio,
    total_trades,
    ROW_NUMBER() OVER (ORDER BY total_pnl DESC) as rank
FROM performance_summaries 
WHERE period = 'all_time' AND total_trades >= 10
ORDER BY total_pnl DESC;

-- Comments for documentation
COMMENT ON TABLE performance_fee_configs IS 'Configuration for performance-based fees per user';
COMMENT ON TABLE trade_records IS 'Individual trade records for performance tracking';
COMMENT ON TABLE performance_summaries IS 'Aggregated performance metrics by period';
COMMENT ON TABLE performance_fee_billing IS 'Monthly billing records for performance fees';
COMMENT ON TABLE portfolio_snapshots IS 'Daily portfolio value snapshots for high-water mark tracking';

COMMENT ON COLUMN performance_fee_configs.high_water_mark IS 'Highest portfolio value achieved - fees only charged above this level';
COMMENT ON COLUMN performance_fee_configs.hurdle_rate IS 'Minimum return percentage before performance fees apply';
COMMENT ON COLUMN trade_records.pnl IS 'Profit and Loss in base currency';
COMMENT ON COLUMN trade_records.performance_fee IS 'Performance fee charged for this trade';
COMMENT ON COLUMN performance_summaries.sharpe_ratio IS 'Risk-adjusted return metric';
COMMENT ON COLUMN performance_summaries.profit_factor IS 'Gross profit divided by gross loss';
