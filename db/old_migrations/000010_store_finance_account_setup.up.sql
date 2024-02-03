CREATE TABLE fiat_accounts (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "balance" NUMERIC(10, 2) NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE crypto_accounts (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "balance" NUMERIC(18, 8) NOT NULL,
  "wallet_address" varchar NOT NULL,
  "crypto_type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create foreign keys for both fiat and crypto accounts to stores table
ALTER TABLE "fiat_accounts" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;
ALTER TABLE "crypto_accounts" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;

CREATE TABLE transactions (
  "id" bigserial PRIMARY KEY,
  "amount" NUMERIC(18, 2) NOT NULL,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "payment_channel" varchar NOT NULL,
  "transaction_fee" NUMERIC(10, 2) NOT NULL,
  "conversion_fee" NUMERIC(10, 2) NOT NULL,
  "description" TEXT NOT NULL,
  "transaction_type" varchar NOT NULL,
  "transaction_ref_id" varchar NOT NULL,
  "status" varchar NOT NULL,
  "account_balance_snapshot" NUMERIC(18, 2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create foreign keys for both fiat and crypto accounts
ALTER TABLE "transactions"
ADD CONSTRAINT fk_from_fiat_account FOREIGN KEY (from_account_id) REFERENCES fiat_accounts(id) ON DELETE CASCADE;

ALTER TABLE "transactions"
ADD CONSTRAINT fk_to_fiat_account FOREIGN KEY (to_account_id) REFERENCES fiat_accounts(id) ON DELETE CASCADE;

ALTER TABLE "transactions"
ADD CONSTRAINT fk_from_crypto_account FOREIGN KEY (from_account_id) REFERENCES crypto_accounts(id) ON DELETE CASCADE;

ALTER TABLE "transactions"
ADD CONSTRAINT fk_to_crypto_account FOREIGN KEY (to_account_id) REFERENCES crypto_accounts(id) ON DELETE CASCADE;

-- Function: get_stores_by_user
-- Description: Fetches stores owned by a specific user and aggregates details of all owners for each store.
CREATE OR REPLACE FUNCTION get_stores_by_user(p_user_id bigint)
RETURNS TABLE (
    store_id bigint,
    store_name varchar,
    store_description varchar,
    store_image varchar,
    store_account_id varchar,
    is_verified boolean,
    category varchar,
    is_frozen boolean,
    store_created_at timestamptz,
    store_owners jsonb
) AS $$
DECLARE
    v_store record;
    v_owners jsonb;
BEGIN
    -- Loop through stores owned by the specified user
    FOR v_store IN (
        SELECT 
            s.id,
            s.name,
            s.description,
            s.profile_image_url,
            s.store_account_id,
            s.is_verified,
            s.category,
            s.is_frozen,
            s.created_at
        FROM 
            stores s
        JOIN 
            store_owners so ON s.id = so.store_id
        WHERE 
            so.user_id = p_user_id
    )
    LOOP
        -- Fetch and aggregate details of all owners for the current store
        SELECT 
            json_agg(json_build_object(
                'account_id', u.account_id,
                'profile_img_url', u.profile_image_url,
                'email', u.email,
                'access_levels', so.access_levels,
                'is_original_owner', so.is_primary,
                'added_at', so.added_at
            )) 
        INTO v_owners
        FROM 
            store_owners so
        JOIN 
            users AS u ON so.user_id = u.id
        WHERE 
            so.store_id = v_store.id;

        -- Populate the result set
        store_id := v_store.id;
        store_name := v_store.name;
        store_description := v_store.description;
        store_image := v_store.profile_image_url;
        store_account_id := v_store.store_account_id;
        is_verified := v_store.is_verified;
        category := v_store.category;
        is_frozen := v_store.is_frozen;
        store_created_at := v_store.created_at;
        store_owners := v_owners;

        -- Return the current row
        RETURN NEXT;
    END LOOP;
END;
$$ LANGUAGE plpgsql;
