-- name: GetDistrict :one
SELECT *
FROM districts
WHERE id = $1
LIMIT 1;
-- name: ListDistricts :many
SELECT *
FROM districts
WHERE city_id = $1
ORDER BY id;