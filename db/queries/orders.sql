-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: CreateOrder :one
INSERT INTO orders (
    users_id
) VALUES (
  $1
)
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;