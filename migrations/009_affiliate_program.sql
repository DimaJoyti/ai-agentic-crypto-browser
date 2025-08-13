-- Affiliate and Referral Program Schema
-- Migration 009: Comprehensive affiliate tracking and commission system

-- Affiliates table for partner management
CREATE TABLE IF NOT EXISTS affiliates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL UNIQUE,
    affiliate_code VARCHAR(50) NOT NULL UNIQUE,
    affiliate_type VARCHAR(20) NOT NULL DEFAULT 'individual', -- individual, business, influencer, partner
    commission_rate DECIMAL(5,4) NOT NULL DEFAULT 0.2000, -- 20% default
    tier_id UUID REFERENCES commission_tiers(id),
    payment_method VARCHAR(20) NOT NULL DEFAULT 'stripe', -- stripe, crypto, bank, paypal
    payment_details TEXT, -- encrypted payment information
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, active, suspended, terminated
    total_referrals BIGINT DEFAULT 0,
    total_clicks BIGINT DEFAULT 0,
    total_signups BIGINT DEFAULT 0,
    total_conversions BIGINT DEFAULT 0,
    total_revenue DECIMAL(15,2) DEFAULT 0.00,
    total_commissions DECIMAL(15,2) DEFAULT 0.00,
    unpaid_commissions DECIMAL(15,2) DEFAULT 0.00,
    last_payout_at TIMESTAMP WITH TIME ZONE,
    approval_date TIMESTAMP WITH TIME ZONE,
    approved_by VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Commission tiers for different affiliate levels
CREATE TABLE IF NOT EXISTS commission_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tier_name VARCHAR(50) NOT NULL UNIQUE,
    min_referrals INTEGER NOT NULL DEFAULT 0,
    max_referrals INTEGER, -- NULL for unlimited
    commission_rate DECIMAL(5,4) NOT NULL,
    bonus_rate DECIMAL(5,4) DEFAULT 0.0000, -- Additional bonus percentage
    required_revenue DECIMAL(15,2) DEFAULT 0.00, -- Minimum revenue to qualify
    tier_benefits TEXT[], -- Array of benefits
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Referrals table for tracking individual referrals
CREATE TABLE IF NOT EXISTS referrals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    affiliate_id UUID NOT NULL REFERENCES affiliates(id) ON DELETE CASCADE,
    referred_user_id VARCHAR(255) NOT NULL,
    referral_code VARCHAR(50) NOT NULL,
    referral_source VARCHAR(50), -- link, social, email, direct, etc.
    campaign_id VARCHAR(100), -- Marketing campaign identifier
    conversion_type VARCHAR(30) NOT NULL, -- signup, subscription, purchase, api_usage
    conversion_value DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    commission_amount DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    commission_rate DECIMAL(5,4) NOT NULL,
    bonus_amount DECIMAL(15,2) DEFAULT 0.00, -- Additional bonuses
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, confirmed, paid, cancelled, disputed
    ip_address INET,
    user_agent TEXT,
    converted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    cancellation_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Affiliate clicks tracking for analytics
CREATE TABLE IF NOT EXISTS affiliate_clicks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    affiliate_id UUID NOT NULL REFERENCES affiliates(id) ON DELETE CASCADE,
    referral_code VARCHAR(50) NOT NULL,
    click_source VARCHAR(50), -- website, social, email, etc.
    campaign_id VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    referer_url TEXT,
    landing_page TEXT,
    converted BOOLEAN DEFAULT false,
    conversion_id UUID REFERENCES referrals(id),
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Affiliate stats for performance tracking
CREATE TABLE IF NOT EXISTS affiliate_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    affiliate_id UUID NOT NULL REFERENCES affiliates(id) ON DELETE CASCADE,
    period VARCHAR(20) NOT NULL, -- daily, weekly, monthly, quarterly, yearly
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    total_clicks BIGINT DEFAULT 0,
    total_signups BIGINT DEFAULT 0,
    total_conversions BIGINT DEFAULT 0,
    conversion_rate DECIMAL(5,4) DEFAULT 0.0000,
    total_revenue DECIMAL(15,2) DEFAULT 0.00,
    total_commissions DECIMAL(15,2) DEFAULT 0.00,
    average_order_value DECIMAL(15,2) DEFAULT 0.00,
    top_referral_source VARCHAR(50),
    top_conversion_type VARCHAR(30),
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(affiliate_id, period, period_start)
);

-- Commission payouts tracking
CREATE TABLE IF NOT EXISTS commission_payouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    affiliate_id UUID NOT NULL REFERENCES affiliates(id) ON DELETE CASCADE,
    payout_period_start DATE NOT NULL,
    payout_period_end DATE NOT NULL,
    total_commissions DECIMAL(15,2) NOT NULL,
    total_bonuses DECIMAL(15,2) DEFAULT 0.00,
    fees DECIMAL(15,2) DEFAULT 0.00, -- Payment processing fees
    net_payout DECIMAL(15,2) NOT NULL,
    payment_method VARCHAR(20) NOT NULL,
    payment_reference VARCHAR(255), -- Stripe payment ID, crypto tx hash, etc.
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed
    processed_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    failure_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Referral campaigns for organized marketing efforts
CREATE TABLE IF NOT EXISTS referral_campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_name VARCHAR(100) NOT NULL,
    campaign_code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    commission_rate DECIMAL(5,4), -- Override default rate
    bonus_rate DECIMAL(5,4) DEFAULT 0.0000,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    target_conversions INTEGER,
    max_budget DECIMAL(15,2),
    current_spend DECIMAL(15,2) DEFAULT 0.00,
    status VARCHAR(20) DEFAULT 'active', -- active, paused, completed, cancelled
    created_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Affiliate applications for new partner onboarding
CREATE TABLE IF NOT EXISTS affiliate_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    business_name VARCHAR(255),
    website_url VARCHAR(500),
    social_media_links JSONB DEFAULT '{}',
    marketing_experience TEXT,
    target_audience TEXT,
    promotional_methods TEXT,
    expected_referrals INTEGER,
    application_status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected
    reviewed_by VARCHAR(255),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_affiliates_user_id ON affiliates(user_id);
CREATE INDEX IF NOT EXISTS idx_affiliates_code ON affiliates(affiliate_code);
CREATE INDEX IF NOT EXISTS idx_affiliates_status ON affiliates(status);
CREATE INDEX IF NOT EXISTS idx_affiliates_type ON affiliates(affiliate_type);
CREATE INDEX IF NOT EXISTS idx_affiliates_commissions ON affiliates(total_commissions DESC);

CREATE INDEX IF NOT EXISTS idx_referrals_affiliate_id ON referrals(affiliate_id);
CREATE INDEX IF NOT EXISTS idx_referrals_referred_user ON referrals(referred_user_id);
CREATE INDEX IF NOT EXISTS idx_referrals_code ON referrals(referral_code);
CREATE INDEX IF NOT EXISTS idx_referrals_status ON referrals(status);
CREATE INDEX IF NOT EXISTS idx_referrals_converted_at ON referrals(converted_at);
CREATE INDEX IF NOT EXISTS idx_referrals_conversion_type ON referrals(conversion_type);

CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_affiliate_id ON affiliate_clicks(affiliate_id);
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_code ON affiliate_clicks(referral_code);
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_clicked_at ON affiliate_clicks(clicked_at);
CREATE INDEX IF NOT EXISTS idx_affiliate_clicks_converted ON affiliate_clicks(converted);

CREATE INDEX IF NOT EXISTS idx_affiliate_stats_affiliate_period ON affiliate_stats(affiliate_id, period, period_start);
CREATE INDEX IF NOT EXISTS idx_affiliate_stats_period_dates ON affiliate_stats(period, period_start, period_end);

CREATE INDEX IF NOT EXISTS idx_commission_payouts_affiliate_id ON commission_payouts(affiliate_id);
CREATE INDEX IF NOT EXISTS idx_commission_payouts_status ON commission_payouts(status);
CREATE INDEX IF NOT EXISTS idx_commission_payouts_period ON commission_payouts(payout_period_start, payout_period_end);

-- Functions for automatic calculations
CREATE OR REPLACE FUNCTION update_affiliate_stats_on_referral()
RETURNS TRIGGER AS $$
BEGIN
    -- Update affiliate totals
    UPDATE affiliates 
    SET total_referrals = total_referrals + 1,
        total_revenue = total_revenue + NEW.conversion_value,
        total_commissions = total_commissions + NEW.commission_amount,
        unpaid_commissions = unpaid_commissions + NEW.commission_amount,
        updated_at = NOW()
    WHERE id = NEW.affiliate_id;
    
    -- Update or insert daily stats
    INSERT INTO affiliate_stats (
        affiliate_id, period, period_start, period_end, total_conversions,
        total_revenue, total_commissions, last_updated
    )
    VALUES (
        NEW.affiliate_id, 'daily', DATE(NEW.converted_at), DATE(NEW.converted_at),
        1, NEW.conversion_value, NEW.commission_amount, NOW()
    )
    ON CONFLICT (affiliate_id, period, period_start)
    DO UPDATE SET
        total_conversions = affiliate_stats.total_conversions + 1,
        total_revenue = affiliate_stats.total_revenue + NEW.conversion_value,
        total_commissions = affiliate_stats.total_commissions + NEW.commission_amount,
        last_updated = NOW();
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update stats on new referrals
CREATE TRIGGER trigger_update_affiliate_stats
    AFTER INSERT ON referrals
    FOR EACH ROW
    EXECUTE FUNCTION update_affiliate_stats_on_referral();

-- Function to update click tracking
CREATE OR REPLACE FUNCTION update_affiliate_clicks()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE affiliates 
    SET total_clicks = total_clicks + 1,
        updated_at = NOW()
    WHERE id = NEW.affiliate_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update click counts
CREATE TRIGGER trigger_update_affiliate_clicks
    AFTER INSERT ON affiliate_clicks
    FOR EACH ROW
    EXECUTE FUNCTION update_affiliate_clicks();

-- Function to calculate conversion rates
CREATE OR REPLACE FUNCTION calculate_conversion_rates()
RETURNS void AS $$
BEGIN
    UPDATE affiliate_stats 
    SET conversion_rate = CASE 
        WHEN total_clicks > 0 THEN total_conversions::decimal / total_clicks::decimal
        ELSE 0 
    END,
    average_order_value = CASE 
        WHEN total_conversions > 0 THEN total_revenue / total_conversions
        ELSE 0 
    END
    WHERE last_updated >= NOW() - INTERVAL '1 day';
END;
$$ LANGUAGE plpgsql;

-- Insert default commission tiers
INSERT INTO commission_tiers (tier_name, min_referrals, commission_rate, description) VALUES
('Bronze', 0, 0.1500, 'Entry level - 15% commission on all referrals'),
('Silver', 10, 0.2000, 'Silver level - 20% commission after 10 referrals'),
('Gold', 50, 0.2500, 'Gold level - 25% commission after 50 referrals'),
('Platinum', 100, 0.3000, 'Platinum level - 30% commission after 100 referrals'),
('Diamond', 500, 0.3500, 'Diamond level - 35% commission after 500 referrals')
ON CONFLICT (tier_name) DO NOTHING;

-- Create views for analytics
CREATE OR REPLACE VIEW affiliate_performance AS
SELECT 
    a.id,
    a.affiliate_code,
    a.affiliate_type,
    a.status,
    a.total_referrals,
    a.total_revenue,
    a.total_commissions,
    a.unpaid_commissions,
    CASE WHEN a.total_clicks > 0 THEN 
        a.total_conversions::decimal / a.total_clicks::decimal 
    ELSE 0 END as conversion_rate,
    CASE WHEN a.total_conversions > 0 THEN 
        a.total_revenue / a.total_conversions 
    ELSE 0 END as average_order_value,
    ct.tier_name,
    ct.commission_rate as tier_rate
FROM affiliates a
LEFT JOIN commission_tiers ct ON a.tier_id = ct.id
WHERE a.status = 'active';

CREATE OR REPLACE VIEW top_affiliates AS
SELECT 
    affiliate_code,
    affiliate_type,
    total_referrals,
    total_revenue,
    total_commissions,
    CASE WHEN total_clicks > 0 THEN 
        total_conversions::decimal / total_clicks::decimal 
    ELSE 0 END as conversion_rate,
    ROW_NUMBER() OVER (ORDER BY total_commissions DESC) as rank
FROM affiliates 
WHERE status = 'active' AND total_referrals >= 5
ORDER BY total_commissions DESC;

-- Comments for documentation
COMMENT ON TABLE affiliates IS 'Affiliate partners and their performance metrics';
COMMENT ON TABLE referrals IS 'Individual referral records and commission tracking';
COMMENT ON TABLE affiliate_clicks IS 'Click tracking for affiliate link analytics';
COMMENT ON TABLE commission_payouts IS 'Commission payment records and status';
COMMENT ON TABLE commission_tiers IS 'Tiered commission structure based on performance';

COMMENT ON COLUMN affiliates.commission_rate IS 'Base commission rate as decimal (0.20 = 20%)';
COMMENT ON COLUMN affiliates.unpaid_commissions IS 'Total unpaid commission amount';
COMMENT ON COLUMN referrals.conversion_value IS 'Value of the conversion (subscription amount, etc.)';
COMMENT ON COLUMN referrals.commission_amount IS 'Commission earned from this referral';
