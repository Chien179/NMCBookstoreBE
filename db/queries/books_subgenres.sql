-- name: GetBookSubgenre :one
SELECT * FROM books_subgenres
WHERE id = $1 LIMIT 1;

-- name: ListBooksSubgenres :many
SELECT * FROM books_subgenres
WHERE subgenres_id = $1
ORDER BY id;

-- name: CreateBookSubgenre :one
INSERT INTO books_subgenres (
  books_id,
  subgenres_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteBookSubgenre :exec
DELETE FROM books_subgenres
WHERE id = $1;