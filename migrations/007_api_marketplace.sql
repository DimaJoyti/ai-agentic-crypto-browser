-- API Marketplace and Usage Tracking Schema
-- Migration 007: API marketplace and usage-based billing

-- API Keys table for marketplace authentication
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    key_hash VARCHAR(255) NOT NULL UNIQUE, -- Hashed API key for security
    key_prefix VARCHAR(20) NOT NULL, -- First few characters for display
    name VARCHAR(255) NOT NULL,
    description TEXT,
    scopes TEXT[] DEFAULT '{}', -- Array of allowed scopes/endpoints
    rate_limit_override INTEGER, -- Custom rate limit if different from default
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    metadata JSONB DEFAULT '{}'
);

-- API Usage Records table for billing and analytics
CREATE TABLE IF NOT EXISTS api_usage_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    api_key_id UUID REFERENCES api_keys(id) ON DELETE SET NULL,
    api_key_prefix VARCHAR(20), -- For quick identification
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    request_size BIGINT DEFAULT 0,
    response_size BIGINT DEFAULT 0,
    duration_ms INTEGER NOT NULL, -- Duration in milliseconds
    cost DECIMAL(10,6) NOT NULL, -- Cost in USD
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    success BOOLEAN DEFAULT true,
    status_code INTEGER,
    error_code VARCHAR(50),
    error_message TEXT,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB DEFAULT '{}'
);

-- API Usage Summary table for faster billing queries
CREATE TABLE IF NOT EXISTS api_usage_summary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    billing_period DATE NOT NULL, -- First day of billing month
    endpoint VARCHAR(255) NOT NULL,
    total_requests BIGINT DEFAULT 0,
    total_cost DECIMAL(12,6) DEFAULT 0,
    total_duration_ms BIGINT DEFAULT 0,
    avg_duration_ms DECIMAL(10,2) DEFAULT 0,
    success_count BIGINT DEFAULT 0,
    error_count BIGINT DEFAULT 0,
    success_rate DECIMAL(5,4) DEFAULT 0, -- Success rate as decimal (0.9500 = 95%)
    first_request_at TIMESTAMP WITH TIME ZONE,
    last_request_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, billing_period, endpoint)
);

-- API Billing Records table for invoicing
CREATE TABLE IF NOT EXISTS api_billing_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    billing_period DATE NOT NULL,
    total_requests BIGINT DEFAULT 0,
    total_cost DECIMAL(12,6) DEFAULT 0,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(12,6) DEFAULT 0,
    final_amount DECIMAL(12,6) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(50) DEFAULT 'pending', -- pending, paid, failed, refunded
    stripe_invoice_id VARCHAR(255),
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    paid_at TIMESTAMP WITH TIME ZONE,
    due_date TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}'
);

-- Rate Limiting table for tracking current usage
CREATE TABLE IF NOT EXISTS api_rate_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    window_start TIMESTAMP WITH TIME ZONE NOT NULL,
    window_duration_seconds INTEGER NOT NULL DEFAULT 60, -- 1 minute default
    request_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, endpoint, window_start)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_api_keys_active ON api_keys(is_active) WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_api_usage_records_user_id ON api_usage_records(user_id);
CREATE INDEX IF NOT EXISTS idx_api_usage_records_timestamp ON api_usage_records(timestamp);
CREATE INDEX IF NOT EXISTS idx_api_usage_records_endpoint ON api_usage_records(endpoint);
CREATE INDEX IF NOT EXISTS idx_api_usage_records_billing_period ON api_usage_records(DATE_TRUNC('month', timestamp));
CREATE INDEX IF NOT EXISTS idx_api_usage_records_user_period ON api_usage_records(user_id, DATE_TRUNC('month', timestamp));

CREATE INDEX IF NOT EXISTS idx_api_usage_summary_user_period ON api_usage_summary(user_id, billing_period);
CREATE INDEX IF NOT EXISTS idx_api_usage_summary_period ON api_usage_summary(billing_period);

CREATE INDEX IF NOT EXISTS idx_api_billing_records_user_id ON api_billing_records(user_id);
CREATE INDEX IF NOT EXISTS idx_api_billing_records_period ON api_billing_records(billing_period);
CREATE INDEX IF NOT EXISTS idx_api_billing_records_status ON api_billing_records(status);

CREATE INDEX IF NOT EXISTS idx_api_rate_limits_user_endpoint ON api_rate_limits(user_id, endpoint);
CREATE INDEX IF NOT EXISTS idx_api_rate_limits_window ON api_rate_limits(window_start);

-- Functions for automatic summary updates
CREATE OR REPLACE FUNCTION update_api_usage_summary()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO api_usage_summary (
        user_id, billing_period, endpoint, total_requests, total_cost,
        total_duration_ms, success_count, error_count, first_request_at, last_request_at
    )
    VALUES (
        NEW.user_id,
        DATE_TRUNC('month', NEW.timestamp)::DATE,
        NEW.endpoint,
        1,
        NEW.cost,
        EXTRACT(EPOCH FROM NEW.duration_ms) * 1000,
        CASE WHEN NEW.success THEN 1 ELSE 0 END,
        CASE WHEN NEW.success THEN 0 ELSE 1 END,
        NEW.timestamp,
        NEW.timestamp
    )
    ON CONFLICT (user_id, billing_period, endpoint)
    DO UPDATE SET
        total_requests = api_usage_summary.total_requests + 1,
        total_cost = api_usage_summary.total_cost + NEW.cost,
        total_duration_ms = api_usage_summary.total_duration_ms + (EXTRACT(EPOCH FROM NEW.duration_ms) * 1000),
        success_count = api_usage_summary.success_count + CASE WHEN NEW.success THEN 1 ELSE 0 END,
        error_count = api_usage_summary.error_count + CASE WHEN NEW.success THEN 0 ELSE 1 END,
        last_request_at = GREATEST(api_usage_summary.last_request_at, NEW.timestamp),
        first_request_at = LEAST(api_usage_summary.first_request_at, NEW.timestamp),
        avg_duration_ms = (api_usage_summary.total_duration_ms + (EXTRACT(EPOCH FROM NEW.duration_ms) * 1000)) / (api_usage_summary.total_requests + 1),
        success_rate = (api_usage_summary.success_count + CASE WHEN NEW.success THEN 1 ELSE 0 END)::DECIMAL / (api_usage_summary.total_requests + 1),
        updated_at = NOW();

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update summary on new usage records
CREATE TRIGGER trigger_update_api_usage_summary
    AFTER INSERT ON api_usage_records
    FOR EACH ROW
    EXECUTE FUNCTION update_api_usage_summary();

-- Function to update API key last used timestamp
CREATE OR REPLACE FUNCTION update_api_key_last_used()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE api_keys 
    SET last_used_at = NEW.timestamp
    WHERE id = NEW.api_key_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update API key last used timestamp
CREATE TRIGGER trigger_update_api_key_last_used
    AFTER INSERT ON api_usage_records
    FOR EACH ROW
    WHEN (NEW.api_key_id IS NOT NULL)
    EXECUTE FUNCTION update_api_key_last_used();

-- Function to clean up old rate limit records
CREATE OR REPLACE FUNCTION cleanup_old_rate_limits()
RETURNS void AS $$
BEGIN
    DELETE FROM api_rate_limits 
    WHERE window_start < NOW() - INTERVAL '1 hour';
END;
$$ LANGUAGE plpgsql;

-- Insert default API pricing configuration
INSERT INTO api_usage_summary (user_id, billing_period, endpoint, total_requests, total_cost) 
VALUES ('system', '2024-01-01', 'pricing_config', 0, 0)
ON CONFLICT DO NOTHING;

-- Create view for API analytics
CREATE OR REPLACE VIEW api_analytics AS
SELECT 
    user_id,
    billing_period,
    SUM(total_requests) as total_requests,
    SUM(total_cost) as total_cost,
    AVG(avg_duration_ms) as avg_duration_ms,
    AVG(success_rate) as avg_success_rate,
    COUNT(DISTINCT endpoint) as unique_endpoints,
    MIN(first_request_at) as first_request_at,
    MAX(last_request_at) as last_request_at
FROM api_usage_summary
GROUP BY user_id, billing_period;

-- Create view for top API users
CREATE OR REPLACE VIEW top_api_users AS
SELECT 
    user_id,
    SUM(total_requests) as total_requests,
    SUM(total_cost) as total_revenue,
    COUNT(DISTINCT billing_period) as active_months,
    AVG(success_rate) as avg_success_rate,
    MAX(last_request_at) as last_activity
FROM api_usage_summary
WHERE billing_period >= DATE_TRUNC('month', NOW() - INTERVAL '12 months')
GROUP BY user_id
ORDER BY total_revenue DESC;

-- Comments for documentation
COMMENT ON TABLE api_keys IS 'API keys for marketplace authentication and authorization';
COMMENT ON TABLE api_usage_records IS 'Detailed records of all API usage for billing and analytics';
COMMENT ON TABLE api_usage_summary IS 'Aggregated usage summary by user, period, and endpoint for faster billing queries';
COMMENT ON TABLE api_billing_records IS 'Monthly billing records and invoice tracking';
COMMENT ON TABLE api_rate_limits IS 'Rate limiting tracking for API endpoints';

COMMENT ON COLUMN api_keys.key_hash IS 'SHA-256 hash of the API key for secure storage';
COMMENT ON COLUMN api_keys.key_prefix IS 'First 8 characters of API key for display purposes';
COMMENT ON COLUMN api_keys.scopes IS 'Array of allowed endpoint categories or specific endpoints';

COMMENT ON COLUMN api_usage_records.duration_ms IS 'Request duration in milliseconds';
COMMENT ON COLUMN api_usage_records.cost IS 'Cost of this API call in USD';

COMMENT ON COLUMN api_usage_summary.success_rate IS 'Success rate as decimal (0.95 = 95%)';
COMMENT ON COLUMN api_usage_summary.billing_period IS 'First day of the billing month';

-- Grant permissions (adjust as needed for your user roles)
-- GRANT SELECT, INSERT, UPDATE ON api_keys TO api_service;
-- GRANT SELECT, INSERT ON api_usage_records TO api_service;
-- GRANT SELECT, INSERT, UPDATE ON api_usage_summary TO api_service;
-- GRANT SELECT, INSERT, UPDATE ON api_billing_records TO billing_service;
