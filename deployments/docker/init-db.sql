-- AI Agentic Crypto Browser Database Initialization

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create schemas
CREATE SCHEMA IF NOT EXISTS ai_data;
CREATE SCHEMA IF NOT EXISTS market_data;
CREATE SCHEMA IF NOT EXISTS user_data;
CREATE SCHEMA IF NOT EXISTS analytics;

-- User behavior tracking
CREATE TABLE IF NOT EXISTS user_data.behavior_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    context JSONB,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_behavior_events_user_id ON user_data.behavior_events(user_id);
CREATE INDEX IF NOT EXISTS idx_behavior_events_timestamp ON user_data.behavior_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_behavior_events_type ON user_data.behavior_events(event_type);

-- User profiles
CREATE TABLE IF NOT EXISTS user_data.user_profiles (
    user_id UUID PRIMARY KEY,
    profile_data JSONB NOT NULL,
    preferences JSONB,
    behavior_patterns JSONB,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Market patterns
CREATE TABLE IF NOT EXISTS market_data.detected_patterns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pattern_type VARCHAR(50) NOT NULL,
    asset VARCHAR(20) NOT NULL,
    timeframe VARCHAR(10) NOT NULL,
    confidence DECIMAL(5,4) NOT NULL,
    strength DECIMAL(5,4) NOT NULL,
    characteristics JSONB NOT NULL,
    market_context JSONB,
    first_detected TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    occurrence_count INTEGER DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_patterns_asset ON market_data.detected_patterns(asset);
CREATE INDEX IF NOT EXISTS idx_patterns_type ON market_data.detected_patterns(pattern_type);
CREATE INDEX IF NOT EXISTS idx_patterns_confidence ON market_data.detected_patterns(confidence);
CREATE INDEX IF NOT EXISTS idx_patterns_timestamp ON market_data.detected_patterns(first_detected);

-- Adaptive strategies
CREATE TABLE IF NOT EXISTS market_data.adaptive_strategies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    strategy_type VARCHAR(50) NOT NULL,
    base_parameters JSONB NOT NULL,
    current_parameters JSONB NOT NULL,
    performance_targets JSONB,
    risk_limits JSONB,
    is_active BOOLEAN DEFAULT true,
    adaptation_count INTEGER DEFAULT 0,
    last_adaptation TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_strategies_type ON market_data.adaptive_strategies(strategy_type);
CREATE INDEX IF NOT EXISTS idx_strategies_active ON market_data.adaptive_strategies(is_active);

-- Strategy adaptations
CREATE TABLE IF NOT EXISTS market_data.strategy_adaptations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    strategy_id UUID NOT NULL REFERENCES market_data.adaptive_strategies(id),
    adaptation_type VARCHAR(50) NOT NULL,
    trigger_reason VARCHAR(100) NOT NULL,
    old_parameters JSONB NOT NULL,
    new_parameters JSONB NOT NULL,
    market_context JSONB,
    confidence DECIMAL(5,4) NOT NULL,
    success BOOLEAN,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_adaptations_strategy ON market_data.strategy_adaptations(strategy_id);
CREATE INDEX IF NOT EXISTS idx_adaptations_timestamp ON market_data.strategy_adaptations(timestamp);

-- Performance metrics
CREATE TABLE IF NOT EXISTS market_data.performance_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    strategy_id UUID NOT NULL REFERENCES market_data.adaptive_strategies(id),
    total_return DECIMAL(10,6),
    sharpe_ratio DECIMAL(8,4),
    max_drawdown DECIMAL(8,4),
    win_rate DECIMAL(5,4),
    profit_factor DECIMAL(8,4),
    total_trades INTEGER,
    metrics_data JSONB NOT NULL,
    calculation_date TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_performance_strategy ON market_data.performance_metrics(strategy_id);
CREATE INDEX IF NOT EXISTS idx_performance_date ON market_data.performance_metrics(calculation_date);

-- AI analysis results
CREATE TABLE IF NOT EXISTS ai_data.analysis_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id VARCHAR(100) NOT NULL,
    user_id UUID,
    analysis_type VARCHAR(50) NOT NULL,
    input_data JSONB NOT NULL,
    result_data JSONB NOT NULL,
    confidence DECIMAL(5,4),
    processing_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_analysis_request ON ai_data.analysis_results(request_id);
CREATE INDEX IF NOT EXISTS idx_analysis_user ON ai_data.analysis_results(user_id);
CREATE INDEX IF NOT EXISTS idx_analysis_type ON ai_data.analysis_results(analysis_type);
CREATE INDEX IF NOT EXISTS idx_analysis_timestamp ON ai_data.analysis_results(created_at);

-- System analytics
CREATE TABLE IF NOT EXISTS analytics.system_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(15,6) NOT NULL,
    metric_unit VARCHAR(20),
    tags JSONB,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_metrics_name ON analytics.system_metrics(metric_name);
CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON analytics.system_metrics(timestamp);

-- API usage tracking
CREATE TABLE IF NOT EXISTS analytics.api_usage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    endpoint VARCHAR(200) NOT NULL,
    method VARCHAR(10) NOT NULL,
    user_id UUID,
    response_status INTEGER NOT NULL,
    response_time_ms INTEGER NOT NULL,
    request_size_bytes INTEGER,
    response_size_bytes INTEGER,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_api_endpoint ON analytics.api_usage(endpoint);
CREATE INDEX IF NOT EXISTS idx_api_user ON analytics.api_usage(user_id);
CREATE INDEX IF NOT EXISTS idx_api_timestamp ON analytics.api_usage(timestamp);

-- Create views for common queries
CREATE OR REPLACE VIEW analytics.daily_pattern_summary AS
SELECT 
    DATE(first_detected) as date,
    pattern_type,
    asset,
    COUNT(*) as pattern_count,
    AVG(confidence) as avg_confidence,
    AVG(strength) as avg_strength
FROM market_data.detected_patterns
GROUP BY DATE(first_detected), pattern_type, asset
ORDER BY date DESC;

CREATE OR REPLACE VIEW analytics.strategy_performance_summary AS
SELECT 
    s.id,
    s.name,
    s.strategy_type,
    s.is_active,
    s.adaptation_count,
    pm.total_return,
    pm.sharpe_ratio,
    pm.max_drawdown,
    pm.win_rate,
    pm.total_trades,
    s.last_adaptation
FROM market_data.adaptive_strategies s
LEFT JOIN LATERAL (
    SELECT * FROM market_data.performance_metrics pm2 
    WHERE pm2.strategy_id = s.id 
    ORDER BY calculation_date DESC 
    LIMIT 1
) pm ON true;

-- Grant permissions
GRANT USAGE ON SCHEMA ai_data TO postgres;
GRANT USAGE ON SCHEMA market_data TO postgres;
GRANT USAGE ON SCHEMA user_data TO postgres;
GRANT USAGE ON SCHEMA analytics TO postgres;

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA ai_data TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA market_data TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA user_data TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA analytics TO postgres;

-- Insert initial data
INSERT INTO market_data.adaptive_strategies (name, strategy_type, base_parameters, current_parameters, performance_targets, risk_limits) VALUES
('Default Trend Following', 'trend_following', 
 '{"position_size": 0.05, "stop_loss": 0.02, "take_profit": 0.04}',
 '{"position_size": 0.05, "stop_loss": 0.02, "take_profit": 0.04}',
 '{"target_return": 0.15, "max_drawdown": 0.1, "min_sharpe_ratio": 1.0}',
 '{"max_position_size": 0.1, "max_leverage": 2.0, "stop_loss_percentage": 0.05}'),
('Default Mean Reversion', 'mean_reversion',
 '{"position_size": 0.03, "reversion_threshold": 2.0, "hold_time": 24}',
 '{"position_size": 0.03, "reversion_threshold": 2.0, "hold_time": 24}',
 '{"target_return": 0.12, "max_drawdown": 0.08, "min_sharpe_ratio": 1.2}',
 '{"max_position_size": 0.08, "max_leverage": 1.5, "stop_loss_percentage": 0.03}')
ON CONFLICT DO NOTHING;
