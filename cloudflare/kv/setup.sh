#!/bin/bash

# Cloudflare KV Namespace Setup Script
# Creates KV namespaces for caching and session management

set -e

echo "🗄️ Setting up Cloudflare KV Namespaces for AI Agentic Crypto Browser"

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

echo "📦 Creating KV namespaces..."

# Create Cache namespace
echo "Creating CACHE namespace..."
CACHE_NAMESPACE=$(wrangler kv:namespace create "CACHE" --preview false)
CACHE_PREVIEW_NAMESPACE=$(wrangler kv:namespace create "CACHE" --preview true)

echo "Creating SESSIONS namespace..."
SESSIONS_NAMESPACE=$(wrangler kv:namespace create "SESSIONS" --preview false)
SESSIONS_PREVIEW_NAMESPACE=$(wrangler kv:namespace create "SESSIONS" --preview true)

echo "Creating RATE_LIMIT namespace..."
RATE_LIMIT_NAMESPACE=$(wrangler kv:namespace create "RATE_LIMIT" --preview false)
RATE_LIMIT_PREVIEW_NAMESPACE=$(wrangler kv:namespace create "RATE_LIMIT" --preview true)

echo "Creating USER_DATA namespace..."
USER_DATA_NAMESPACE=$(wrangler kv:namespace create "USER_DATA" --preview false)
USER_DATA_PREVIEW_NAMESPACE=$(wrangler kv:namespace create "USER_DATA" --preview true)

echo "✅ All KV namespaces created successfully!"
echo ""
echo "📝 Please update your wrangler.toml file with the following KV namespace bindings:"
echo ""
echo "[[kv_namespaces]]"
echo "binding = \"CACHE\""
echo "id = \"<CACHE_NAMESPACE_ID>\""
echo "preview_id = \"<CACHE_PREVIEW_NAMESPACE_ID>\""
echo ""
echo "[[kv_namespaces]]"
echo "binding = \"SESSIONS\""
echo "id = \"<SESSIONS_NAMESPACE_ID>\""
echo "preview_id = \"<SESSIONS_PREVIEW_NAMESPACE_ID>\""
echo ""
echo "[[kv_namespaces]]"
echo "binding = \"RATE_LIMIT\""
echo "id = \"<RATE_LIMIT_NAMESPACE_ID>\""
echo "preview_id = \"<RATE_LIMIT_PREVIEW_NAMESPACE_ID>\""
echo ""
echo "[[kv_namespaces]]"
echo "binding = \"USER_DATA\""
echo "id = \"<USER_DATA_NAMESPACE_ID>\""
echo "preview_id = \"<USER_DATA_PREVIEW_NAMESPACE_ID>\""
echo ""
echo "🔍 To get the namespace IDs, run:"
echo "wrangler kv:namespace list"
echo ""
echo "📋 KV Namespace Usage:"
echo "- CACHE: General application caching (API responses, computed data)"
echo "- SESSIONS: User session management and authentication"
echo "- RATE_LIMIT: Rate limiting counters and tracking"
echo "- USER_DATA: User preferences and temporary data"
echo ""
echo "🎉 KV setup completed!"
