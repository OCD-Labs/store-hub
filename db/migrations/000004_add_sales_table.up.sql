CREATE TABLE "sales" (
  "id" bigserial PRIMARY KEY,
  "item_id" bigint NOT NULL,
  "customer_id" bigint NOT NULL,
  "order_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sales" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("customer_id") REFERENCES "users" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");