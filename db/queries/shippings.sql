-- name: CreateShipping :one
INSERT INTO shippings (
  from_address,
  to_address,
  total
) VALUES (
  $1, $2, $3
)
RETURNING *;