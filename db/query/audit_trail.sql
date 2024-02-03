-- name: LogAction :exec
INSERT INTO store_audit_trail (
  -- Turn this into a function instead
  store_id, user_id, action, details
) VALUES (
  $1, $2, $3, $4
);