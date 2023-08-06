ALTER TABLE "stores" ADD COLUMN "currency" varchar NOT NULL DEFAULT 'USD';
ALTER TABLE "items" ADD COLUMN "cover_img_url" varchar NOT NULL DEFAULT 'http://res.cloudinary.com/duxnx9n5t/image/upload/v1690708857/raczxa9rcxlo35odjp5x.png';

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "delivery_status" varchar NOT NULL,
  "delivered_on" timestamptz NOT NULL DEFAULT '0001-01-01T00:00:00Z',
  "expected_delivery_date" timestamptz NOT NULL DEFAULT (now() + interval '3 days'),
  "item_id" bigint NOT NULL,
  "order_quantity" int NOT NULL,
  "buyer_id" bigint NOT NULL,
  "seller_id" bigint NOT NULL,
  "store_id" bigint NOT NULL,
  "delivery_fee" NUMERIC(10, 2) NOT NULL,
  "payment_channel" varchar NOT NULL,
  "payment_method" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "orders" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("buyer_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("seller_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");
