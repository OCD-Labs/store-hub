-- Drop the function
DROP FUNCTION IF EXISTS upsert_cart_item(bigint, bigint, bigint, int);

-- Drop foreign key constraints from cart_items table
ALTER TABLE cart_items
DROP CONSTRAINT IF EXISTS cart_items_cart_id_fkey,
DROP CONSTRAINT IF EXISTS cart_items_item_id_fkey;

-- Drop unique constraint from cart_items table
ALTER TABLE cart_items 
DROP CONSTRAINT IF EXISTS unique_item_in_cart;

-- Drop cart_items table
DROP TABLE IF EXISTS cart_items;

-- Drop foreign key constraint from carts table
ALTER TABLE carts
DROP CONSTRAINT IF EXISTS carts_user_id_fkey;

-- Drop carts table
DROP TABLE IF EXISTS carts;
