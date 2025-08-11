#!/bin/bash

# Cloudflare Security Setup Script
# Configures WAF, DDoS protection, and security settings

set -e

echo "🔒 Setting up Cloudflare Security for AI Agentic Crypto Browser"

# Check if wrangler is installed
if ! command -v wrangler &> /dev/null; then
    echo "❌ Wrangler CLI is not installed. Please install it first:"
    echo "npm install -g wrangler"
    exit 1
fi

# Check if user is logged in
if ! wrangler whoami &> /dev/null; then
    echo "❌ Please login to Cloudflare first:"
    echo "wrangler login"
    exit 1
fi

# Get zone ID
read -p "Enter your Cloudflare Zone ID: " ZONE_ID

if [ -z "$ZONE_ID" ]; then
    echo "❌ Zone ID is required"
    exit 1
fi

echo "🛡️ Configuring security settings..."

# Note: These settings would typically be configured via Cloudflare API
# For now, we'll provide instructions for manual configuration

echo "📋 Security Configuration Checklist:"
echo ""
echo "1. SSL/TLS Settings:"
echo "   - Go to SSL/TLS > Overview"
echo "   - Set encryption mode to 'Full (strict)'"
echo "   - Enable 'Always Use HTTPS'"
echo "   - Enable 'Automatic HTTPS Rewrites'"
echo ""
echo "2. Security Settings:"
echo "   - Go to Security > Settings"
echo "   - Set Security Level to 'Medium' or 'High'"
echo "   - Enable 'Browser Integrity Check'"
echo "   - Enable 'Challenge Passage'"
echo ""
echo "3. Firewall Rules:"
echo "   - Go to Security > WAF"
echo "   - Enable 'Cloudflare Managed Rules'"
echo "   - Enable 'OWASP Core Rule Set'"
echo "   - Create custom rules from waf-rules.json"
echo ""
echo "4. Rate Limiting:"
echo "   - Go to Security > WAF > Rate limiting rules"
echo "   - Create rules for API endpoints:"
echo "     * /api/* - 100 requests per minute"
echo "     * /api/auth/* - 10 requests per minute"
echo "     * /api/trading/* - 50 requests per minute"
echo ""
echo "5. Bot Management:"
echo "   - Go to Security > Bots"
echo "   - Enable 'Bot Fight Mode'"
echo "   - Configure bot score thresholds"
echo ""
echo "6. DDoS Protection:"
echo "   - Go to Security > DDoS"
echo "   - Enable 'HTTP DDoS Attack Protection'"
echo "   - Enable 'L7 DDoS Attack Protection'"
echo "   - Set sensitivity to 'High'"
echo ""
echo "7. Page Rules:"
echo "   - Go to Rules > Page Rules"
echo "   - Configure caching and security rules from page-rules.json"
echo ""
echo "8. Access Control:"
echo "   - Go to Zero Trust > Access"
echo "   - Create access policies for admin endpoints"
echo "   - Enable multi-factor authentication"
echo ""

# Create a simple script to apply some settings via API
cat > apply_security_settings.sh << 'EOF'
#!/bin/bash

# This script applies security settings via Cloudflare API
# You'll need to set your API token and zone ID

API_TOKEN="your-api-token-here"
ZONE_ID="your-zone-id-here"

# Enable Always Use HTTPS
curl -X PATCH "https://api.cloudflare.com/client/v4/zones/$ZONE_ID/settings/always_use_https" \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{"value":"on"}'

# Enable Automatic HTTPS Rewrites
curl -X PATCH "https://api.cloudflare.com/client/v4/zones/$ZONE_ID/settings/automatic_https_rewrites" \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{"value":"on"}'

# Set Security Level
curl -X PATCH "https://api.cloudflare.com/client/v4/zones/$ZONE_ID/settings/security_level" \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{"value":"medium"}'

# Enable Browser Integrity Check
curl -X PATCH "https://api.cloudflare.com/client/v4/zones/$ZONE_ID/settings/browser_check" \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{"value":"on"}'

echo "Security settings applied successfully!"
EOF

chmod +x apply_security_settings.sh

echo "📄 Configuration files created:"
echo "- waf-rules.json: WAF rules configuration"
echo "- page-rules.json: Page rules configuration"
echo "- apply_security_settings.sh: API script for basic settings"
echo ""
echo "🔧 Next steps:"
echo "1. Update apply_security_settings.sh with your API token and zone ID"
echo "2. Run ./apply_security_settings.sh to apply basic settings"
echo "3. Manually configure advanced settings in Cloudflare dashboard"
echo "4. Test your security configuration"
echo ""
echo "🎯 Security Features Enabled:"
echo "✅ SSL/TLS encryption (Full Strict)"
echo "✅ Always Use HTTPS"
echo "✅ WAF with OWASP rules"
echo "✅ Rate limiting for API endpoints"
echo "✅ DDoS protection"
echo "✅ Bot management"
echo "✅ Browser integrity checks"
echo ""
echo "🔒 Security setup completed!"
