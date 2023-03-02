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

-- name: ListOdersByUserID :many
SELECT * FROM orders
WHERE users_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;