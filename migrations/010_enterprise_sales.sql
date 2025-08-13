-- Enterprise Sales Pipeline Schema
-- Migration 010: Comprehensive enterprise sales and CRM system

-- Enterprise leads table for prospect management
CREATE TABLE IF NOT EXISTS enterprise_leads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_name VARCHAR(255) NOT NULL,
    contact_name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(50),
    contact_title VARCHAR(100),
    company_size VARCHAR(20) NOT NULL, -- startup, small, medium, large, enterprise
    company_type VARCHAR(30) NOT NULL, -- hedge_fund, family_office, prop_trading, crypto_fund, institution
    aum DECIMAL(20,2) DEFAULT 0.00, -- Assets Under Management
    trading_volume DECIMAL(20,2) DEFAULT 0.00, -- Monthly trading volume
    current_solutions TEXT[], -- Existing trading platforms
    pain_points TEXT[], -- Key challenges
    budget DECIMAL(15,2) DEFAULT 0.00, -- Annual budget
    timeline VARCHAR(50), -- Implementation timeline
    decision_makers TEXT[], -- Key stakeholders
    source VARCHAR(50) NOT NULL, -- lead_source (website, referral, cold_outreach, event, etc.)
    status VARCHAR(20) NOT NULL DEFAULT 'new', -- new, contacted, qualified, proposal, negotiation, closed_won, closed_lost
    priority VARCHAR(10) NOT NULL DEFAULT 'medium', -- low, medium, high, critical
    assigned_sales_rep VARCHAR(255),
    estimated_value DECIMAL(15,2) DEFAULT 0.00, -- Deal size estimate
    probability_score DECIMAL(5,2) DEFAULT 0.00, -- 0-100%
    last_contact_date TIMESTAMP WITH TIME ZONE,
    next_follow_up_date TIMESTAMP WITH TIME ZONE,
    qualification_notes TEXT,
    competitor_info TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Enterprise deals table for active opportunities
CREATE TABLE IF NOT EXISTS enterprise_deals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lead_id UUID NOT NULL REFERENCES enterprise_leads(id) ON DELETE CASCADE,
    deal_name VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    deal_value DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    contract_length INTEGER DEFAULT 12, -- months
    stage VARCHAR(20) NOT NULL DEFAULT 'discovery', -- discovery, demo, proposal, negotiation, contract, closed_won, closed_lost
    probability DECIMAL(5,2) DEFAULT 0.00, -- 0-100%
    expected_close_date DATE,
    actual_close_date DATE,
    sales_rep VARCHAR(255) NOT NULL,
    sales_engineer VARCHAR(255),
    products TEXT[], -- List of products/services
    custom_pricing BOOLEAN DEFAULT false,
    white_label BOOLEAN DEFAULT false,
    on_premise BOOLEAN DEFAULT false,
    sla VARCHAR(50), -- Service level agreement tier
    support VARCHAR(50), -- Support tier
    lost_reason TEXT, -- If deal is lost
    competitor_info TEXT,
    contract_terms JSONB DEFAULT '{}',
    technical_requirements JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Sales activities table for tracking interactions
CREATE TABLE IF NOT EXISTS sales_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lead_id UUID REFERENCES enterprise_leads(id) ON DELETE CASCADE,
    deal_id UUID REFERENCES enterprise_deals(id) ON DELETE CASCADE,
    activity_type VARCHAR(30) NOT NULL, -- call, email, meeting, demo, proposal, contract, follow_up
    subject VARCHAR(255) NOT NULL,
    description TEXT,
    outcome VARCHAR(50), -- positive, negative, neutral, no_response
    next_steps TEXT,
    sales_rep VARCHAR(255) NOT NULL,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration INTEGER DEFAULT 0, -- minutes
    attendees TEXT[], -- List of attendees
    recording_url VARCHAR(500), -- Meeting recording link
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Proposals table for custom enterprise proposals
CREATE TABLE IF NOT EXISTS enterprise_proposals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deal_id UUID NOT NULL REFERENCES enterprise_deals(id) ON DELETE CASCADE,
    proposal_name VARCHAR(255) NOT NULL,
    version INTEGER DEFAULT 1,
    total_value DECIMAL(15,2) NOT NULL,
    contract_length INTEGER NOT NULL, -- months
    products JSONB NOT NULL DEFAULT '[]', -- Array of proposal products
    custom_features TEXT[],
    sla JSONB DEFAULT '{}', -- Service level agreement
    pricing JSONB DEFAULT '{}', -- Pricing structure
    terms JSONB DEFAULT '{}', -- Contract terms
    status VARCHAR(20) DEFAULT 'draft', -- draft, sent, viewed, accepted, rejected, expired
    sent_at TIMESTAMP WITH TIME ZONE,
    viewed_at TIMESTAMP WITH TIME ZONE,
    responded_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    document_url VARCHAR(500), -- PDF/document link
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Enterprise pricing tiers for different customer segments
CREATE TABLE IF NOT EXISTS enterprise_pricing_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tier_name VARCHAR(50) NOT NULL UNIQUE,
    min_aum DECIMAL(20,2) DEFAULT 0.00, -- Minimum AUM requirement
    max_aum DECIMAL(20,2), -- Maximum AUM (NULL for unlimited)
    base_monthly_fee DECIMAL(10,2) NOT NULL,
    setup_fee DECIMAL(10,2) DEFAULT 0.00,
    api_rate_limit INTEGER DEFAULT 1000, -- requests per minute
    performance_fee_rate DECIMAL(5,4) DEFAULT 0.0000, -- 0-100%
    volume_discount_rate DECIMAL(5,4) DEFAULT 0.0000,
    white_label_available BOOLEAN DEFAULT false,
    on_premise_available BOOLEAN DEFAULT false,
    dedicated_support BOOLEAN DEFAULT false,
    sla_uptime DECIMAL(5,4) DEFAULT 0.9900, -- 99.00%
    features JSONB DEFAULT '{}',
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Sales team members and territories
CREATE TABLE IF NOT EXISTS sales_team (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- sales_rep, sales_engineer, sales_manager, vp_sales
    territory VARCHAR(100), -- geographic or vertical territory
    quota DECIMAL(15,2) DEFAULT 0.00, -- Annual quota
    commission_rate DECIMAL(5,4) DEFAULT 0.0000, -- Commission percentage
    manager_id UUID REFERENCES sales_team(id),
    hire_date DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Sales targets and quotas
CREATE TABLE IF NOT EXISTS sales_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sales_rep_id UUID NOT NULL REFERENCES sales_team(id) ON DELETE CASCADE,
    target_period VARCHAR(20) NOT NULL, -- monthly, quarterly, annual
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    revenue_target DECIMAL(15,2) NOT NULL,
    deals_target INTEGER DEFAULT 0,
    new_leads_target INTEGER DEFAULT 0,
    activities_target INTEGER DEFAULT 0,
    actual_revenue DECIMAL(15,2) DEFAULT 0.00,
    actual_deals INTEGER DEFAULT 0,
    actual_leads INTEGER DEFAULT 0,
    actual_activities INTEGER DEFAULT 0,
    achievement_percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(sales_rep_id, target_period, period_start)
);

-- Lead scoring rules for qualification
CREATE TABLE IF NOT EXISTS lead_scoring_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_name VARCHAR(100) NOT NULL,
    criteria JSONB NOT NULL, -- Scoring criteria
    score_value INTEGER NOT NULL, -- Points to add/subtract
    rule_type VARCHAR(20) NOT NULL, -- demographic, behavioral, engagement
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Email templates for sales outreach
CREATE TABLE IF NOT EXISTS sales_email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_name VARCHAR(100) NOT NULL,
    template_type VARCHAR(30) NOT NULL, -- cold_outreach, follow_up, proposal, demo_invite
    subject_line VARCHAR(255) NOT NULL,
    email_body TEXT NOT NULL,
    variables JSONB DEFAULT '{}', -- Template variables
    open_rate DECIMAL(5,4) DEFAULT 0.0000,
    response_rate DECIMAL(5,4) DEFAULT 0.0000,
    usage_count INTEGER DEFAULT 0,
    created_by VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_enterprise_leads_status ON enterprise_leads(status);
CREATE INDEX IF NOT EXISTS idx_enterprise_leads_assigned_rep ON enterprise_leads(assigned_sales_rep);
CREATE INDEX IF NOT EXISTS idx_enterprise_leads_company_type ON enterprise_leads(company_type);
CREATE INDEX IF NOT EXISTS idx_enterprise_leads_priority ON enterprise_leads(priority);
CREATE INDEX IF NOT EXISTS idx_enterprise_leads_follow_up ON enterprise_leads(next_follow_up_date);
CREATE INDEX IF NOT EXISTS idx_enterprise_leads_created ON enterprise_leads(created_at);

CREATE INDEX IF NOT EXISTS idx_enterprise_deals_stage ON enterprise_deals(stage);
CREATE INDEX IF NOT EXISTS idx_enterprise_deals_sales_rep ON enterprise_deals(sales_rep);
CREATE INDEX IF NOT EXISTS idx_enterprise_deals_close_date ON enterprise_deals(expected_close_date);
CREATE INDEX IF NOT EXISTS idx_enterprise_deals_value ON enterprise_deals(deal_value);

CREATE INDEX IF NOT EXISTS idx_sales_activities_lead_id ON sales_activities(lead_id);
CREATE INDEX IF NOT EXISTS idx_sales_activities_deal_id ON sales_activities(deal_id);
CREATE INDEX IF NOT EXISTS idx_sales_activities_type ON sales_activities(activity_type);
CREATE INDEX IF NOT EXISTS idx_sales_activities_rep ON sales_activities(sales_rep);
CREATE INDEX IF NOT EXISTS idx_sales_activities_scheduled ON sales_activities(scheduled_at);

CREATE INDEX IF NOT EXISTS idx_enterprise_proposals_deal_id ON enterprise_proposals(deal_id);
CREATE INDEX IF NOT EXISTS idx_enterprise_proposals_status ON enterprise_proposals(status);

-- Functions for automatic calculations
CREATE OR REPLACE FUNCTION update_lead_score()
RETURNS TRIGGER AS $$
DECLARE
    total_score INTEGER := 0;
    rule RECORD;
BEGIN
    -- Calculate lead score based on scoring rules
    FOR rule IN SELECT * FROM lead_scoring_rules WHERE is_active = true LOOP
        -- Simplified scoring logic (would be more complex in real implementation)
        IF rule.rule_type = 'demographic' THEN
            IF NEW.company_type = 'hedge_fund' THEN
                total_score := total_score + 20;
            END IF;
            IF NEW.aum > 100000000 THEN -- $100M+
                total_score := total_score + 30;
            END IF;
        END IF;
    END LOOP;
    
    -- Update probability score based on total score
    NEW.probability_score := LEAST(total_score, 100);
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update lead scores
CREATE TRIGGER trigger_update_lead_score
    BEFORE INSERT OR UPDATE ON enterprise_leads
    FOR EACH ROW
    EXECUTE FUNCTION update_lead_score();

-- Function to update sales targets
CREATE OR REPLACE FUNCTION update_sales_targets()
RETURNS TRIGGER AS $$
BEGIN
    -- Update actual revenue when deal is closed won
    IF NEW.stage = 'closed_won' AND (OLD.stage IS NULL OR OLD.stage != 'closed_won') THEN
        UPDATE sales_targets 
        SET actual_revenue = actual_revenue + NEW.deal_value,
            actual_deals = actual_deals + 1,
            achievement_percentage = (actual_revenue + NEW.deal_value) / revenue_target * 100,
            updated_at = NOW()
        WHERE sales_rep_id = (
            SELECT id FROM sales_team WHERE user_id = NEW.sales_rep
        ) AND period_start <= NEW.actual_close_date 
          AND period_end >= NEW.actual_close_date;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update sales targets
CREATE TRIGGER trigger_update_sales_targets
    AFTER UPDATE ON enterprise_deals
    FOR EACH ROW
    EXECUTE FUNCTION update_sales_targets();

-- Insert default pricing tiers
INSERT INTO enterprise_pricing_tiers (tier_name, min_aum, base_monthly_fee, setup_fee, description) VALUES
('Startup Fund', 0, 5000.00, 2500.00, 'For emerging crypto funds and prop trading firms'),
('Growth Fund', 10000000, 15000.00, 7500.00, 'For established funds with $10M+ AUM'),
('Institutional', 100000000, 50000.00, 25000.00, 'For large institutions with $100M+ AUM'),
('Enterprise', 1000000000, 150000.00, 75000.00, 'For major institutions with $1B+ AUM')
ON CONFLICT (tier_name) DO NOTHING;

-- Insert default lead scoring rules
INSERT INTO lead_scoring_rules (rule_name, criteria, score_value, rule_type) VALUES
('Hedge Fund Type', '{"company_type": "hedge_fund"}', 25, 'demographic'),
('Large AUM', '{"aum_threshold": 100000000}', 30, 'demographic'),
('High Trading Volume', '{"volume_threshold": 10000000}', 20, 'behavioral'),
('Enterprise Size', '{"company_size": "enterprise"}', 15, 'demographic'),
('Crypto Fund', '{"company_type": "crypto_fund"}', 20, 'demographic')
ON CONFLICT DO NOTHING;

-- Insert default email templates
INSERT INTO sales_email_templates (template_name, template_type, subject_line, email_body) VALUES
('Cold Outreach - Hedge Fund', 'cold_outreach', 
'AI Trading Platform for {{company_name}} - 85%+ Accuracy',
'Hi {{contact_name}},

I noticed {{company_name}} is actively trading crypto assets. Our AI platform has helped similar {{company_type}}s achieve:

• 85%+ prediction accuracy
• 40% reduction in drawdowns  
• $2M+ additional alpha per $100M AUM

Would you be interested in a 15-minute demo to see how this could benefit {{company_name}}?

Best regards,
{{sales_rep_name}}'),

('Demo Follow-up', 'follow_up',
'Next steps for {{company_name}} - AI Trading Implementation',
'Hi {{contact_name}},

Thank you for the demo yesterday. Based on our discussion about {{pain_points}}, I believe our platform could deliver significant value for {{company_name}}.

Next steps:
1. Technical deep-dive with your team
2. Custom proposal with {{company_type}} pricing
3. Pilot program setup

When would be a good time for the technical session?

Best regards,
{{sales_rep_name}}')
ON CONFLICT DO NOTHING;

-- Create views for sales analytics
CREATE OR REPLACE VIEW sales_pipeline_summary AS
SELECT 
    s.name as sales_rep,
    COUNT(l.*) as total_leads,
    COUNT(CASE WHEN l.status = 'qualified' THEN 1 END) as qualified_leads,
    COUNT(d.*) as total_deals,
    COUNT(CASE WHEN d.stage = 'proposal' THEN 1 END) as proposal_stage,
    COUNT(CASE WHEN d.stage = 'negotiation' THEN 1 END) as negotiation_stage,
    COUNT(CASE WHEN d.stage = 'closed_won' THEN 1 END) as closed_won,
    COALESCE(SUM(d.deal_value), 0) as total_pipeline_value,
    COALESCE(SUM(CASE WHEN d.stage = 'closed_won' THEN d.deal_value ELSE 0 END), 0) as closed_revenue
FROM sales_team s
LEFT JOIN enterprise_leads l ON s.user_id = l.assigned_sales_rep
LEFT JOIN enterprise_deals d ON l.id = d.lead_id
WHERE s.is_active = true
GROUP BY s.id, s.name;

CREATE OR REPLACE VIEW enterprise_lead_funnel AS
SELECT 
    status,
    COUNT(*) as lead_count,
    COALESCE(SUM(estimated_value), 0) as total_value,
    COALESCE(AVG(probability_score), 0) as avg_probability
FROM enterprise_leads 
GROUP BY status
ORDER BY 
    CASE status 
        WHEN 'new' THEN 1
        WHEN 'contacted' THEN 2
        WHEN 'qualified' THEN 3
        WHEN 'proposal' THEN 4
        WHEN 'negotiation' THEN 5
        WHEN 'closed_won' THEN 6
        WHEN 'closed_lost' THEN 7
    END;

-- Comments for documentation
COMMENT ON TABLE enterprise_leads IS 'Enterprise sales leads and prospects';
COMMENT ON TABLE enterprise_deals IS 'Active enterprise sales opportunities';
COMMENT ON TABLE sales_activities IS 'Sales team activities and touchpoints';
COMMENT ON TABLE enterprise_proposals IS 'Custom proposals for enterprise clients';
COMMENT ON TABLE enterprise_pricing_tiers IS 'Pricing tiers for different enterprise segments';

COMMENT ON COLUMN enterprise_leads.aum IS 'Assets Under Management in USD';
COMMENT ON COLUMN enterprise_leads.probability_score IS 'Lead qualification score 0-100%';
COMMENT ON COLUMN enterprise_deals.deal_value IS 'Total contract value in USD';
COMMENT ON COLUMN enterprise_deals.probability IS 'Deal close probability 0-100%';
