-- name: GetCart :one
SELECT * FROM carts
WHERE id = $1 LIMIT 1;

-- name: CreateCart :one
INSERT INTO carts (
    users_id
) VALUES (
  $1
)
RETURNING *;