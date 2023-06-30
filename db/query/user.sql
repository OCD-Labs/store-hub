-- name: CreateUser :one
INSERT INTO users (
  first_name,
  last_name,
  account_id,
  status,
  hashed_password,
  about,
  email,
  socials,
  profile_image_url
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = sqlc.arg(user_id) LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = sqlc.arg(user_email) LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  last_name = COALESCE(sqlc.narg(last_name), last_name),
  email = COALESCE(sqlc.narg(email), email),
  is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified),
  is_active = COALESCE(sqlc.narg(is_active), is_active),
  profile_image_url = COALESCE(sqlc.narg(profile_image_url), profile_image_url),
  socials = COALESCE(sqlc.narg(socials), socials),
  status = COALESCE(sqlc.narg(status), status),
  about = COALESCE(sqlc.narg(about), about)
WHERE 
  id = sqlc.narg(id) OR email = sqlc.narg(email)
RETURNING *;