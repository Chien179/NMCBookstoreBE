-- name: GetCart :one
SELECT *
FROM carts
WHERE id = $1
LIMIT 1;
-- name: CreateCart :one
INSERT INTO carts (books_id, username, amount, total)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: ListCartsByUsername :many
SELECT *
FROM carts
WHERE username = $1
ORDER BY id;
-- name: DeleteCart :exec
DELETE FROM carts
WHERE id = $1
  AND username = $2;
-- name: UpdateAmount :one
UPDATE carts
SET amount = $2,
  total = $3
WHERE id = $1
RETURNING *;