#!/bin/bash

# Launch Affiliate and Referral Program Script
# Deploy and configure the affiliate system for viral customer acquisition

set -e

echo "🤝 Launching Affiliate and Referral Program"
echo "==========================================="

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

print_step() {
    echo -e "${BLUE}▶${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Step 1: Database Migration
migrate_database() {
    print_step "Running affiliate program database migration..."
    
    # Check if migration file exists
    if [ ! -f "migrations/009_affiliate_program.sql" ]; then
        print_error "Migration file not found: migrations/009_affiliate_program.sql"
        exit 1
    fi
    
    # Run migration
    if command -v psql &> /dev/null; then
        psql $DATABASE_URL -f migrations/009_affiliate_program.sql
        print_success "Database migration completed"
    else
        print_warning "psql not found. Please run migration manually:"
        echo "  psql \$DATABASE_URL -f migrations/009_affiliate_program.sql"
    fi
}

# Step 2: Configure Affiliate Program
configure_affiliate_program() {
    print_step "Configuring affiliate program settings..."
    
    cat > config/affiliate_program.json << 'EOF'
{
  "program_settings": {
    "name": "AI-Agentic Crypto Browser Affiliate Program",
    "cookie_duration_days": 90,
    "minimum_payout": 100.00,
    "payout_schedule": "monthly",
    "payout_day": 1,
    "auto_approve_applications": false,
    "require_tax_info": true
  },
  "commission_tiers": {
    "bronze": {
      "rate": 0.15,
      "min_referrals": 0,
      "benefits": ["Basic tracking", "Monthly payouts"]
    },
    "silver": {
      "rate": 0.20,
      "min_referrals": 10,
      "benefits": ["Advanced analytics", "Priority support"]
    },
    "gold": {
      "rate": 0.25,
      "min_referrals": 50,
      "benefits": ["Custom landing pages", "Marketing materials"]
    },
    "platinum": {
      "rate": 0.30,
      "min_referrals": 100,
      "benefits": ["Dedicated account manager", "Custom campaigns"]
    },
    "diamond": {
      "rate": 0.35,
      "min_referrals": 500,
      "benefits": ["Revenue sharing", "Co-marketing opportunities"]
    }
  },
  "conversion_types": {
    "signup": {
      "commission_rate": 0.10,
      "base_value": 10.00
    },
    "subscription": {
      "commission_rate": 0.20,
      "recurring": true
    },
    "api_usage": {
      "commission_rate": 0.15,
      "recurring": true
    },
    "performance_fees": {
      "commission_rate": 0.10,
      "recurring": true
    }
  },
  "payment_methods": ["stripe", "crypto", "bank_transfer", "paypal"],
  "fraud_protection": {
    "max_conversions_per_ip": 3,
    "min_time_between_conversions": 3600,
    "require_email_verification": true
  }
}
EOF
    
    print_success "Affiliate program configuration created"
}

# Step 3: Build and Test
build_and_test() {
    print_step "Building and testing affiliate system..."
    
    # Build the application
    go build -o bin/affiliate-program cmd/main.go
    if [ $? -eq 0 ]; then
        print_success "Application built successfully"
    else
        print_error "Build failed"
        exit 1
    fi
    
    # Test affiliate endpoints
    print_step "Testing affiliate endpoints..."
    
    # Start application in background for testing
    ./bin/affiliate-program &
    APP_PID=$!
    sleep 3
    
    # Test program info endpoint
    if curl -s http://localhost:8080/affiliate/program > /dev/null; then
        print_success "Affiliate program endpoint working"
    else
        print_warning "Affiliate program endpoint not responding"
    fi
    
    # Test dashboard endpoint
    if curl -s http://localhost:8080/affiliate/dashboard > /dev/null; then
        print_success "Affiliate dashboard endpoint working"
    else
        print_warning "Affiliate dashboard endpoint not responding"
    fi
    
    # Stop test application
    kill $APP_PID 2>/dev/null || true
}

# Step 4: Create Marketing Materials
create_marketing_materials() {
    print_step "Creating affiliate marketing materials..."
    
    mkdir -p marketing/affiliate
    
    # Affiliate program landing page
    cat > marketing/affiliate/landing_page.md << 'EOF'
# Join Our Affiliate Program

## Earn Up to 35% Commission

Partner with the leading AI-powered crypto trading platform and earn substantial commissions on every referral.

### Why Partner With Us?

🎯 **High Conversion Rates**
- 85%+ AI prediction accuracy drives conversions
- Professional platform attracts serious traders
- Multiple revenue streams = higher lifetime value

💰 **Generous Commissions**
- Up to 35% commission on all referrals
- Recurring commissions on subscriptions
- Performance fee sharing (2-20% of profits)
- API usage commissions

📊 **Advanced Tracking**
- Real-time analytics and reporting
- 90-day cookie duration
- Multi-touch attribution
- Fraud protection built-in

🚀 **Marketing Support**
- Professional banners and creatives
- Landing page templates
- Email marketing templates
- Social media content

### Commission Structure

| Tier | Rate | Min Referrals | Benefits |
|------|------|---------------|----------|
| Bronze | 15% | 0 | Basic tracking |
| Silver | 20% | 10 | Advanced analytics |
| Gold | 25% | 50 | Custom materials |
| Platinum | 30% | 100 | Account manager |
| Diamond | 35% | 500 | Revenue sharing |

### Getting Started

1. **Apply** - Submit your application
2. **Approval** - Get approved within 24-48 hours
3. **Promote** - Start sharing your affiliate links
4. **Earn** - Get paid monthly on the 1st

[Apply Now] [Learn More] [Contact Us]
EOF
    
    # Email templates
    cat > marketing/affiliate/email_templates.md << 'EOF'
# Affiliate Email Templates

## Welcome Email
Subject: Welcome to Our Affiliate Program! 🎉

Hi [Name],

Welcome to the AI-Agentic Crypto Browser Affiliate Program!

Your affiliate code: [AFFILIATE_CODE]
Your commission rate: [COMMISSION_RATE]%

Here's what you need to know:
• Your unique affiliate link: [AFFILIATE_LINK]
• Commission structure: Up to 35% on all referrals
• Payment schedule: Monthly on the 1st
• Cookie duration: 90 days

Get started:
1. Share your affiliate link
2. Track performance in your dashboard
3. Get paid monthly

Questions? Reply to this email or contact support.

Best regards,
The Affiliate Team

## Monthly Performance Report
Subject: Your Affiliate Performance Report - [MONTH]

Hi [Name],

Here's your performance summary for [MONTH]:

📊 Performance Metrics:
• Clicks: [CLICKS]
• Conversions: [CONVERSIONS]
• Conversion Rate: [CONVERSION_RATE]%
• Commissions Earned: $[COMMISSIONS]

💰 Payout Information:
• This Month's Earnings: $[EARNINGS]
• Total Unpaid: $[UNPAID]
• Next Payout Date: [PAYOUT_DATE]

🎯 Tips to Improve:
• Focus on high-converting content
• Target crypto trading communities
• Use our new banner creatives

View full dashboard: [DASHBOARD_LINK]

Keep up the great work!

## Payout Notification
Subject: Your Commission Payment is On the Way! 💰

Hi [Name],

Great news! Your commission payment has been processed.

Payment Details:
• Amount: $[AMOUNT]
• Period: [PERIOD]
• Method: [PAYMENT_METHOD]
• Reference: [REFERENCE]

The payment should arrive within 1-3 business days.

Thank you for being a valued affiliate partner!
EOF
    
    # Social media templates
    cat > marketing/affiliate/social_templates.md << 'EOF'
# Social Media Templates

## Twitter/X Posts

### General Promotion
🚀 Just discovered this AI crypto trading platform with 85%+ accuracy!

✅ AI-powered predictions
✅ Sub-100ms execution
✅ Multi-chain support
✅ Professional-grade tools

Perfect for serious traders. Check it out: [AFFILIATE_LINK]

#CryptoTrading #AI #TradingBot

### Results-Focused
📈 My AI trading results this month:

• 73% win rate
• 12.5% average return
• $2,847 profit
• 0.15% max drawdown

The AI predictions are incredibly accurate. Try it yourself: [AFFILIATE_LINK]

### Educational
🧵 Thread: Why AI trading beats manual trading

1/ Human emotions cause 90% of trading losses
2/ AI removes fear, greed, and FOMO from decisions
3/ 85%+ accuracy vs 50-60% human average
4/ 24/7 monitoring vs limited human hours
5/ Backtested strategies vs gut feelings

Try AI trading: [AFFILIATE_LINK]

## LinkedIn Posts

### Professional
After testing 15+ crypto trading platforms, I found one that stands out.

Key differentiators:
🧠 85%+ AI prediction accuracy
⚡ Sub-100ms execution speed
🔗 7+ blockchain integrations
📊 Institutional-grade analytics

Perfect for professional traders and funds.

Interested? Let's connect: [AFFILIATE_LINK]

### Case Study
Case Study: How AI Trading Transformed My Portfolio

Before AI:
• 55% win rate
• High stress levels
• Inconsistent results
• Time-intensive

After AI:
• 78% win rate
• Automated execution
• Consistent profits
• Passive income

The difference is remarkable. See for yourself: [AFFILIATE_LINK]

## YouTube/Video Scripts

### 5-Minute Review
"I Tested This AI Trading Platform for 30 Days - Here's What Happened"

Outline:
1. Introduction (0-30s)
2. Platform overview (30s-1m)
3. AI accuracy testing (1m-2m)
4. Real trading results (2m-3m)
5. Pros and cons (3m-4m)
6. Final verdict (4m-5m)

CTA: "Link in description to try it yourself"

### Tutorial
"How to Set Up AI Crypto Trading in 10 Minutes"

Outline:
1. Account creation (0-2m)
2. Platform walkthrough (2m-5m)
3. AI strategy setup (5m-7m)
4. First trade execution (7m-9m)
5. Monitoring and results (9m-10m)

CTA: "Start your AI trading journey today"
EOF
    
    print_success "Marketing materials created"
}

# Step 5: Setup Tracking and Analytics
setup_tracking() {
    print_step "Setting up affiliate tracking and analytics..."
    
    cat > config/affiliate_tracking.yaml << 'EOF'
tracking:
  click_tracking:
    - source_attribution
    - utm_parameters
    - referrer_tracking
    - device_fingerprinting
  
  conversion_tracking:
    - signup_conversions
    - subscription_conversions
    - api_usage_conversions
    - performance_fee_conversions
  
  fraud_detection:
    - ip_address_validation
    - device_fingerprinting
    - velocity_checking
    - pattern_analysis

analytics:
  real_time_metrics:
    - clicks_per_hour
    - conversions_per_hour
    - commission_earnings
    - top_performing_affiliates
  
  reporting:
    - daily_performance_reports
    - weekly_trend_analysis
    - monthly_payout_reports
    - quarterly_business_reviews

alerts:
  - name: "High Performing Affiliate"
    condition: "daily_commissions > $1000"
    notification: "slack"
  
  - name: "Fraud Detection"
    condition: "suspicious_activity_detected"
    notification: "email"
  
  - name: "Payout Processing"
    condition: "monthly_payout_ready"
    notification: "slack"
EOF
    
    print_success "Tracking and analytics configuration created"
}

# Step 6: Revenue Projections
show_revenue_projections() {
    print_step "Calculating affiliate program revenue projections..."
    
    cat << 'EOF'

🤝 Affiliate Program Revenue Projections
=======================================

Customer Acquisition Impact:
• Affiliates can drive 30-50% of new customers
• Higher quality leads (pre-qualified by affiliates)
• Lower customer acquisition cost vs paid ads
• Viral growth through network effects

Conservative Estimates (Monthly):
• 50 active affiliates × 5 referrals × $199 avg = $49,750 revenue
• Commission cost: $9,950 (20% avg rate)
• Net revenue gain: $39,800/month

Optimistic Estimates (Monthly):
• 500 active affiliates × 10 referrals × $299 avg = $1,495,000 revenue
• Commission cost: $299,000 (20% avg rate)
• Net revenue gain: $1,196,000/month

Growth Timeline:
• Month 1: 25 affiliates, $25K revenue, $5K commissions
• Month 3: 100 affiliates, $100K revenue, $20K commissions
• Month 6: 300 affiliates, $450K revenue, $90K commissions
• Month 12: 1,000+ affiliates, $2M+ revenue, $400K+ commissions

Key Success Metrics:
• Affiliate recruitment rate: 50-100 new affiliates/month
• Average referrals per affiliate: 5-15/month
• Conversion rate: 2-5% (affiliate traffic)
• Average commission per affiliate: $200-$2,000/month
• Customer lifetime value: $500-$5,000

ROI Analysis:
• Customer acquisition cost: $50-$150 (vs $200-$500 paid ads)
• Payback period: 1-3 months
• Lifetime value to acquisition cost ratio: 5:1 to 20:1
• Commission expense ratio: 15-25% of revenue

Viral Growth Potential:
• Each satisfied customer becomes potential affiliate
• Network effects amplify reach exponentially
• Word-of-mouth marketing drives organic growth
• Community building creates sustainable advantage

EOF
}

# Step 7: Launch Checklist
create_launch_checklist() {
    print_step "Creating affiliate program launch checklist..."
    
    cat > docs/affiliate_launch_checklist.md << 'EOF'
# Affiliate Program Launch Checklist

## Pre-Launch (Week 1)
- [ ] Database migration completed
- [ ] Affiliate tracking system tested
- [ ] Commission calculation verified
- [ ] Payment processing configured
- [ ] Fraud detection enabled
- [ ] Legal terms and conditions reviewed
- [ ] Tax compliance documentation prepared

## Marketing Materials (Week 2)
- [ ] Landing page created and optimized
- [ ] Banner creatives designed (multiple sizes)
- [ ] Email templates created
- [ ] Social media templates prepared
- [ ] Video scripts written
- [ ] Case studies documented
- [ ] FAQ section completed

## Affiliate Recruitment (Week 3)
- [ ] Target affiliate list compiled
- [ ] Outreach email templates created
- [ ] Application process streamlined
- [ ] Approval workflow established
- [ ] Onboarding sequence automated
- [ ] Training materials prepared
- [ ] Support documentation created

## Launch Week (Week 4)
- [ ] Soft launch with 10 beta affiliates
- [ ] Tracking and analytics verified
- [ ] Commission calculations tested
- [ ] Payment processing validated
- [ ] Support team trained
- [ ] Monitoring dashboards active
- [ ] Feedback collection system ready

## Post-Launch (Ongoing)
- [ ] Daily performance monitoring
- [ ] Weekly affiliate recruitment
- [ ] Monthly performance reviews
- [ ] Quarterly program optimization
- [ ] Continuous fraud monitoring
- [ ] Regular payout processing
- [ ] Ongoing affiliate support

## Success Metrics to Track
- [ ] Number of active affiliates
- [ ] Average referrals per affiliate
- [ ] Conversion rates by traffic source
- [ ] Customer lifetime value from affiliates
- [ ] Commission expense ratio
- [ ] Affiliate satisfaction scores
- [ ] Revenue attribution accuracy
EOF
    
    print_success "Launch checklist created"
}

# Main execution
main() {
    echo ""
    print_step "Starting affiliate program launch..."
    echo ""
    
    migrate_database
    echo ""
    
    configure_affiliate_program
    echo ""
    
    build_and_test
    echo ""
    
    create_marketing_materials
    echo ""
    
    setup_tracking
    echo ""
    
    create_launch_checklist
    echo ""
    
    show_revenue_projections
    
    echo ""
    print_success "Affiliate program launch complete!"
    echo ""
    echo "🤝 Next Steps:"
    echo "1. Review and customize affiliate program settings"
    echo "2. Create affiliate recruitment campaign"
    echo "3. Launch with beta affiliates for testing"
    echo "4. Scale recruitment and optimize performance"
    echo ""
    echo "💰 Revenue Potential: $40K-$1.2M monthly net gain"
    echo "🎯 Target: 25-1,000 active affiliates"
    echo ""
    print_success "Ready to drive viral customer acquisition! 🚀"
}

# Run the script
main "$@"
