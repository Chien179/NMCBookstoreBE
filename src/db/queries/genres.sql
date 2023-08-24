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

-- name: SoftDeleteGenre :one
UPDATE genres
SET is_deleted = true
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateGenre :one
UPDATE genres
SET name = $2
WHERE id = $1
RETURNING *;