-- name: GetAddress :one
SELECT *
FROM address
WHERE id = $1
LIMIT 1;
-- name: ListAddresses :many
SELECT address.id AS id,
  address.address AS address,
  districts.name AS district,
  cities.name AS city
FROM address
  INNER JOIN cities ON cities.id = address.city_id
  INNER JOIN districts ON districts.id = address.district_id
WHERE address.username = $1
ORDER BY address.id;
-- name: CreateAddress :one
INSERT INTO address (
    username,
    address,
    district_id,
    city_id
  )
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: DeleteAddress :exec
DELETE FROM address
WHERE id = $1;
-- name: UpdateAddress :one
UPDATE address
SET address = COALESCE(sqlc.narg(address), address),
  district_id = COALESCE(sqlc.narg(district_id), district_id),
  city_id = COALESCE(sqlc.narg(city_id), city_id)
WHERE id = sqlc.arg(id)
RETURNING *;