-- Migration 004: Add currency column to users
ALTER TABLE users ADD COLUMN IF NOT EXISTS currency VARCHAR(10) DEFAULT 'EUR';
