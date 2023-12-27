-- name: CreatePayment :one
INSERT INTO payments (
    id,
    username,
    order_id,
    shipping_id,
    subtotal,
    status,
    created_at
  )
VALUES ($1, $2, $3, $4, $5, $6, NOW())
RETURNING *;