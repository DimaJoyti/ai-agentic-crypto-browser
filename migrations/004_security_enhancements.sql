-- Security Enhancements Migration
-- This migration adds comprehensive security features to the authentication system

-- Add security-related columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS mfa_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS mfa_secret TEXT DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS mfa_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS failed_login_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS locked_until TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_ip INET;
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_changed_at TIMESTAMP DEFAULT NOW();
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_email_verified BOOLEAN DEFAULT FALSE;

-- Create login attempts table for security monitoring
CREATE TABLE IF NOT EXISTS login_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN NOT NULL,
    reason TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create user sessions table for session management
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(255),
    ip_address INET,
    user_agent TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create MFA backup codes table
CREATE TABLE IF NOT EXISTS mfa_backup_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash VARCHAR(255) NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create security audit log table
CREATE TABLE IF NOT EXISTS security_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100),
    resource_id UUID,
    ip_address INET,
    user_agent TEXT,
    details JSONB,
    success BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create password history table to prevent password reuse
CREATE TABLE IF NOT EXISTS password_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create blacklisted tokens table
CREATE TABLE IF NOT EXISTS blacklisted_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    reason VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create security policies table
CREATE TABLE IF NOT EXISTS security_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    policy_type VARCHAR(50) NOT NULL,
    rules JSONB NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 100,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create user roles table
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_name VARCHAR(100) NOT NULL,
    team_id UUID,
    granted_by UUID REFERENCES users(id),
    granted_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    UNIQUE(user_id, role_name, team_id)
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    conditions JSONB,
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create role permissions table
CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name VARCHAR(100) NOT NULL,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(role_name, permission_id)
);

-- Create WebAuthn credentials table
CREATE TABLE IF NOT EXISTS webauthn_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credential_id BYTEA NOT NULL UNIQUE,
    public_key BYTEA NOT NULL,
    attestation_type VARCHAR(50),
    transport TEXT[],
    flags INTEGER,
    sign_count INTEGER DEFAULT 0,
    aaguid BYTEA,
    clone_warning BOOLEAN DEFAULT FALSE,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    last_used_at TIMESTAMP
);

-- Create rate limiting table
CREATE TABLE IF NOT EXISTS rate_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier VARCHAR(255) NOT NULL, -- IP address or user ID
    identifier_type VARCHAR(20) NOT NULL, -- 'ip' or 'user'
    endpoint VARCHAR(255),
    requests_count INTEGER DEFAULT 1,
    window_start TIMESTAMP DEFAULT NOW(),
    window_duration INTERVAL DEFAULT '1 hour',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(identifier, identifier_type, endpoint, window_start)
);

-- Create security settings table
CREATE TABLE IF NOT EXISTS security_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    setting_key VARCHAR(255) NOT NULL UNIQUE,
    setting_value JSONB NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_login_attempts_email ON login_attempts(email);
CREATE INDEX IF NOT EXISTS idx_login_attempts_ip ON login_attempts(ip_address);
CREATE INDEX IF NOT EXISTS idx_login_attempts_created_at ON login_attempts(created_at);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_is_active ON user_sessions(is_active);

CREATE INDEX IF NOT EXISTS idx_mfa_backup_codes_user_id ON mfa_backup_codes(user_id);
CREATE INDEX IF NOT EXISTS idx_mfa_backup_codes_used_at ON mfa_backup_codes(used_at);

CREATE INDEX IF NOT EXISTS idx_security_audit_log_user_id ON security_audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_security_audit_log_action ON security_audit_log(action);
CREATE INDEX IF NOT EXISTS idx_security_audit_log_created_at ON security_audit_log(created_at);

CREATE INDEX IF NOT EXISTS idx_password_history_user_id ON password_history(user_id);
CREATE INDEX IF NOT EXISTS idx_password_history_created_at ON password_history(created_at);

CREATE INDEX IF NOT EXISTS idx_blacklisted_tokens_token_hash ON blacklisted_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_blacklisted_tokens_expires_at ON blacklisted_tokens(expires_at);

CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_name ON user_roles(role_name);
CREATE INDEX IF NOT EXISTS idx_user_roles_team_id ON user_roles(team_id);

CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action);

CREATE INDEX IF NOT EXISTS idx_webauthn_credentials_user_id ON webauthn_credentials(user_id);
CREATE INDEX IF NOT EXISTS idx_webauthn_credentials_credential_id ON webauthn_credentials(credential_id);

CREATE INDEX IF NOT EXISTS idx_rate_limits_identifier ON rate_limits(identifier, identifier_type);
CREATE INDEX IF NOT EXISTS idx_rate_limits_window_start ON rate_limits(window_start);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_locked_until ON users(locked_until);
CREATE INDEX IF NOT EXISTS idx_users_last_login_at ON users(last_login_at);

-- Insert default security policies
INSERT INTO security_policies (name, description, policy_type, rules, priority) VALUES
('password_complexity', 'Default password complexity requirements', 'password', 
 '{"min_length": 8, "require_uppercase": true, "require_lowercase": true, "require_numbers": true, "require_symbols": false}', 100),
('session_timeout', 'Default session timeout policy', 'session',
 '{"timeout_minutes": 1440, "extend_on_activity": true, "max_concurrent_sessions": 5}', 100),
('mfa_requirement', 'MFA requirement policy', 'mfa',
 '{"require_for_admin": true, "require_for_sensitive_operations": true, "grace_period_hours": 24}', 100),
('rate_limiting', 'Default rate limiting policy', 'rate_limit',
 '{"requests_per_minute": 60, "requests_per_hour": 1000, "burst_limit": 10}', 100);

-- Insert default permissions
INSERT INTO permissions (name, description, resource, action, is_system) VALUES
('users.read', 'Read user information', 'users', 'read', true),
('users.write', 'Create and update users', 'users', 'write', true),
('users.delete', 'Delete users', 'users', 'delete', true),
('workflows.read', 'Read workflows', 'workflows', 'read', true),
('workflows.write', 'Create and update workflows', 'workflows', 'write', true),
('workflows.execute', 'Execute workflows', 'workflows', 'execute', true),
('workflows.delete', 'Delete workflows', 'workflows', 'delete', true),
('teams.read', 'Read team information', 'teams', 'read', true),
('teams.write', 'Create and update teams', 'teams', 'write', true),
('teams.manage', 'Manage team members and settings', 'teams', 'manage', true),
('analytics.read', 'Read analytics data', 'analytics', 'read', true),
('system.admin', 'Full system administration', 'system', 'admin', true);

-- Insert default role permissions
INSERT INTO role_permissions (role_name, permission_id) 
SELECT 'super_admin', id FROM permissions WHERE name = 'system.admin';

INSERT INTO role_permissions (role_name, permission_id) 
SELECT 'admin', id FROM permissions WHERE name IN ('users.read', 'users.write', 'workflows.read', 'workflows.write', 'workflows.execute', 'teams.read', 'teams.write', 'analytics.read');

INSERT INTO role_permissions (role_name, permission_id) 
SELECT 'team_owner', id FROM permissions WHERE name IN ('teams.manage', 'users.read', 'users.write', 'workflows.read', 'workflows.write', 'workflows.execute', 'analytics.read');

INSERT INTO role_permissions (role_name, permission_id) 
SELECT 'team_admin', id FROM permissions WHERE name IN ('users.read', 'workflows.read', 'workflows.write', 'workflows.execute', 'analytics.read');

INSERT INTO role_permissions (role_name, permission_id) 
SELECT 'user', id FROM permissions WHERE name IN ('workflows.read', 'workflows.execute');

INSERT INTO role_permissions (role_name, permission_id) 
SELECT 'viewer', id FROM permissions WHERE name = 'workflows.read';

-- Insert default security settings
INSERT INTO security_settings (setting_key, setting_value, description, is_system) VALUES
('password_policy', '{"min_length": 8, "max_length": 128, "require_uppercase": true, "require_lowercase": true, "require_numbers": true, "require_symbols": false, "forbidden_words": ["password", "123456", "qwerty"]}', 'Password complexity requirements', true),
('session_policy', '{"timeout_hours": 24, "extend_on_activity": true, "max_concurrent_sessions": 5, "require_mfa_for_sensitive": true}', 'Session management policy', true),
('rate_limiting', '{"global_requests_per_second": 100, "user_requests_per_second": 10, "ip_requests_per_second": 20, "burst_multiplier": 2}', 'Rate limiting configuration', true),
('security_headers', '{"hsts_max_age": 31536000, "csp_enabled": true, "frame_options": "DENY", "content_type_options": "nosniff"}', 'Security headers configuration', true);

-- Create function to clean up expired data
CREATE OR REPLACE FUNCTION cleanup_expired_security_data() RETURNS void AS $$
BEGIN
    -- Clean up expired sessions
    DELETE FROM user_sessions WHERE expires_at < NOW();
    
    -- Clean up expired blacklisted tokens
    DELETE FROM blacklisted_tokens WHERE expires_at < NOW();
    
    -- Clean up old login attempts (keep last 30 days)
    DELETE FROM login_attempts WHERE created_at < NOW() - INTERVAL '30 days';
    
    -- Clean up old audit logs (keep last 90 days)
    DELETE FROM security_audit_log WHERE created_at < NOW() - INTERVAL '90 days';
    
    -- Clean up old password history (keep last 12 passwords per user)
    DELETE FROM password_history 
    WHERE id NOT IN (
        SELECT id FROM (
            SELECT id, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) as rn
            FROM password_history
        ) ranked WHERE rn <= 12
    );
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically hash and store password history
CREATE OR REPLACE FUNCTION store_password_history() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' AND OLD.password_hash != NEW.password_hash THEN
        INSERT INTO password_history (user_id, password_hash)
        VALUES (NEW.id, OLD.password_hash);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_store_password_history
    AFTER UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION store_password_history();

-- Create function to check password reuse
CREATE OR REPLACE FUNCTION check_password_reuse(user_uuid UUID, new_password_hash TEXT, history_limit INTEGER DEFAULT 5) 
RETURNS BOOLEAN AS $$
DECLARE
    reused BOOLEAN := FALSE;
BEGIN
    SELECT EXISTS(
        SELECT 1 FROM password_history 
        WHERE user_id = user_uuid 
        AND password_hash = new_password_hash
        ORDER BY created_at DESC 
        LIMIT history_limit
    ) INTO reused;
    
    RETURN reused;
END;
$$ LANGUAGE plpgsql;
