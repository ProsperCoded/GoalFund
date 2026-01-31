-- Migration: Add KYC verification fields to users table
-- Description: Adds NIN (National Identification Number) and KYC verification tracking fields

-- Add KYC verification columns to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS nin VARCHAR(11),
ADD COLUMN IF NOT EXISTS kyc_verified BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS kyc_verified_at TIMESTAMP;

-- Add index on NIN for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_nin ON users(nin);

-- Add index on kyc_verified_at for reporting
CREATE INDEX IF NOT EXISTS idx_users_kyc_verified_at ON users(kyc_verified_at);

-- Add comment to document the NIN column
COMMENT ON COLUMN users.nin IS 'National Identification Number (11 digits) - used for basic KYC verification';
COMMENT ON COLUMN users.kyc_verified IS 'Indicates whether the user has completed KYC verification';
COMMENT ON COLUMN users.kyc_verified_at IS 'Timestamp when KYC verification was completed';
