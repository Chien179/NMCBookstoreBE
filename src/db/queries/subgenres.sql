-- name: GetSubgenre :one
SELECT *
FROM subgenres
WHERE id = $1
LIMIT 1;
-- name: ListSubgenres :many
SELECT *
FROM subgenres
WHERE genres_id = $1
ORDER BY id;
-- name: ListAllSubgenres :many
SELECT *
FROM subgenres
ORDER BY id;
-- name: CreateSubgenre :one
INSERT INTO subgenres (genres_id, name)
VALUES ($1, $2)
RETURNING *;
-- name: DeleteSubgenre :exec
DELETE FROM subgenres
WHERE id = $1;
-- name: SoftDeleteSubgenre :one
UPDATE subgenres
SET is_deleted = true
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: UpdateSubgenre :one
UPDATE subgenres
SET genres_id = COALESCE(sqlc.narg(genres_id), genres_id),
  name = COALESCE(sqlc.narg(name), name)
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: ListSubgenresNoticeable :many
SELECT subgenres.id,
  subgenres.name,
  subgenres.genres_id
FROM subgenres
  INNER JOIN books_subgenres ON subgenres.id = books_subgenres.subgenres_id
GROUP BY subgenres.id
LIMIT 6;