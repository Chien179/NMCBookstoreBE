-- name: GetCity :one
SELECT * FROM cities
WHERE id = $1 LIMIT 1;

-- name: ListCities :many
SELECT * FROM cities
ORDER BY id;