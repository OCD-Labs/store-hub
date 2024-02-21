-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  token,
  scope,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: CheckSessionExists :one
SELECT EXISTS(SELECT 1 FROM sessions WHERE token = $1 AND scope = $2) AS session_exists;

-- name: DeleteExpiredSession :exec
SELECT delete_expired_sessions();