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

-- Sales Overview Table
CREATE TABLE "sales_overview" (
  "id" bigserial PRIMARY KEY,
  "number_of_sales" bigint NOT NULL DEFAULT 0,
  "sales_percentage" NUMERIC(6, 4) NOT NULL DEFAULT 0,
  "revenue" NUMERIC(10, 2) NOT NULL DEFAULT 0,
  "item_id" bigint NOT NULL,
  "store_id" bigint NOT NULL
);

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
ALTER TABLE cart_items 
ADD CONSTRAINT unique_item_in_cart UNIQUE (cart_id, item_id);

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

-- Review Likes Table
CREATE TABLE "review_likes" (
  "id" bigserial PRIMARY KEY,
  "review_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  liked BOOLEAN NOT NULL,
  UNIQUE (review_id, user_id)
);
ALTER TABLE "review_likes" ADD FOREIGN KEY ("review_id") REFERENCES "reviews" ("id") ON DELETE CASCADE;

-- Store Audit Trail Table
CREATE TABLE "store_audit_trail" (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "action" varchar NOT NULL,
  "details" jsonb,
  "timestamp" timestamptz NOT NULL DEFAULT (now())
);


-- Functions and Triggers

-- Function to update the sales_overview table
CREATE OR REPLACE FUNCTION update_sales_overview()
RETURNS TRIGGER AS $$
DECLARE
    orderQty int;
    itemPrice NUMERIC(10, 2);
    supplyQuantity NUMERIC;
    v_item_exists bool;
    v_store_exists bool;
BEGIN
    -- Check if the item exists
    SELECT EXISTS(SELECT 1 FROM items WHERE id = NEW.item_id) INTO v_item_exists;
    IF NOT v_item_exists THEN
        RAISE EXCEPTION 'Item with ID % does not exist', NEW.item_id;
    END IF;

    -- Check if the store exists
    SELECT EXISTS(SELECT 1 FROM stores WHERE id = NEW.store_id) INTO v_store_exists;
    IF NOT v_store_exists THEN
        RAISE EXCEPTION 'Store with ID % does not exist', NEW.store_id;
    END IF;

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

-- Create the reduce_sales_overview function
CREATE OR REPLACE FUNCTION reduce_sales_overview(item_id_arg bigint, store_id_arg bigint, order_id_arg bigint)
RETURNS void AS $$
DECLARE
    orderQty int;
    itemPrice NUMERIC(10, 2);
    supplyQuantity NUMERIC;
BEGIN
    -- Fetch the order quantity and item details for the sale
    SELECT o.order_quantity, i.price, i.supply_quantity 
    INTO orderQty, itemPrice, supplyQuantity 
    FROM orders o
    JOIN items i ON i.id = item_id_arg
    WHERE o.id = order_id_arg;

    -- Reduce the number of sales by the order quantity
    UPDATE sales_overview
    SET 
        number_of_sales = GREATEST(number_of_sales - orderQty, 0), -- Ensure it doesn't go below 0
        sales_percentage = (GREATEST(number_of_sales - orderQty, 0) / supplyQuantity) * 100,
        revenue = GREATEST(number_of_sales - orderQty, 0) * itemPrice
    WHERE item_id = item_id_arg AND store_id = store_id_arg;
END;
$$ LANGUAGE plpgsql;

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

-- Function to create an order, with checks to ensure the supposed referenced rows 
-- exist in their respective tables.
CREATE OR REPLACE FUNCTION create_order(
    p_item_id bigint,
    p_order_quantity int,
    p_buyer_id bigint,
    p_seller_id bigint,
    p_store_id bigint,
    p_delivery_fee NUMERIC(10, 2),
    p_payment_channel varchar,
    p_payment_method varchar
)
RETURNS orders AS $$
DECLARE
    v_item_exists bool;
    v_buyer_exists bool;
    v_seller_exists bool;
    v_store_exists bool;
    v_result orders%ROWTYPE;
BEGIN
    -- Check if item exists
    SELECT EXISTS(SELECT 1 FROM items WHERE id = p_item_id) INTO v_item_exists;
    IF NOT v_item_exists THEN
        RAISE EXCEPTION 'Item with ID % does not exist', p_item_id;
    END IF;

    -- Check if buyer exists
    SELECT EXISTS(SELECT 1 FROM users WHERE id = p_buyer_id) INTO v_buyer_exists;
    IF NOT v_buyer_exists THEN
        RAISE EXCEPTION 'Buyer with ID % does not exist', p_buyer_id;
    END IF;

    -- Check if seller exists
    SELECT EXISTS(SELECT 1 FROM users WHERE id = p_seller_id) INTO v_seller_exists;
    IF NOT v_seller_exists THEN
        RAISE EXCEPTION 'Seller with ID % does not exist', p_seller_id;
    END IF;

    -- Check if store exists
    SELECT EXISTS(SELECT 1 FROM stores WHERE id = p_store_id) INTO v_store_exists;
    IF NOT v_store_exists THEN
        RAISE EXCEPTION 'Store with ID % does not exist', p_store_id;
    END IF;

    -- If all checks pass, insert the order
    INSERT INTO orders (
        item_id,
        order_quantity,
        buyer_id,
        seller_id,
        store_id,
        delivery_fee,
        payment_channel,
        payment_method
    ) VALUES (
        p_item_id, p_order_quantity, p_buyer_id, p_seller_id, p_store_id, p_delivery_fee, p_payment_channel, p_payment_method
    ) RETURNING * INTO v_result;

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

-- Function to create an sale, with checks to ensure the supposed referenced rows 
-- exist in their respective tables.
CREATE OR REPLACE FUNCTION create_sale(
    p_store_id bigint,
    p_item_id bigint,
    p_customer_id bigint,
    p_seller_id bigint,
    p_order_id bigint
)
RETURNS sales AS $$
DECLARE
    v_store_exists bool;
    v_item_exists bool;
    v_customer_exists bool;
    v_seller_exists bool;
    v_order_exists bool;
    v_result sales%ROWTYPE;
BEGIN
    -- Check if store exists
    SELECT EXISTS(SELECT 1 FROM stores WHERE id = p_store_id) INTO v_store_exists;
    IF NOT v_store_exists THEN
        RAISE EXCEPTION 'Store with ID % does not exist', p_store_id;
    END IF;

    -- Check if item exists
    SELECT EXISTS(SELECT 1 FROM items WHERE id = p_item_id) INTO v_item_exists;
    IF NOT v_item_exists THEN
        RAISE EXCEPTION 'Item with ID % does not exist', p_item_id;
    END IF;

    -- Check if customer exists
    SELECT EXISTS(SELECT 1 FROM users WHERE id = p_customer_id) INTO v_customer_exists;
    IF NOT v_customer_exists THEN
        RAISE EXCEPTION 'Customer with ID % does not exist', p_customer_id;
    END IF;

    -- Check if seller exists
    SELECT EXISTS(SELECT 1 FROM users WHERE id = p_seller_id) INTO v_seller_exists;
    IF NOT v_seller_exists THEN
        RAISE EXCEPTION 'Seller with ID % does not exist', p_seller_id;
    END IF;

    -- Check if order exists
    SELECT EXISTS(SELECT 1 FROM orders WHERE id = p_order_id) INTO v_order_exists;
    IF NOT v_order_exists THEN
        RAISE EXCEPTION 'Order with ID % does not exist', p_order_id;
    END IF;

    -- If all checks pass, insert the sale
    INSERT INTO sales (
        store_id,
        item_id,
        customer_id,
        seller_id,
        order_id
    ) VALUES (
        p_store_id, p_item_id, p_customer_id, p_seller_id, p_order_id
    ) RETURNING * INTO v_result;

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

-- Function to create a review
CREATE OR REPLACE FUNCTION create_review(
    p_store_id bigint,
    p_user_id bigint,
    p_item_id bigint,
    p_rating NUMERIC(2, 1),
    p_review_type varchar,
    p_comment TEXT,
    p_is_verified_purchase BOOLEAN
) RETURNS void AS $$
DECLARE
    v_store_exists bool;
    v_user_exists bool;
    v_item_exists bool;
    v_order_exists bool;
BEGIN
    -- Check if the store exists
    SELECT EXISTS(SELECT 1 FROM stores WHERE id = p_store_id) INTO v_store_exists;
    IF NOT v_store_exists THEN
        RAISE EXCEPTION 'Store with ID % does not exist', p_store_id;
    END IF;

    -- Check if the user exists
    SELECT EXISTS(SELECT 1 FROM users WHERE id = p_user_id) INTO v_user_exists;
    IF NOT v_user_exists THEN
        RAISE EXCEPTION 'User with ID % does not exist', p_user_id;
    END IF;

    -- Check if the item exists
    SELECT EXISTS(SELECT 1 FROM items WHERE id = p_item_id) INTO v_item_exists;
    IF NOT v_item_exists THEN
        RAISE EXCEPTION 'Item with ID % does not exist', p_item_id;
    END IF;

    -- If is_verified_purchase is true, check for a corresponding order
    IF p_is_verified_purchase THEN
        SELECT EXISTS(SELECT 1 FROM orders WHERE user_id = p_user_id AND item_id = p_item_id) INTO v_order_exists;
        IF NOT v_order_exists THEN
            RAISE EXCEPTION 'No verified purchase found for User ID % and Item ID %', p_user_id, p_item_id;
        END IF;
    END IF;

    -- Insert the review
    INSERT INTO reviews (
        store_id,
        user_id,
        item_id,
        rating,
        review_type,
        comment,
        is_verified_purchase
    ) VALUES (
        p_store_id,
        p_user_id,
        p_item_id,
        p_rating,
        p_review_type,
        p_comment,
        p_is_verified_purchase
    );
END;
$$ LANGUAGE plpgsql;

-- Function to delete expired session token
CREATE OR REPLACE FUNCTION delete_expired_sessions()
RETURNS void AS $$
BEGIN
  DELETE FROM sessions WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

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