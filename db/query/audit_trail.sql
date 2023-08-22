-- name: LogAction :exec
INSERT INTO store_audit_trail (
  store_id, user_id, action, details
) VALUES (
  $1, $2, $3, $4
);