-- name: CreatePayment :one
INSERT INTO payments (
  username,
  order_id,
  shipping_id,
  subtotal,
  status
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;