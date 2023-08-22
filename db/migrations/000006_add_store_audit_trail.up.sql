ALTER TABLE store_owners ADD COLUMN is_primary boolean NOT NULL DEFAULT false;
ALTER TABLE store_owners DROP COLUMN access_level;
ALTER TABLE store_owners ADD COLUMN access_levels int ARRAY NOT NULL DEFAULT '{}';


CREATE TABLE "store_audit_trail" (
  "id" bigserial PRIMARY KEY,
  "store_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "action" varchar NOT NULL,
  "details" jsonb,
  "timestamp" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "store_audit_trail" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");
ALTER TABLE "store_audit_trail" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
