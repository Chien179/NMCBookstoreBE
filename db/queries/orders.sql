-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: GetOrderToPayment :one
SELECT * FROM orders
WHERE username = $1 
LIMIT 1;

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

-- name: UpdateOrder :one
UPDATE orders
SET 
  status = COALESCE(sqlc.narg(status), status),
  sub_amount = COALESCE(sqlc.narg(sub_amount), sub_amount),
  sub_total = COALESCE(sqlc.narg(sub_total), sub_total)
WHERE 
  id = sqlc.arg(id)
RETURNING *;