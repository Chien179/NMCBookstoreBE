-- name: GetSubgenre :one
SELECT * FROM subgenres
WHERE id = $1 LIMIT 1;

-- name: ListSubgenres :many
SELECT * FROM subgenres
WHERE genres_id = $1
ORDER BY id;

-- name: CreateSubgenre :one
INSERT INTO subgenres (
    genres_id,
    name
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteSubgenre :exec
DELETE FROM subgenres
WHERE id = $1;

-- name: UpdateSubgenre :one
UPDATE subgenres
SET 
  genres_id = $2,
  name = $3
WHERE id = $1
RETURNING *;