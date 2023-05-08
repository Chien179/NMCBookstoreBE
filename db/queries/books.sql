-- name: GetBook :one
SELECT * FROM books
WHERE id = $1 LIMIT 1;

-- name: ListBooks :many
SELECT
    (SELECT (COUNT(*)/sqlc.arg('limit'))
     FROM books) 
     as total_page, 
    (SELECT JSON_AGG(t.*) FROM (
        SELECT * FROM books
        ORDER BY id
        LIMIT sqlc.arg('limit')
        OFFSET sqlc.arg('offset')
    ) AS t) AS books ;

-- name: ListTop10TheBestBooks :many
SELECT * FROM books
ORDER BY rating DESC
LIMIT 10;

-- name: ListTop10NewestBooks :many
SELECT * FROM books
ORDER BY created_at DESC
LIMIT 10;

-- name: CreateBook :one
INSERT INTO books (
  name,
  price,
  image,
  description,
  author,
  publisher,
  quantity
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;

-- name: UpdateBook :one
UPDATE books
SET name = COALESCE(sqlc.narg(name), name),
  price = COALESCE(sqlc.narg(price), price),
  image = COALESCE(sqlc.narg(image), image),
  description = COALESCE(sqlc.narg(description), description),
  author = COALESCE(sqlc.narg(author), author),
  publisher = COALESCE(sqlc.narg(publisher), publisher),
  quantity = COALESCE(sqlc.narg(quantity), quantity)
WHERE id = sqlc.arg(id)
RETURNING *;