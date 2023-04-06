-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: ListTransactionsByOrderID :many
SELECT * FROM transactions
WHERE orders_id = $1
ORDER BY id;

-- name: CreateTransaction :one
INSERT INTO transactions (
    orders_id,
    books_id,
    amount,
    total
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;