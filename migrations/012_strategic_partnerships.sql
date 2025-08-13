-- Strategic Partnerships Schema
-- Migration 012: Comprehensive partnership management and integration system

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Partners table for strategic partnerships
CREATE TABLE IF NOT EXISTS partners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    partner_type VARCHAR(50) NOT NULL, -- exchange, defi_protocol, media, technology, financial
    category VARCHAR(20) NOT NULL, -- tier_1, tier_2, tier_3, strategic
    status VARCHAR(20) NOT NULL DEFAULT 'prospect', -- prospect, negotiating, active, paused, terminated
    contract_type VARCHAR(30) NOT NULL, -- revenue_share, integration, white_label, licensing
    description TEXT,
    website VARCHAR(500),
    logo_url VARCHAR(500),
    headquarters_country VARCHAR(3),
    founded_year INTEGER,
    employee_count VARCHAR(20), -- 1-10, 11-50, 51-200, 201-1000, 1000+
    funding_stage VARCHAR(20), -- seed, series_a, series_b, series_c, public
    total_funding DECIMAL(15,2) DEFAULT 0.00,
    annual_revenue DECIMAL(15,2) DEFAULT 0.00,
    customer_count INTEGER DEFAULT 0,
    market_cap DECIMAL(20,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Partner contacts for relationship management
CREATE TABLE IF NOT EXISTS partner_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    contact_type VARCHAR(20) NOT NULL, -- primary, technical, business, legal, support
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    linkedin_url VARCHAR(500),
    time_zone VARCHAR(50),
    language VARCHAR(10) DEFAULT 'en',
    is_decision_maker BOOLEAN DEFAULT false,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partnership agreements and contracts
CREATE TABLE IF NOT EXISTS partnership_agreements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    agreement_type VARCHAR(30) NOT NULL, -- msa, sow, amendment, nda
    status VARCHAR(20) NOT NULL DEFAULT 'draft', -- draft, negotiating, signed, active, expired, terminated
    title VARCHAR(255) NOT NULL,
    description TEXT,
    signed_date DATE,
    effective_date DATE,
    expiration_date DATE,
    auto_renewal BOOLEAN DEFAULT false,
    renewal_term_months INTEGER DEFAULT 12,
    notice_period_days INTEGER DEFAULT 30,
    termination_fee DECIMAL(10,2) DEFAULT 0.00,
    governing_law VARCHAR(100),
    legal_entity VARCHAR(255),
    document_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Revenue sharing configurations
CREATE TABLE IF NOT EXISTS revenue_sharing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    agreement_id UUID REFERENCES partnership_agreements(id),
    sharing_model VARCHAR(20) NOT NULL, -- percentage, fixed_fee, tiered, hybrid
    percentage_rate DECIMAL(5,4) DEFAULT 0.0000, -- 0-100% as decimal (e.g., 0.2000 = 20%)
    fixed_fee_monthly DECIMAL(10,2) DEFAULT 0.00,
    minimum_payment DECIMAL(10,2) DEFAULT 0.00,
    payment_schedule VARCHAR(20) DEFAULT 'monthly', -- monthly, quarterly, annual
    payment_method VARCHAR(20) DEFAULT 'bank_transfer', -- bank_transfer, crypto, check
    currency VARCHAR(3) DEFAULT 'USD',
    revenue_types TEXT[], -- subscription, trading_fees, api_usage, performance_fees
    exclusions TEXT[], -- excluded revenue streams
    reporting_frequency VARCHAR(20) DEFAULT 'monthly', -- daily, weekly, monthly
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tiered revenue sharing rates
CREATE TABLE IF NOT EXISTS revenue_sharing_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    revenue_sharing_id UUID NOT NULL REFERENCES revenue_sharing(id) ON DELETE CASCADE,
    tier_order INTEGER NOT NULL,
    min_amount DECIMAL(15,2) NOT NULL,
    max_amount DECIMAL(15,2), -- NULL for unlimited
    rate DECIMAL(5,4) NOT NULL, -- percentage as decimal
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partner integrations and technical details
CREATE TABLE IF NOT EXISTS partner_integrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    integration_type VARCHAR(30) NOT NULL, -- api, webhook, white_label, embedded, data_feed
    status VARCHAR(20) NOT NULL DEFAULT 'planned', -- planned, in_progress, testing, live, deprecated
    name VARCHAR(255) NOT NULL,
    description TEXT,
    api_base_url VARCHAR(500),
    webhook_url VARCHAR(500),
    authentication_type VARCHAR(20), -- api_key, oauth2, jwt, mutual_tls
    rate_limit_per_minute INTEGER DEFAULT 1000,
    data_format VARCHAR(10) DEFAULT 'json', -- json, xml, csv, binary
    is_real_time BOOLEAN DEFAULT false,
    requires_approval BOOLEAN DEFAULT true,
    documentation_url VARCHAR(500),
    test_environment_url VARCHAR(500),
    go_live_date DATE,
    last_health_check TIMESTAMP WITH TIME ZONE,
    health_status VARCHAR(20) DEFAULT 'unknown', -- healthy, degraded, down, unknown
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- API endpoints for partner integrations
CREATE TABLE IF NOT EXISTS integration_endpoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integration_id UUID NOT NULL REFERENCES partner_integrations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    endpoint_url VARCHAR(500) NOT NULL,
    http_method VARCHAR(10) NOT NULL, -- GET, POST, PUT, DELETE, PATCH
    purpose TEXT,
    rate_limit_per_minute INTEGER DEFAULT 100,
    requires_auth BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Webhook configurations
CREATE TABLE IF NOT EXISTS integration_webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integration_id UUID NOT NULL REFERENCES partner_integrations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    webhook_url VARCHAR(500) NOT NULL,
    events TEXT[] NOT NULL, -- array of event types
    secret_key VARCHAR(255),
    max_retries INTEGER DEFAULT 3,
    retry_interval_seconds INTEGER DEFAULT 60,
    is_active BOOLEAN DEFAULT true,
    last_triggered TIMESTAMP WITH TIME ZONE,
    success_count INTEGER DEFAULT 0,
    failure_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partnership performance metrics
CREATE TABLE IF NOT EXISTS partnership_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    metric_date DATE NOT NULL,
    total_revenue DECIMAL(15,2) DEFAULT 0.00,
    shared_revenue DECIMAL(15,2) DEFAULT 0.00,
    referred_users INTEGER DEFAULT 0,
    active_users INTEGER DEFAULT 0,
    api_requests INTEGER DEFAULT 0,
    integration_uptime DECIMAL(5,4) DEFAULT 1.0000, -- 0-1 as decimal
    avg_response_time_ms INTEGER DEFAULT 0,
    error_rate DECIMAL(5,4) DEFAULT 0.0000,
    customer_satisfaction DECIMAL(3,2) DEFAULT 0.00, -- 1-10 scale
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(partner_id, metric_date)
);

-- Revenue sharing payments tracking
CREATE TABLE IF NOT EXISTS revenue_sharing_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    revenue_sharing_id UUID NOT NULL REFERENCES revenue_sharing(id),
    payment_period_start DATE NOT NULL,
    payment_period_end DATE NOT NULL,
    total_revenue DECIMAL(15,2) NOT NULL,
    shared_amount DECIMAL(15,2) NOT NULL,
    payment_status VARCHAR(20) DEFAULT 'pending', -- pending, processing, paid, failed, disputed
    payment_method VARCHAR(20),
    payment_reference VARCHAR(255),
    payment_date DATE,
    currency VARCHAR(3) DEFAULT 'USD',
    exchange_rate DECIMAL(10,6) DEFAULT 1.000000,
    fees DECIMAL(10,2) DEFAULT 0.00,
    net_amount DECIMAL(15,2),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partner onboarding tracking
CREATE TABLE IF NOT EXISTS partner_onboarding (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    stage VARCHAR(30) NOT NULL DEFAULT 'initiated', -- initiated, documentation, legal, integration, testing, live
    progress_percentage DECIMAL(5,2) DEFAULT 0.00, -- 0-100%
    estimated_go_live DATE,
    actual_go_live DATE,
    project_manager VARCHAR(255),
    blockers TEXT[],
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Onboarding milestones
CREATE TABLE IF NOT EXISTS onboarding_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    onboarding_id UUID NOT NULL REFERENCES partner_onboarding(id) ON DELETE CASCADE,
    milestone_name VARCHAR(255) NOT NULL,
    description TEXT,
    due_date DATE,
    completed_date DATE,
    status VARCHAR(20) DEFAULT 'pending', -- pending, in_progress, completed, blocked, skipped
    owner VARCHAR(255),
    milestone_order INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partner compliance tracking
CREATE TABLE IF NOT EXISTS partner_compliance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    framework VARCHAR(50) NOT NULL, -- SOC2, ISO27001, GDPR, PCI_DSS, etc.
    certification_name VARCHAR(255),
    certification_authority VARCHAR(255),
    valid_from DATE,
    valid_until DATE,
    certificate_number VARCHAR(255),
    status VARCHAR(20) DEFAULT 'pending', -- pending, valid, expired, revoked
    audit_frequency VARCHAR(20), -- annual, quarterly, monthly
    last_audit_date DATE,
    next_audit_date DATE,
    document_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partner communication log
CREATE TABLE IF NOT EXISTS partner_communications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    contact_id UUID REFERENCES partner_contacts(id),
    communication_type VARCHAR(20) NOT NULL, -- email, call, meeting, demo, contract_review
    subject VARCHAR(255),
    summary TEXT,
    outcome VARCHAR(50), -- positive, negative, neutral, follow_up_needed
    next_steps TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    our_attendees TEXT[],
    their_attendees TEXT[],
    meeting_url VARCHAR(500),
    recording_url VARCHAR(500),
    documents TEXT[], -- URLs to related documents
    created_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partner opportunities pipeline
CREATE TABLE IF NOT EXISTS partner_opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID REFERENCES partners(id), -- NULL for prospects
    opportunity_name VARCHAR(255) NOT NULL,
    opportunity_type VARCHAR(30) NOT NULL, -- integration, revenue_share, white_label, acquisition
    stage VARCHAR(20) NOT NULL DEFAULT 'prospecting', -- prospecting, qualification, proposal, negotiation, closed_won, closed_lost
    estimated_value DECIMAL(15,2) DEFAULT 0.00,
    probability DECIMAL(5,2) DEFAULT 0.00, -- 0-100%
    expected_close_date DATE,
    actual_close_date DATE,
    lead_source VARCHAR(50), -- inbound, outbound, referral, event, cold_outreach
    assigned_to VARCHAR(255),
    description TEXT,
    competitive_situation TEXT,
    decision_criteria TEXT,
    key_stakeholders TEXT[],
    next_steps TEXT,
    lost_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_partners_type ON partners(partner_type);
CREATE INDEX IF NOT EXISTS idx_partners_category ON partners(category);
CREATE INDEX IF NOT EXISTS idx_partners_status ON partners(status);
CREATE INDEX IF NOT EXISTS idx_partners_created ON partners(created_at);

CREATE INDEX IF NOT EXISTS idx_partner_contacts_partner ON partner_contacts(partner_id);
CREATE INDEX IF NOT EXISTS idx_partner_contacts_type ON partner_contacts(contact_type);

CREATE INDEX IF NOT EXISTS idx_partnership_agreements_partner ON partnership_agreements(partner_id);
CREATE INDEX IF NOT EXISTS idx_partnership_agreements_status ON partnership_agreements(status);

CREATE INDEX IF NOT EXISTS idx_revenue_sharing_partner ON revenue_sharing(partner_id);
CREATE INDEX IF NOT EXISTS idx_revenue_sharing_active ON revenue_sharing(is_active);

CREATE INDEX IF NOT EXISTS idx_partner_integrations_partner ON partner_integrations(partner_id);
CREATE INDEX IF NOT EXISTS idx_partner_integrations_status ON partner_integrations(status);

CREATE INDEX IF NOT EXISTS idx_partnership_metrics_partner ON partnership_metrics(partner_id);
CREATE INDEX IF NOT EXISTS idx_partnership_metrics_date ON partnership_metrics(metric_date);

CREATE INDEX IF NOT EXISTS idx_revenue_payments_partner ON revenue_sharing_payments(partner_id);
CREATE INDEX IF NOT EXISTS idx_revenue_payments_status ON revenue_sharing_payments(payment_status);
CREATE INDEX IF NOT EXISTS idx_revenue_payments_period ON revenue_sharing_payments(payment_period_start, payment_period_end);

-- Functions for automatic calculations
CREATE OR REPLACE FUNCTION update_partnership_metrics()
RETURNS TRIGGER AS $$
BEGIN
    -- Update partner metrics when revenue sharing payments are made
    IF NEW.payment_status = 'paid' AND (OLD.payment_status IS NULL OR OLD.payment_status != 'paid') THEN
        INSERT INTO partnership_metrics (
            partner_id, metric_date, total_revenue, shared_revenue
        ) VALUES (
            NEW.partner_id, 
            NEW.payment_date,
            NEW.total_revenue,
            NEW.shared_amount
        )
        ON CONFLICT (partner_id, metric_date) 
        DO UPDATE SET 
            total_revenue = partnership_metrics.total_revenue + NEW.total_revenue,
            shared_revenue = partnership_metrics.shared_revenue + NEW.shared_amount;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update partnership metrics
CREATE TRIGGER trigger_update_partnership_metrics
    AFTER UPDATE ON revenue_sharing_payments
    FOR EACH ROW
    EXECUTE FUNCTION update_partnership_metrics();

-- Function to calculate revenue sharing amounts
CREATE OR REPLACE FUNCTION calculate_revenue_share(
    p_partner_id UUID,
    p_revenue_amount DECIMAL,
    p_revenue_type TEXT
) RETURNS DECIMAL AS $$
DECLARE
    v_sharing_config RECORD;
    v_tier_config RECORD;
    v_share_amount DECIMAL := 0;
BEGIN
    -- Get revenue sharing configuration
    SELECT * INTO v_sharing_config
    FROM revenue_sharing 
    WHERE partner_id = p_partner_id 
    AND is_active = true
    AND (revenue_types IS NULL OR p_revenue_type = ANY(revenue_types))
    AND (exclusions IS NULL OR p_revenue_type != ALL(exclusions))
    LIMIT 1;
    
    IF NOT FOUND THEN
        RETURN 0;
    END IF;
    
    -- Calculate based on sharing model
    CASE v_sharing_config.sharing_model
        WHEN 'percentage' THEN
            v_share_amount := p_revenue_amount * v_sharing_config.percentage_rate;
        
        WHEN 'fixed_fee' THEN
            v_share_amount := v_sharing_config.fixed_fee_monthly;
        
        WHEN 'tiered' THEN
            -- Find appropriate tier
            SELECT * INTO v_tier_config
            FROM revenue_sharing_tiers
            WHERE revenue_sharing_id = v_sharing_config.id
            AND p_revenue_amount >= min_amount
            AND (max_amount IS NULL OR p_revenue_amount <= max_amount)
            ORDER BY tier_order
            LIMIT 1;
            
            IF FOUND THEN
                v_share_amount := p_revenue_amount * v_tier_config.rate;
            END IF;
        
        ELSE
            v_share_amount := 0;
    END CASE;
    
    -- Apply minimum payment
    IF v_share_amount < v_sharing_config.minimum_payment THEN
        v_share_amount := v_sharing_config.minimum_payment;
    END IF;
    
    RETURN v_share_amount;
END;
$$ LANGUAGE plpgsql;

-- Insert sample partner types and categories
INSERT INTO partners (name, partner_type, category, status, contract_type, description, website) VALUES
('Binance', 'exchange', 'tier_1', 'active', 'integration', 'World largest cryptocurrency exchange', 'https://binance.com'),
('Uniswap', 'defi_protocol', 'tier_1', 'active', 'integration', 'Leading decentralized exchange protocol', 'https://uniswap.org'),
('CoinDesk', 'media', 'tier_2', 'active', 'revenue_share', 'Leading cryptocurrency news and media', 'https://coindesk.com'),
('Chainlink', 'technology', 'tier_1', 'active', 'integration', 'Decentralized oracle network', 'https://chain.link'),
('Galaxy Digital', 'financial', 'tier_1', 'negotiating', 'white_label', 'Institutional crypto services', 'https://galaxydigital.io')
ON CONFLICT DO NOTHING;

-- Create views for partnership analytics
CREATE OR REPLACE VIEW partnership_summary AS
SELECT 
    p.id,
    p.name,
    p.partner_type,
    p.category,
    p.status,
    COUNT(DISTINCT pa.id) as active_agreements,
    COUNT(DISTINCT pi.id) as active_integrations,
    COALESCE(SUM(pm.total_revenue), 0) as total_revenue_generated,
    COALESCE(SUM(pm.shared_revenue), 0) as total_revenue_shared,
    COALESCE(AVG(pm.customer_satisfaction), 0) as avg_satisfaction,
    p.created_at
FROM partners p
LEFT JOIN partnership_agreements pa ON p.id = pa.partner_id AND pa.status = 'active'
LEFT JOIN partner_integrations pi ON p.id = pi.partner_id AND pi.status = 'live'
LEFT JOIN partnership_metrics pm ON p.id = pm.partner_id
GROUP BY p.id, p.name, p.partner_type, p.category, p.status, p.created_at;

CREATE OR REPLACE VIEW revenue_sharing_summary AS
SELECT 
    p.name as partner_name,
    p.partner_type,
    rs.sharing_model,
    rs.percentage_rate,
    rs.payment_schedule,
    COUNT(rsp.id) as total_payments,
    COALESCE(SUM(rsp.total_revenue), 0) as total_revenue,
    COALESCE(SUM(rsp.shared_amount), 0) as total_shared,
    COALESCE(AVG(rsp.shared_amount), 0) as avg_payment
FROM partners p
JOIN revenue_sharing rs ON p.id = rs.partner_id
LEFT JOIN revenue_sharing_payments rsp ON rs.id = rsp.revenue_sharing_id AND rsp.payment_status = 'paid'
WHERE rs.is_active = true
GROUP BY p.name, p.partner_type, rs.sharing_model, rs.percentage_rate, rs.payment_schedule;

-- Comments for documentation
COMMENT ON TABLE partners IS 'Strategic partners and their basic information';
COMMENT ON TABLE partnership_agreements IS 'Legal agreements and contracts with partners';
COMMENT ON TABLE revenue_sharing IS 'Revenue sharing configurations and terms';
COMMENT ON TABLE partner_integrations IS 'Technical integrations with partners';
COMMENT ON TABLE partnership_metrics IS 'Partnership performance metrics and KPIs';

COMMENT ON COLUMN revenue_sharing.percentage_rate IS 'Revenue share percentage as decimal (0.2000 = 20%)';
COMMENT ON COLUMN partnership_metrics.integration_uptime IS 'Integration uptime as decimal (0.9999 = 99.99%)';
COMMENT ON COLUMN partnership_metrics.error_rate IS 'API error rate as decimal (0.0100 = 1%)';
