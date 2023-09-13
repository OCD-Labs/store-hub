-- Drop the unique constraint on the token column
ALTER TABLE "sessions" DROP CONSTRAINT "unique_token";

-- Drop the index on the token column
DROP INDEX "idx_token";

-- Drop the trigger if it exists
DROP TRIGGER IF EXISTS trigger_delete_expired_sessions ON sessions;

-- Drop the function
DROP FUNCTION public.delete_expired_sessions_trigger();

-- Recreate the existing function
CREATE OR REPLACE FUNCTION delete_expired_sessions()
RETURNS void AS $$
BEGIN
  DELETE FROM sessions WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;
