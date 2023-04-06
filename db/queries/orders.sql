-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: CreateOrder :one
INSERT INTO orders (
    username
) VALUES (
  $1
)
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: ListOdersByUserName :many
SELECT * FROM orders
WHERE username = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateStatus :one
UPDATE orders
SET 
  status = $2
WHERE 
  id = $1
RETURNING *;