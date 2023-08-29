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

CREATE TABLE "review_likes" (
  "id" bigserial PRIMARY KEY,
  "review_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  liked BOOLEAN NOT NULL,
  UNIQUE (review_id, user_id)
);

-- New ALTER statements for reviews and review_likes
ALTER TABLE "reviews" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;
ALTER TABLE "reviews" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "reviews" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id");

ALTER TABLE "review_likes" ADD FOREIGN KEY ("review_id") REFERENCES "reviews" ("id") ON DELETE CASCADE;
ALTER TABLE "review_likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "orders" ADD COLUMN is_reviewed BOOLEAN NOT NULL DEFAULT FALSE;

-- Add a status column to the items table with possible values 'VISIBLE' and 'HIDDEN'
ALTER TABLE items ADD COLUMN status VARCHAR(10) NOT NULL DEFAULT 'VISIBLE' CHECK (status IN ('VISIBLE', 'HIDDEN'));

-- Define function to hide item from the storefront if out of stock.
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

-- Trigger the update_item_status_on_out_of_stock on items UPDATE.
CREATE TRIGGER trigger_update_item_status
BEFORE UPDATE ON items
FOR EACH ROW
WHEN (NEW.supply_quantity = 0)
EXECUTE FUNCTION update_item_status_on_out_of_stock();
