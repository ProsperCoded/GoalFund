-- Migration: Add settlement account fields and password tracking
-- Description: Adds bank account details for refunds/settlements and password setup tracking

-- Add settlement account columns to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS has_set_password BOOLEAN DEFAULT TRUE,
ADD COLUMN IF NOT EXISTS settlement_bank_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS settlement_account_number VARCHAR(20),
ADD COLUMN IF NOT EXISTS settlement_account_name VARCHAR(255);

-- Add indexes for settlement account lookups
CREATE INDEX IF NOT EXISTS idx_users_settlement_account_number ON users(settlement_account_number);
CREATE INDEX IF NOT EXISTS idx_users_has_set_password ON users(has_set_password);

-- Add comments to document the new columns
COMMENT ON COLUMN users.has_set_password IS 'Indicates whether the user has set their password (false for email-only contributions)';
COMMENT ON COLUMN users.settlement_bank_name IS 'Bank name for user settlements and refunds';
COMMENT ON COLUMN users.settlement_account_number IS 'Account number for user settlements and refunds';
COMMENT ON COLUMN users.settlement_account_name IS 'Account name for user settlements and refunds';
