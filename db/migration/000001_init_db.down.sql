-- Start by dropping the foreign key constraints
ALTER TABLE "item_ratings" DROP CONSTRAINT IF EXISTS item_ratings_user_id_fkey;
ALTER TABLE "item_ratings" DROP CONSTRAINT IF EXISTS item_ratings_item_id_fkey;
ALTER TABLE "items" DROP CONSTRAINT IF EXISTS items_store_id_fkey;
ALTER TABLE "store_owners" DROP CONSTRAINT IF EXISTS store_owners_user_id_fkey;
ALTER TABLE "store_owners" DROP CONSTRAINT IF EXISTS store_owners_store_id_fkey;
ALTER TABLE "sessions" DROP CONSTRAINT IF EXISTS sessions_user_id_fkey;

-- Drop the tables
DROP TABLE IF EXISTS "item_ratings";
DROP TABLE IF EXISTS "items";
DROP TABLE IF EXISTS "store_owners";
DROP TABLE IF EXISTS "stores";
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "users";