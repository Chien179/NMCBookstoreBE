-- name: GetBookCart :one
SELECT * FROM books_carts
WHERE id = $1 LIMIT 1;

-- name: ListBooksCartsByCartID :many
SELECT * FROM books_carts
WHERE carts_id = $1
ORDER BY id;

-- name: ListBooksCartsByBookID :many
SELECT * FROM books_carts
WHERE books_id = $1
ORDER BY id;

-- name: CreateBookCart :one
INSERT INTO books_carts (
  books_id,
  carts_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteBookCart :exec
DELETE FROM books_carts
WHERE id = $1;