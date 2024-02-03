CREATE TABLE "sales" (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "customer_id" bigint NOT NULL,
  "seller_id" bigint NOT NULL,
  "order_id" bigint UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sales" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("customer_id") REFERENCES "users" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("seller_id") REFERENCES "users" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
