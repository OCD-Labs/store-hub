-- DOWN Migration

-- Drop the functions and triggers
DROP TRIGGER IF EXISTS trigger_update_sales_overview ON sales;
DROP TRIGGER IF EXISTS trigger_update_item_status ON items;
DROP TRIGGER IF EXISTS trigger_distinct_access_levels ON store_owners;
DROP FUNCTION IF EXISTS initialize_transaction(bigint, NUMERIC(18, 2), varchar, varchar, varchar);
DROP FUNCTION IF EXISTS upsert_pending_funds(bigint, varchar, NUMERIC(18, 2));
DROP FUNCTION IF EXISTS process_transaction_completion(varchar, varchar, NUMERIC(10, 2), jsonb);
DROP FUNCTION IF EXISTS release_pending_funds(bigint);
DROP FUNCTION IF EXISTS update_sales_overview();
DROP FUNCTION IF EXISTS reduce_sales_overview(bigint, bigint, bigint);
DROP FUNCTION IF EXISTS upsert_cart_item(bigint, bigint, bigint, int);
DROP FUNCTION IF EXISTS update_item_status_on_out_of_stock();
DROP FUNCTION IF EXISTS fn_distinct_access_levels();
DROP FUNCTION IF EXISTS create_order(bigint, int, bigint, bigint, bigint, NUMERIC(10, 2), varchar, varchar);
DROP FUNCTION IF EXISTS create_sale(bigint, bigint, bigint, bigint, bigint);
DROP FUNCTION IF EXISTS create_review(bigint, bigint, bigint, NUMERIC(2, 1), varchar, TEXT, BOOLEAN);
DROP FUNCTION IF EXISTS delete_expired_sessions();
DROP FUNCTION IF EXISTS get_stores_by_user(bigint);

DROP FUNCTION IF EXISTS release_pending_funds;
DROP FUNCTION IF EXISTS upsert_pending_funds;

-- Drop the foreign key constraints
ALTER TABLE "pending_transaction_funds" DROP CONSTRAINT IF EXISTS "valid_account_type";
ALTER TABLE "pending_transaction_funds" DROP CONSTRAINT IF EXISTS "pending_transaction_funds_store_id_fkey";

ALTER TABLE "transactions" DROP CONSTRAINT IF EXISTS "valid_transaction_status";
ALTER TABLE "transactions" DROP CONSTRAINT IF EXISTS "valid_transaction_provider";
ALTER TABLE "transactions" DROP CONSTRAINT IF EXISTS "transactions_customer_id_fkey";

ALTER TABLE "fiat_accounts" DROP CONSTRAINT IF EXISTS "fiat_accounts_store_id_fkey";
ALTER TABLE "crypto_accounts" DROP CONSTRAINT IF EXISTS "crypto_accounts_store_id_fkey";

ALTER TABLE "cart_items" DROP CONSTRAINT IF EXISTS cart_items_cart_id_fkey;
ALTER TABLE "cart_items" DROP CONSTRAINT IF EXISTS cart_items_item_id_fkey;
ALTER TABLE "cart_items" DROP CONSTRAINT IF EXISTS unique_item_in_cart;
ALTER TABLE "carts" DROP CONSTRAINT IF EXISTS carts_user_id_fkey;

ALTER TABLE "review_likes" DROP CONSTRAINT IF EXISTS review_likes_review_id_fkey;

ALTER TABLE "store_owners" DROP CONSTRAINT IF EXISTS store_owners_user_id_fkey;
ALTER TABLE "store_owners" DROP CONSTRAINT IF EXISTS store_owners_store_id_fkey;

ALTER TABLE "items" DROP CONSTRAINT IF EXISTS items_store_id_fkey;

ALTER TABLE "sessions" DROP CONSTRAINT IF EXISTS sessions_user_id_fkey;

-- Drop the tables
DROP TABLE IF EXISTS "pending_transaction_funds";
DROP TABLE IF EXISTS "transactions";
DROP TABLE IF EXISTS "fiat_accounts";
DROP TABLE IF EXISTS "crypto_accounts";
DROP TABLE IF EXISTS "cart_items";
DROP TABLE IF EXISTS "carts";
DROP TABLE IF EXISTS "reviews";
DROP TABLE IF EXISTS "review_likes";
DROP TABLE IF EXISTS "sales_overview";
DROP TABLE IF EXISTS "sales";
DROP TABLE IF EXISTS "orders";
DROP TABLE IF EXISTS "store_audit_trail";
DROP TABLE IF EXISTS "store_owners";
DROP TABLE IF EXISTS "item_ratings";
DROP TABLE IF EXISTS "items";
DROP TABLE IF EXISTS "stores";
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "users";
