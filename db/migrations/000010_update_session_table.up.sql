-- Add a unique constraint to the token column
ALTER TABLE "sessions" ADD CONSTRAINT "unique_token" UNIQUE ("token");

-- Create an index on the token column for efficient searching
CREATE INDEX "idx_token" ON "sessions" ("token");

-- Create a new function with the desired return type (in this case, trigger).
CREATE OR REPLACE FUNCTION public.delete_expired_sessions_trigger()
RETURNS trigger
LANGUAGE plpgsql
AS $function$
BEGIN
  DELETE FROM sessions WHERE expires_at < NOW();
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
