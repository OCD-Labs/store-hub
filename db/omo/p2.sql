-- UP Migration

-- Users Table
CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "account_id" varchar UNIQUE NOT NULL,
  "status" varchar NOT NULL,
  "about" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "socials" jsonb NOT NULL,
  "profile_image_url" varchar,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "is_active" boolean NOT NULL DEFAULT true,
  "is_email_verified" boolean NOT NULL DEFAULT false
);

-- Sessions Table
CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "token" varchar NOT NULL,
  "scope" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

-- Stores Table
CREATE TABLE "stores" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "store_account_id" varchar UNIQUE NOT NULL,
  "profile_image_url" varchar NOT NULL,
  "is_verified" boolean NOT NULL DEFAULT false,
  "category" varchar NOT NULL,
  "is_frozen" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Store Owners Table
CREATE TABLE "store_owners" (
  "user_id" bigint NOT NULL,
  "store_id" bigint NOT NULL,
  "access_levels" int[] NOT NULL DEFAULT '{}',
  "is_primary" boolean NOT NULL DEFAULT false,
  "added_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "store_owners" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "store_owners" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;

-- Items Table
CREATE TABLE "items" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "price" NUMERIC(10, 2) NOT NULL,
  "store_id" bigint NOT NULL,
  "image_urls" text[] NOT NULL,
  "category" varchar NOT NULL,
  "discount_percentage" NUMERIC(6, 4) NOT NULL,
  "supply_quantity" bigint NOT NULL,
  "extra" jsonb NOT NULL,
  "is_frozen" bool NOT NULL DEFAULT false,
  "currency" varchar NOT NULL DEFAULT 'NGN',
  "cover_img_url" varchar NOT NULL DEFAULT '',
  "status" VARCHAR(10) NOT NULL DEFAULT 'VISIBLE' CHECK (status IN ('VISIBLE', 'HIDDEN')),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "items" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;

-- Item Ratings Table
CREATE TABLE "item_ratings" (
  "user_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "rating" smallint NOT NULL,
  "comment" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "item_ratings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "item_ratings" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON DELETE CASCADE;

-- Orders Table
CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "delivery_status" varchar NOT NULL DEFAULT 'PENDING',
  "delivered_on" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "expected_delivery_date" timestamptz NOT NULL DEFAULT (now() + interval '3 days'),
  "item_id" bigint NOT NULL,
  "order_quantity" int NOT NULL,
  "buyer_id" bigint NOT NULL,
  "seller_id" bigint NOT NULL,
  "store_id" bigint NOT NULL,
  "delivery_fee" NUMERIC(10, 2) NOT NULL,
  "payment_channel" varchar NOT NULL,
  "payment_method" varchar NOT NULL,
  "is_reviewed" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "orders" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON DELETE CASCADE;
ALTER TABLE "orders" ADD FOREIGN KEY ("buyer_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "orders" ADD FOREIGN KEY ("seller_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "orders" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;

-- Sales Table
CREATE TABLE "sales" (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "customer_id" bigint NOT NULL,
  "seller_id" bigint NOT NULL,
  "order_id" bigint UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "sales" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;
ALTER TABLE "sales" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON DELETE CASCADE;
ALTER TABLE "sales" ADD FOREIGN KEY ("customer_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "sales" ADD FOREIGN KEY ("seller_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "sales" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE CASCADE;

-- Sales Overview Table
CREATE TABLE "sales_overview" (
  "id" bigserial PRIMARY KEY,
  "number_of_sales" bigint NOT NULL DEFAULT 0,
  "sales_percentage" NUMERIC(6, 4) NOT NULL DEFAULT 0,
  "revenue" NUMERIC(10, 2) NOT NULL DEFAULT 0,
  "item_id" bigint NOT NULL,
  "store_id" bigint NOT NULL
);
ALTER TABLE "sales_overview" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON DELETE CASCADE;
ALTER TABLE "sales_overview" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;

-- Carts Table
CREATE TABLE "carts" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL UNIQUE,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "carts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

-- Cart Items Table
CREATE TABLE "cart_items" (
  "id" bigserial PRIMARY KEY,
  "cart_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "store_id" bigint NOT NULL,
  "quantity" int NOT NULL DEFAULT 1,
  "added_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "cart_items" ADD FOREIGN KEY ("cart_id") REFERENCES "carts" ("id") ON DELETE CASCADE;
ALTER TABLE "cart_items" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON DELETE CASCADE;

-- Fiat Accounts Table
CREATE TABLE "fiat_accounts" (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "balance" NUMERIC(10, 2) NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "fiat_accounts" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;

-- Crypto Accounts Table
CREATE TABLE "crypto_accounts" (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "balance" NUMERIC(18, 8) NOT NULL,
  "wallet_address" varchar NOT NULL,
  "crypto_type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "crypto_accounts" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;

-- Transactions Table
CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "amount" NUMERIC(18, 2) NOT NULL,
  "from_account_id" bigint NOT NULL,
  "TO_account_id" bigint NOT NULL,
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
ALTER TABLE "transactions" ADD CONSTRAINT fk_from_fiat_account FOREIGN KEY (from_account_id) REFERENCES fiat_accounts(id) ON DELETE CASCADE;
ALTER TABLE "transactions" ADD CONSTRAINT fk_to_fiat_account FOREIGN KEY (to_account_id) REFERENCES fiat_accounts(id) ON DELETE CASCADE;
ALTER TABLE "transactions" ADD CONSTRAINT fk_from_crypto_account FOREIGN KEY (from_account_id) REFERENCES crypto_accounts(id) ON DELETE CASCADE;
ALTER TABLE "transactions" ADD CONSTRAINT fk_to_crypto_account FOREIGN KEY (to_account_id) REFERENCES crypto_accounts(id) ON DELETE CASCADE;

-- Reviews Table
CREATE TABLE "reviews" (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "rating" NUMERIC(2, 1) NOT NULL CHECK (rating >= 1 AND rating <= 5),
  "review_type" varchar NOT NULL,
  "comment" TEXT NOT NULL DEFAULT '',
  "is_verified_purchase" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "reviews" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;
ALTER TABLE "reviews" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "reviews" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON DELETE CASCADE;

-- Review Likes Table
CREATE TABLE "review_likes" (
  "id" bigserial PRIMARY KEY,
  "review_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  liked BOOLEAN NOT NULL,
  UNIQUE (review_id, user_id)
);
ALTER TABLE "review_likes" ADD FOREIGN KEY ("review_id") REFERENCES "reviews" ("id") ON DELETE CASCADE;
ALTER TABLE "review_likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

-- Store Audit Trail Table
CREATE TABLE "store_audit_trail" (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "action" varchar NOT NULL,
  "details" jsonb,
  "timestamp" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "store_audit_trail" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;
ALTER TABLE "store_audit_trail" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;


-- Functions and Triggers

-- Function to update the sales_overview table
CREATE OR REPLACE FUNCTION update_sales_overview()
RETURNS TRIGGER AS $$
DECLARE
    orderQty int;
    itemPrice NUMERIC(10, 2);
    supplyQuantity NUMERIC;
BEGIN
    -- Fetch the order quantity and item details for the sale
    SELECT o.order_quantity, i.price, i.supply_quantity 
    INTO orderQty, itemPrice, supplyQuantity 
    FROM orders o
    JOIN items i ON o.item_id = i.id
    WHERE o.id = NEW.order_id;

    -- Check if the item and store combination already exists in sales_overview
    IF EXISTS (SELECT 1 FROM sales_overview WHERE item_id = NEW.item_id AND store_id = NEW.store_id) THEN
        -- Update the existing record
        UPDATE sales_overview
        SET 
            number_of_sales = number_of_sales + orderQty,
            sales_percentage = ((number_of_sales + orderQty) / supplyQuantity) * 100,
            revenue = (number_of_sales + orderQty) * itemPrice
        WHERE item_id = NEW.item_id AND store_id = NEW.store_id;
    ELSE
        -- Insert a new record
        INSERT INTO sales_overview (number_of_sales, sales_percentage, revenue, item_id, store_id)
        VALUES (
            orderQty,
            (orderQty / supplyQuantity) * 100,
            orderQty * itemPrice,
            NEW.item_id,
            NEW.store_id
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the update_sales_overview function after a sale is inserted
CREATE TRIGGER trigger_update_sales_overview
AFTER INSERT ON sales
FOR EACH ROW
EXECUTE FUNCTION update_sales_overview();

-- Function to ensure distinct access levels for store owners
CREATE OR REPLACE FUNCTION fn_distinct_access_levels()
RETURNS TRIGGER AS $$
BEGIN
    -- Ensure distinct access levels
    NEW.access_levels := ARRAY(
        SELECT DISTINCT unnest(NEW.access_levels)
    );

    -- If access_levels is empty, delete the row
    IF array_length(NEW.access_levels, 1) IS NULL THEN
        DELETE FROM store_owners WHERE user_id = NEW.user_id AND store_id = NEW.store_id;
        RETURN NULL; -- Important to return NULL for DELETE operation in BEFORE trigger
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to ensure distinct access levels for store owners
CREATE TRIGGER trigger_distinct_access_levels
BEFORE INSERT OR UPDATE ON store_owners
FOR EACH ROW
EXECUTE FUNCTION fn_distinct_access_levels();

-- Function to hide item from the storefront if out of stock
CREATE OR REPLACE FUNCTION update_item_status_on_out_of_stock()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the supply_quantity of the item is 0
    IF NEW.supply_quantity = 0 THEN
        -- Set the status to 'HIDDEN'
        NEW.status := 'HIDDEN';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update item status on items UPDATE
CREATE TRIGGER trigger_update_item_status
BEFORE UPDATE ON items
FOR EACH ROW
WHEN (NEW.supply_quantity = 0)
EXECUTE FUNCTION update_item_status_on_out_of_stock();

-- Function to upsert cart items
CREATE OR REPLACE FUNCTION upsert_cart_item(p_cart_id bigint, p_item_id bigint, p_store_id bigint, p_quantity int DEFAULT 1)
RETURNS cart_items AS $$
DECLARE
    v_existing_quantity int;
    v_result cart_items%ROWTYPE;
BEGIN
    SELECT quantity INTO v_existing_quantity FROM cart_items WHERE cart_id = p_cart_id AND item_id = p_item_id;

    IF FOUND THEN
        UPDATE cart_items 
        SET quantity = v_existing_quantity + p_quantity
        WHERE cart_id = p_cart_id AND item_id = p_item_id
        RETURNING * INTO v_result;
    ELSE
        INSERT INTO cart_items (cart_id, item_id, store_id, quantity)
        VALUES (p_cart_id, p_item_id, p_store_id, p_quantity)
        RETURNING * INTO v_result;
    END IF;

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

