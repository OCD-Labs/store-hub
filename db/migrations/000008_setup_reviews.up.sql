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