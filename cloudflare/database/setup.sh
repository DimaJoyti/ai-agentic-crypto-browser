#!/bin/bash

# Cloudflare D1 Database Setup Script
# This script creates and configures the D1 database for AI Agentic Crypto Browser

set -e

echo "🚀 Setting up Cloudflare D1 Database for AI Agentic Crypto Browser"

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

# Database name
DB_NAME="ai-crypto-browser-db"

echo "📊 Creating D1 database: $DB_NAME"

# Create the database
wrangler d1 create $DB_NAME

echo "✅ Database created successfully!"
echo ""
echo "📝 Please update your wrangler.toml file with the database ID shown above."
echo "Add this to your [[d1_databases]] section:"
echo ""
echo "[[d1_databases]]"
echo "binding = \"DB\""
echo "database_name = \"$DB_NAME\""
echo "database_id = \"<DATABASE_ID_FROM_OUTPUT_ABOVE>\""
echo ""

read -p "Press Enter after updating wrangler.toml to continue with schema setup..."

echo "🔧 Running database migrations..."

# Run migrations in order
echo "Running migration 001: Initial schema..."
wrangler d1 execute $DB_NAME --file=./migrations/001_initial_schema.sql

echo "Running migration 002: Trading tables..."
wrangler d1 execute $DB_NAME --file=./migrations/002_trading_tables.sql

echo "Running migration 003: AI analytics tables..."
wrangler d1 execute $DB_NAME --file=./migrations/003_ai_analytics_tables.sql

echo "Running migration 004: User preferences..."
wrangler d1 execute $DB_NAME --file=./migrations/004_user_preferences.sql

echo "✅ All migrations completed successfully!"

# Ask if user wants to seed data
read -p "Do you want to seed the database with sample data? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🌱 Seeding database with sample data..."
    wrangler d1 execute $DB_NAME --file=./seeds/001_sample_data.sql
    echo "✅ Sample data seeded successfully!"
fi

echo ""
echo "🎉 Database setup completed!"
echo ""
echo "📋 Next steps:"
echo "1. Update your Worker's wrangler.toml with the database ID"
echo "2. Deploy your Worker: wrangler deploy"
echo "3. Test the API endpoints"
echo ""
echo "🔍 Useful commands:"
echo "- List databases: wrangler d1 list"
echo "- Query database: wrangler d1 execute $DB_NAME --command=\"SELECT * FROM users LIMIT 5\""
echo "- Backup database: wrangler d1 export $DB_NAME --output=backup.sql"
echo ""
