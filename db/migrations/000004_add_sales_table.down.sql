ALTER TABLE "sales" DROP CONSTRAINT IF EXISTS sales_store_id_fkey;
ALTER TABLE "sales" DROP CONSTRAINT IF EXISTS sales_item_id_fkey;
ALTER TABLE "sales" DROP CONSTRAINT IF EXISTS sales_customer_id_fkey;
ALTER TABLE "sales" DROP CONSTRAINT IF EXISTS sales_seller_id_fkey;
ALTER TABLE "sales" DROP CONSTRAINT IF EXISTS sales_order_id_fkey;

DROP TABLE IF EXISTS "sales";