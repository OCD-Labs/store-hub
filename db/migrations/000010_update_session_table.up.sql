-- Create a new function with the desired return type (in this case, trigger).
CREATE OR REPLACE FUNCTION public.delete_expired_sessions_trigger()
RETURNS trigger
LANGUAGE plpgsql
AS $function$
BEGIN
  -- Check if there are expired sessions before attempting to delete.
  IF EXISTS (SELECT 1 FROM sessions WHERE expires_at < NOW()) THEN
    DELETE FROM sessions WHERE expires_at < NOW();
  END IF;
  RETURN NULL;
END;
$function$;

-- Create a trigger that fires before INSERT or UPDATE on the sessions table
CREATE TRIGGER trigger_delete_expired_sessions
BEFORE INSERT OR UPDATE
ON sessions
FOR EACH ROW
EXECUTE FUNCTION delete_expired_sessions_trigger();

DROP FUNCTION public.delete_expired_sessions();
