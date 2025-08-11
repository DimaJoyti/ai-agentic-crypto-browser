-- Migration: 005_solana_integration.sql
-- Description: Add Solana blockchain integration tables
-- Created: 2024-01-09
-- Database: PostgreSQL
-- Note: This file contains valid PostgreSQL syntax. IDE errors are due to SQL Server parser being used instead of PostgreSQL parser.

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Solana wallets table
CREATE TABLE IF NOT EXISTS solana_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    public_key VARCHAR(44) UNIQUE NOT NULL,
    wallet_type VARCHAR(50) NOT NULL, -- phantom, solflare, backpack, glow, ledger, trezor
    is_active BOOLEAN DEFAULT true,
    balance DECIMAL(36,18) DEFAULT 0,
    last_balance_update TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for solana_wallets
CREATE INDEX IF NOT EXISTS idx_solana_wallets_user_id ON solana_wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_solana_wallets_public_key ON solana_wallets(public_key);
CREATE INDEX IF NOT EXISTS idx_solana_wallets_active ON solana_wallets(is_active) WHERE is_active = true;

-- Solana transactions table
CREATE TABLE IF NOT EXISTS solana_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID REFERENCES solana_wallets(id) ON DELETE CASCADE,
    signature VARCHAR(88) UNIQUE NOT NULL,
    from_address VARCHAR(44) NOT NULL,
    to_address VARCHAR(44) NOT NULL,
    amount DECIMAL(36,18) NOT NULL,
    token_mint VARCHAR(44), -- NULL for SOL transfers
    transaction_type VARCHAR(50) NOT NULL, -- transfer, token_transfer, swap, stake, unstake, nft_transfer
    status VARCHAR(20) DEFAULT 'pending', -- pending, confirmed, finalized, failed
    priority VARCHAR(20) DEFAULT 'medium', -- low, medium, high, max
    max_fee DECIMAL(36,18),
    actual_fee DECIMAL(36,18),
    memo TEXT,
    block_time TIMESTAMP,
    slot BIGINT,
    error_message TEXT,
    logs JSONB,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for solana_transactions
CREATE INDEX IF NOT EXISTS idx_solana_transactions_wallet_id ON solana_transactions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_solana_transactions_signature ON solana_transactions(signature);
CREATE INDEX IF NOT EXISTS idx_solana_transactions_from_address ON solana_transactions(from_address);
CREATE INDEX IF NOT EXISTS idx_solana_transactions_to_address ON solana_transactions(to_address);
CREATE INDEX IF NOT EXISTS idx_solana_transactions_status ON solana_transactions(status);
CREATE INDEX IF NOT EXISTS idx_solana_transactions_type ON solana_transactions(transaction_type);
CREATE INDEX IF NOT EXISTS idx_solana_transactions_created_at ON solana_transactions(created_at DESC);

-- Solana DeFi positions table
CREATE TABLE IF NOT EXISTS solana_defi_positions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL REFERENCES solana_wallets(id) ON DELETE CASCADE,
    protocol VARCHAR(50) NOT NULL, -- jupiter, raydium, orca, marinade, etc.
    position_type VARCHAR(50) NOT NULL, -- liquidity, stake, lend, borrow, farm
    pool_address VARCHAR(44),
    token_a VARCHAR(44),
    token_b VARCHAR(44),
    amount_a DECIMAL(36,18),
    amount_b DECIMAL(36,18),
    lp_tokens DECIMAL(36,18),
    apy DECIMAL(10,4),
    entry_price DECIMAL(36,18),
    current_value DECIMAL(36,18),
    unrealized_pnl DECIMAL(36,18),
    fees_earned DECIMAL(36,18),
    rewards_earned DECIMAL(36,18),
    is_active BOOLEAN DEFAULT true,
    opened_at TIMESTAMP DEFAULT NOW(),
    closed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for solana_defi_positions
CREATE INDEX IF NOT EXISTS idx_solana_defi_positions_wallet_id ON solana_defi_positions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_solana_defi_positions_protocol ON solana_defi_positions(protocol);
CREATE INDEX IF NOT EXISTS idx_solana_defi_positions_type ON solana_defi_positions(position_type);
CREATE INDEX IF NOT EXISTS idx_solana_defi_positions_active ON solana_defi_positions(is_active) WHERE is_active = true;

-- Solana NFT holdings table
CREATE TABLE IF NOT EXISTS solana_nft_holdings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL REFERENCES solana_wallets(id) ON DELETE CASCADE,
    mint_address VARCHAR(44) UNIQUE NOT NULL,
    collection_address VARCHAR(44),
    name VARCHAR(255),
    symbol VARCHAR(10),
    description TEXT,
    image_url TEXT,
    metadata_uri TEXT,
    attributes JSONB,
    creators JSONB,
    rarity_rank INTEGER,
    rarity_score DECIMAL(10,4),
    floor_price DECIMAL(36,18),
    last_sale_price DECIMAL(36,18),
    estimated_value DECIMAL(36,18),
    marketplace VARCHAR(50), -- magic_eden, tensor, opensea, etc.
    is_listed BOOLEAN DEFAULT false,
    listing_price DECIMAL(36,18),
    listing_marketplace VARCHAR(50),
    acquired_at TIMESTAMP DEFAULT NOW(),
    acquired_price DECIMAL(36,18),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for solana_nft_holdings
CREATE INDEX IF NOT EXISTS idx_solana_nft_holdings_wallet_id ON solana_nft_holdings(wallet_id);
CREATE INDEX IF NOT EXISTS idx_solana_nft_holdings_mint ON solana_nft_holdings(mint_address);
CREATE INDEX IF NOT EXISTS idx_solana_nft_holdings_collection ON solana_nft_holdings(collection_address);
CREATE INDEX IF NOT EXISTS idx_solana_nft_holdings_listed ON solana_nft_holdings(is_listed) WHERE is_listed = true;

-- Solana token balances table (for tracking all token holdings)
CREATE TABLE IF NOT EXISTS solana_token_balances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL REFERENCES solana_wallets(id) ON DELETE CASCADE,
    mint_address VARCHAR(44) NOT NULL,
    token_account VARCHAR(44) NOT NULL,
    symbol VARCHAR(20),
    name VARCHAR(100),
    decimals INTEGER DEFAULT 9,
    balance DECIMAL(36,18) NOT NULL DEFAULT 0,
    usd_value DECIMAL(36,18),
    price_per_token DECIMAL(36,18),
    last_updated TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(wallet_id, mint_address)
);

-- Create indexes for solana_token_balances
CREATE INDEX IF NOT EXISTS idx_solana_token_balances_wallet_id ON solana_token_balances(wallet_id);
CREATE INDEX IF NOT EXISTS idx_solana_token_balances_mint ON solana_token_balances(mint_address);
CREATE INDEX IF NOT EXISTS idx_solana_token_balances_symbol ON solana_token_balances(symbol);

-- Solana staking positions table
CREATE TABLE IF NOT EXISTS solana_staking_positions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL REFERENCES solana_wallets(id) ON DELETE CASCADE,
    stake_account VARCHAR(44) UNIQUE NOT NULL,
    validator_address VARCHAR(44) NOT NULL,
    validator_name VARCHAR(100),
    staked_amount DECIMAL(36,18) NOT NULL,
    rewards_earned DECIMAL(36,18) DEFAULT 0,
    apy DECIMAL(10,4),
    status VARCHAR(20) DEFAULT 'active', -- active, deactivating, inactive
    activation_epoch BIGINT,
    deactivation_epoch BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for solana_staking_positions
CREATE INDEX IF NOT EXISTS idx_solana_staking_positions_wallet_id ON solana_staking_positions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_solana_staking_positions_validator ON solana_staking_positions(validator_address);
CREATE INDEX IF NOT EXISTS idx_solana_staking_positions_status ON solana_staking_positions(status);

-- Solana program interactions table (for tracking smart contract interactions)
CREATE TABLE IF NOT EXISTS solana_program_interactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL REFERENCES solana_wallets(id) ON DELETE CASCADE,
    program_id VARCHAR(44) NOT NULL,
    program_name VARCHAR(100),
    instruction_name VARCHAR(100),
    transaction_signature VARCHAR(88) NOT NULL,
    accounts JSONB,
    instruction_data BYTEA,
    compute_units_used BIGINT,
    success BOOLEAN DEFAULT true,
    error_message TEXT,
    logs JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for solana_program_interactions
CREATE INDEX IF NOT EXISTS idx_solana_program_interactions_wallet_id ON solana_program_interactions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_solana_program_interactions_program_id ON solana_program_interactions(program_id);
CREATE INDEX IF NOT EXISTS idx_solana_program_interactions_signature ON solana_program_interactions(transaction_signature);

-- Create updated_at trigger function if it doesn't exist
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Add updated_at triggers for all tables
CREATE TRIGGER update_solana_wallets_updated_at BEFORE UPDATE ON solana_wallets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_solana_transactions_updated_at BEFORE UPDATE ON solana_transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_solana_defi_positions_updated_at BEFORE UPDATE ON solana_defi_positions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_solana_nft_holdings_updated_at BEFORE UPDATE ON solana_nft_holdings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_solana_token_balances_updated_at BEFORE UPDATE ON solana_token_balances FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_solana_staking_positions_updated_at BEFORE UPDATE ON solana_staking_positions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert some sample data for testing (optional)
-- This would be removed in production
INSERT INTO solana_wallets (user_id, public_key, wallet_type) 
SELECT id, 'DemoSolanaWallet1234567890123456789012345', 'phantom' 
FROM users LIMIT 1
ON CONFLICT (public_key) DO NOTHING;

-- Add comments for documentation
COMMENT ON TABLE solana_wallets IS 'Stores Solana wallet connections for users';
COMMENT ON TABLE solana_transactions IS 'Tracks all Solana transactions (SOL and token transfers)';
COMMENT ON TABLE solana_defi_positions IS 'Tracks DeFi positions across Solana protocols';
COMMENT ON TABLE solana_nft_holdings IS 'Stores NFT holdings and metadata';
COMMENT ON TABLE solana_token_balances IS 'Tracks token balances for all wallets';
COMMENT ON TABLE solana_staking_positions IS 'Tracks native Solana staking positions';
COMMENT ON TABLE solana_program_interactions IS 'Logs all smart contract interactions';

-- Migration completed successfully
SELECT 'Solana integration tables created successfully' as result;
