DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'stores' AND column_name = 'currency') THEN
        ALTER TABLE "stores" DROP COLUMN "currency";
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'items' AND column_name = 'cover_img_url') THEN
        ALTER TABLE "items" DROP COLUMN "cover_img_url";
    END IF;
END $$;

ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS orders_item_id_fkey;
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS orders_buyer_id_fkey;
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS orders_store_id_fkey;

DROP TABLE IF EXISTS "orders";