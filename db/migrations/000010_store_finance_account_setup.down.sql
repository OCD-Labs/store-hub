-- Drop foreign key constraints first
ALTER TABLE "transactions" DROP CONSTRAINT IF EXISTS "fk_to_crypto_account";
ALTER TABLE "transactions" DROP CONSTRAINT IF EXISTS "fk_from_crypto_account";
ALTER TABLE "transactions" DROP CONSTRAINT IF EXISTS "fk_to_fiat_account";
ALTER TABLE "transactions" DROP CONSTRAINT IF EXISTS "fk_from_fiat_account";

-- Drop the transactions table
DROP TABLE IF EXISTS "transactions";

-- Drop foreign keys for fiat and crypto accounts
ALTER TABLE "fiat_accounts" DROP CONSTRAINT IF EXISTS "fiat_accounts_store_id_fkey";
ALTER TABLE "crypto_accounts" DROP CONSTRAINT IF EXISTS "crypto_accounts_store_id_fkey";

-- Drop the fiat_accounts and crypto_accounts tables
DROP TABLE IF EXISTS "fiat_accounts";
DROP TABLE IF EXISTS "crypto_accounts";

-- Drop get_stores_by_user function.
DROP FUNCTION IF EXISTS get_stores_by_user(bigint);
