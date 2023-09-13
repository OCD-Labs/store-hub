-- Drop the unique constraint on the token column
ALTER TABLE "sessions" DROP CONSTRAINT "unique_token";

-- Drop the index on the token column
DROP INDEX "idx_token";

-- Drop the trigger if it exists
DROP TRIGGER IF EXISTS trigger_delete_expired_sessions ON sessions;
