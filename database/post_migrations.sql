-- Additional indexes for better query performance
-- Run these after GORM auto-migration

-- Enable trigram extension for better text search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Composite index for user_id + status (common filter combination)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tasks_user_status ON tasks(user_id, status);

-- Index for title searches (supports ILIKE queries better)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tasks_title_gin ON tasks USING gin(title gin_trgm_ops);

-- Index for description searches
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tasks_description_gin ON tasks USING gin(description gin_trgm_ops);

-- Composite index for user_id + created_at (for sorting)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tasks_user_created ON tasks(user_id, created_at DESC);

-- Composite index for user_id + deadline (for deadline sorting)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tasks_user_deadline ON tasks(user_id, deadline DESC NULLS LAST);

-- Partial index for active tasks only (excludes soft-deleted)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tasks_active ON tasks(user_id, created_at) WHERE deleted_at IS NULL;
