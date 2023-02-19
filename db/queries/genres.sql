-- name: GetGenre :one
SELECT * FROM genres
WHERE id = $1 LIMIT 1;

-- name: ListGenres :many
SELECT * FROM genres
ORDER BY id;

-- name: CreateGenre :one
INSERT INTO genres (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: DeleteGenre :exec
DELETE FROM genres
WHERE id = $1;

-- name: UpdateGenre :one
UPDATE genres
SET name = $2
WHERE id = $1
RETURNING *;