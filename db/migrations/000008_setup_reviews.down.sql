-- Drop the foreign key constraints
ALTER TABLE "review_likes" DROP CONSTRAINT IF EXISTS review_likes_review_id_fkey;
ALTER TABLE "review_likes" DROP CONSTRAINT IF EXISTS review_likes_user_id_fkey;

ALTER TABLE "reviews" DROP CONSTRAINT IF EXISTS reviews_store_id_fkey;
ALTER TABLE "reviews" DROP CONSTRAINT IF EXISTS reviews_user_id_fkey;
ALTER TABLE "reviews" DROP CONSTRAINT IF EXISTS reviews_item_id_fkey;

-- Drop the tables
DROP TABLE IF EXISTS "review_likes";
DROP TABLE IF EXISTS "reviews";

ALTER TABLE orders DROP COLUMN IF EXISTS is_reviewed;
ALTER TABLE items DROP COLUMN IF EXISTS status;
