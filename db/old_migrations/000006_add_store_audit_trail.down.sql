-- Drop the foreign keys from store_audit_trail
ALTER TABLE "store_audit_trail" DROP CONSTRAINT IF EXISTS store_audit_trail_store_id_fkey;
ALTER TABLE "store_audit_trail" DROP CONSTRAINT IF EXISTS store_audit_trail_user_id_fkey;

-- Drop the store_audit_trail table
DROP TABLE IF EXISTS "store_audit_trail";

-- Revert changes to store_owners table
ALTER TABLE store_owners DROP COLUMN IF EXISTS is_primary;
ALTER TABLE store_owners DROP COLUMN IF EXISTS access_levels;
ALTER TABLE store_owners ADD COLUMN IF NOT EXISTS access_level smallint;
