-- carts table
CREATE TABLE carts (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL UNIQUE,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE carts
ADD FOREIGN KEY ("user_id") REFERENCES users(id) ON DELETE CASCADE;

-- cart_items table
CREATE TABLE cart_items (
  "id" bigserial PRIMARY KEY,
  "cart_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "store_id" bigint NOT NULL,
  "quantity" int NOT NULL DEFAULT 1,
  "added_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE cart_items
ADD FOREIGN KEY ("cart_id") REFERENCES carts(id) ON DELETE CASCADE,
ADD FOREIGN KEY ("item_id") REFERENCES items(id) ON DELETE CASCADE;

ALTER TABLE cart_items 
ADD CONSTRAINT unique_item_in_cart UNIQUE (cart_id, item_id);

-- The upsert_cart_item function updates or create a new cart_items row
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

