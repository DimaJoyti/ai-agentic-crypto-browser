#!/bin/bash

# Setup Stripe Integration for AI-Agentic Crypto Browser
# This script helps you configure Stripe for immediate monetization

set -e

echo "💳 Setting up Stripe Integration"
echo "================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

# Check if .env file exists
check_env_file() {
    if [ ! -f .env ]; then
        print_info "Creating .env file from template..."
        cp .env.example .env 2>/dev/null || {
            print_warning ".env.example not found, creating new .env file"
            touch .env
        }
    fi
    print_status ".env file ready"
}

# Guide user through Stripe setup
setup_stripe_account() {
    echo ""
    print_info "Stripe Account Setup Instructions:"
    echo ""
    echo "1. Go to https://stripe.com and create an account"
    echo "2. Complete account verification (required for live payments)"
    echo "3. Go to Developers > API Keys in your Stripe dashboard"
    echo "4. Copy your keys and paste them below"
    echo ""
    
    # Get Stripe keys from user
    read -p "Enter your Stripe Publishable Key (pk_test_... or pk_live_...): " STRIPE_PUBLISHABLE_KEY
    read -p "Enter your Stripe Secret Key (sk_test_... or sk_live_...): " STRIPE_SECRET_KEY
    
    # Validate keys
    if [[ ! $STRIPE_PUBLISHABLE_KEY =~ ^pk_(test_|live_) ]]; then
        print_error "Invalid publishable key format"
        exit 1
    fi
    
    if [[ ! $STRIPE_SECRET_KEY =~ ^sk_(test_|live_) ]]; then
        print_error "Invalid secret key format"
        exit 1
    fi
    
    print_status "Stripe keys validated"
}

# Update .env file with Stripe keys
update_env_file() {
    print_info "Updating .env file with Stripe configuration..."
    
    # Remove existing Stripe keys if they exist
    sed -i '/^STRIPE_/d' .env
    
    # Add new Stripe configuration
    cat >> .env << EOF

# Stripe Configuration
STRIPE_PUBLISHABLE_KEY=$STRIPE_PUBLISHABLE_KEY
STRIPE_SECRET_KEY=$STRIPE_SECRET_KEY
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret_here

# Stripe Product IDs (will be created automatically)
STRIPE_STARTER_MONTHLY_PRICE=price_starter_monthly
STRIPE_STARTER_ANNUAL_PRICE=price_starter_annual
STRIPE_PROFESSIONAL_MONTHLY_PRICE=price_professional_monthly
STRIPE_PROFESSIONAL_ANNUAL_PRICE=price_professional_annual
STRIPE_ENTERPRISE_MONTHLY_PRICE=price_enterprise_monthly
STRIPE_ENTERPRISE_ANNUAL_PRICE=price_enterprise_annual

# Beta Pricing (50% off)
BETA_DISCOUNT_PERCENT=50
BETA_COUPON_CODE=BETA50
EOF
    
    print_status ".env file updated with Stripe configuration"
}

# Create Stripe products and prices
create_stripe_products() {
    print_info "Creating Stripe products and prices..."
    
    # Check if stripe CLI is installed
    if ! command -v stripe &> /dev/null; then
        print_warning "Stripe CLI not found. Installing..."
        
        # Install Stripe CLI based on OS
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            curl -s https://packages.stripe.dev/api/security/keypair/stripe-cli-gpg/public | gpg --dearmor | sudo tee /usr/share/keyrings/stripe.gpg
            echo "deb [signed-by=/usr/share/keyrings/stripe.gpg] https://packages.stripe.dev/stripe-cli-debian-local stable main" | sudo tee -a /etc/apt/sources.list.d/stripe.list
            sudo apt update
            sudo apt install stripe
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            brew install stripe/stripe-cli/stripe
        else
            print_error "Please install Stripe CLI manually from https://stripe.com/docs/stripe-cli"
            exit 1
        fi
    fi
    
    # Login to Stripe
    print_info "Please login to Stripe CLI..."
    stripe login
    
    # Create products using Stripe CLI
    print_info "Creating Stripe products..."
    
    # Starter Plan
    STARTER_PRODUCT=$(stripe products create \
        --name="AI Crypto Trading - Starter" \
        --description="Perfect for beginners getting started with AI trading" \
        --type=service \
        --format=json | jq -r '.id')
    
    STARTER_MONTHLY_PRICE=$(stripe prices create \
        --product=$STARTER_PRODUCT \
        --unit-amount=4900 \
        --currency=usd \
        --recurring='{"interval":"month"}' \
        --nickname="Starter Monthly" \
        --format=json | jq -r '.id')
    
    STARTER_ANNUAL_PRICE=$(stripe prices create \
        --product=$STARTER_PRODUCT \
        --unit-amount=49000 \
        --currency=usd \
        --recurring='{"interval":"year"}' \
        --nickname="Starter Annual" \
        --format=json | jq -r '.id')
    
    # Professional Plan
    PROFESSIONAL_PRODUCT=$(stripe products create \
        --name="AI Crypto Trading - Professional" \
        --description="Advanced features for serious traders" \
        --type=service \
        --format=json | jq -r '.id')
    
    PROFESSIONAL_MONTHLY_PRICE=$(stripe prices create \
        --product=$PROFESSIONAL_PRODUCT \
        --unit-amount=19900 \
        --currency=usd \
        --recurring='{"interval":"month"}' \
        --nickname="Professional Monthly" \
        --format=json | jq -r '.id')
    
    PROFESSIONAL_ANNUAL_PRICE=$(stripe prices create \
        --product=$PROFESSIONAL_PRODUCT \
        --unit-amount=199000 \
        --currency=usd \
        --recurring='{"interval":"year"}' \
        --nickname="Professional Annual" \
        --format=json | jq -r '.id')
    
    # Enterprise Plan
    ENTERPRISE_PRODUCT=$(stripe products create \
        --name="AI Crypto Trading - Enterprise" \
        --description="Full platform access with enterprise features" \
        --type=service \
        --format=json | jq -r '.id')
    
    ENTERPRISE_MONTHLY_PRICE=$(stripe prices create \
        --product=$ENTERPRISE_PRODUCT \
        --unit-amount=99900 \
        --currency=usd \
        --recurring='{"interval":"month"}' \
        --nickname="Enterprise Monthly" \
        --format=json | jq -r '.id')
    
    ENTERPRISE_ANNUAL_PRICE=$(stripe prices create \
        --product=$ENTERPRISE_PRODUCT \
        --unit-amount=999000 \
        --currency=usd \
        --recurring='{"interval":"year"}' \
        --nickname="Enterprise Annual" \
        --format=json | jq -r '.id')
    
    # Update .env with actual price IDs
    sed -i "s/price_starter_monthly/$STARTER_MONTHLY_PRICE/g" .env
    sed -i "s/price_starter_annual/$STARTER_ANNUAL_PRICE/g" .env
    sed -i "s/price_professional_monthly/$PROFESSIONAL_MONTHLY_PRICE/g" .env
    sed -i "s/price_professional_annual/$PROFESSIONAL_ANNUAL_PRICE/g" .env
    sed -i "s/price_enterprise_monthly/$ENTERPRISE_MONTHLY_PRICE/g" .env
    sed -i "s/price_enterprise_annual/$ENTERPRISE_ANNUAL_PRICE/g" .env
    
    print_status "Stripe products and prices created successfully"
}

# Create beta discount coupon
create_beta_coupon() {
    print_info "Creating beta discount coupon..."
    
    stripe coupons create \
        --id=BETA50 \
        --percent-off=50 \
        --duration=repeating \
        --duration-in-months=6 \
        --name="Beta Program 50% Off" \
        --max-redemptions=500
    
    print_status "Beta discount coupon created (BETA50 - 50% off for 6 months)"
}

# Setup webhook endpoint
setup_webhook() {
    print_info "Setting up webhook endpoint..."
    
    echo ""
    echo "To complete the webhook setup:"
    echo "1. Go to your Stripe Dashboard > Developers > Webhooks"
    echo "2. Click 'Add endpoint'"
    echo "3. Enter your endpoint URL: https://yourdomain.com/webhooks/stripe"
    echo "4. Select these events:"
    echo "   - customer.subscription.created"
    echo "   - customer.subscription.updated"
    echo "   - customer.subscription.deleted"
    echo "   - invoice.payment_succeeded"
    echo "   - invoice.payment_failed"
    echo "5. Copy the webhook signing secret and update your .env file"
    echo ""
    
    read -p "Press Enter when you've completed the webhook setup..."
    print_status "Webhook setup instructions provided"
}

# Test Stripe integration
test_stripe_integration() {
    print_info "Testing Stripe integration..."
    
    # Test API connection
    stripe customers list --limit=1 > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        print_status "Stripe API connection successful"
    else
        print_error "Stripe API connection failed"
        exit 1
    fi
    
    # Test products exist
    stripe products list --limit=10 | grep -q "AI Crypto Trading"
    if [ $? -eq 0 ]; then
        print_status "Stripe products found"
    else
        print_warning "Stripe products not found - may need manual creation"
    fi
}

# Generate summary
generate_summary() {
    echo ""
    echo "🎉 Stripe Integration Setup Complete!"
    echo "====================================="
    echo ""
    echo "✅ Stripe account configured"
    echo "✅ Products and prices created"
    echo "✅ Beta discount coupon ready"
    echo "✅ Environment variables set"
    echo ""
    echo "📋 Next Steps:"
    echo "1. Complete webhook setup in Stripe dashboard"
    echo "2. Update STRIPE_WEBHOOK_SECRET in .env file"
    echo "3. Deploy your application"
    echo "4. Test subscription flow"
    echo ""
    echo "💰 You're ready to start making money!"
    echo ""
    echo "🔗 Useful Links:"
    echo "- Stripe Dashboard: https://dashboard.stripe.com"
    echo "- Webhook Setup: https://dashboard.stripe.com/webhooks"
    echo "- API Keys: https://dashboard.stripe.com/apikeys"
    echo ""
}

# Main execution
main() {
    echo ""
    print_info "Starting Stripe integration setup..."
    echo ""
    
    check_env_file
    setup_stripe_account
    update_env_file
    
    # Ask if user wants to create products automatically
    read -p "Do you want to create Stripe products automatically? (y/n): " create_products
    if [[ $create_products =~ ^[Yy]$ ]]; then
        create_stripe_products
        create_beta_coupon
        test_stripe_integration
    else
        print_info "Skipping automatic product creation"
        print_info "You can create products manually in your Stripe dashboard"
    fi
    
    setup_webhook
    generate_summary
}

# Run the script
main "$@"
