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
