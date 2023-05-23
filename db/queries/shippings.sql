-- name: CreateShipping :one
INSERT INTO shippings (
  to_address,
  total
) VALUES (
$1, $2
)
RETURNING *;