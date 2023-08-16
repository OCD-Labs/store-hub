-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-08-16T12:42:24.657Z

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

CREATE TABLE "stores" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "profile_image_url" varchar NOT NULL,
  "store_account_id" varchar UNIQUE NOT NULL,
  "is_verified" boolean NOT NULL DEFAULT false,
  "category" varchar NOT NULL,
  "is_frozen" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "store_owners" (
  "user_id" bigint NOT NULL,
  "store_id" bigint NOT NULL,
  "access_level" smallint NOT NULL,
  "added_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "items" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "price" "NUMERIC(10, 2)" NOT NULL,
  "store_id" bigint NOT NULL,
  "image_urls" text[] NOT NULL,
  "category" varchar NOT NULL,
  "discount_percentage" "NUMERIC(6, 4)" NOT NULL,
  "supply_quantity" bigint NOT NULL,
  "extra" jsonb NOT NULL,
  "is_frozen" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "item_ratings" (
  "user_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "rating" char NOT NULL,
  "comment" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

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
  "delivery_fee" "NUMERIC(10, 2)" NOT NULL,
  "payment_channel" varchar NOT NULL,
  "payment_method" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sales" (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "customer_id" bigint NOT NULL,
  "seller_id" bigint NOT NULL,
  "order_id" bigint UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "store_owners" ("user_id", "store_id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "store_owners" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "store_owners" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

ALTER TABLE "items" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

ALTER TABLE "item_ratings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "item_ratings" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("buyer_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("seller_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("customer_id") REFERENCES "users" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("seller_id") REFERENCES "users" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
