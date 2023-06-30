-- name: CreateStore :one
INSERT INTO stores (
  name,
  description,
  profile_image_url,
  category
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetStoreByID :one
SELECT 
  s.*, 
  json_agg(json_build_object(
      'user', json_build_object('id', u.id, 'first_name', u.first_name, 'last_name', u.last_name, 'email', u.email),
      'store_owners', json_build_object('user_id', so.user_id, 'store_id', so.store_id, 'added_at', so.added_at)
  )) AS owners
FROM 
  stores AS s
JOIN 
  store_owners AS so ON s.id = so.store_id
JOIN 
  users AS u ON so.user_id = u.id
WHERE 
  s.id = sqlc.arg(store_id)
GROUP BY 
  s.id;