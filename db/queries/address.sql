-- name: GetAddress :one
SELECT * FROM address
WHERE id = $1 LIMIT 1;

-- name: ListAddresses :one
SELECT
  (SELECT (COUNT(*)/sqlc.arg('limit'))
     FROM address
     WHERE address.username = $1) 
     as total_page, 
  (SELECT JSON_AGG(t.*) FROM (
    SELECT * FROM address
    WHERE address.username = $1
    ORDER BY id
    LIMIT sqlc.arg('limit')
    OFFSET sqlc.arg('offset')
  ) AS t) AS address;

-- name: CreateAddress :one
INSERT INTO address (
  username,
  address,
  district,
  city
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteAddress :exec
DELETE FROM address
WHERE id = $1;

-- name: UpdateAddress :one
UPDATE address
SET  address = COALESCE(sqlc.narg(address), address),
  district = COALESCE(sqlc.narg(district), district),
  city = COALESCE(sqlc.narg(city), city)
WHERE id = sqlc.arg(id)
RETURNING *;