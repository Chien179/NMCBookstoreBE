-- name: GetAddress :one
SELECT * FROM address
WHERE id = $1 LIMIT 1;

-- name: ListAddresses :one
(SELECT t.total_page, JSON_AGG(json_build_object
    ('id',id,
    'address',address,
    'username',username,
    'created_at',created_at)
    ) AS addresses
	FROM (
      SELECT 
        CEILING(CAST(COUNT(id) OVER () AS FLOAT)/sqlc.arg('limit')) AS total_page, * 
      FROM address
      WHERE address.username = $1
      ORDER BY id
      LIMIT sqlc.arg('limit')
      OFFSET sqlc.arg('offset')
    ) AS t
    GROUP BY t.total_page);

-- name: CreateAddress :one
INSERT INTO address (
  username,
  address,
  district_id,
  city_id
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
  district_id = COALESCE(sqlc.narg(district_id), district_id),
  city_id = COALESCE(sqlc.narg(city_id), city_id)
WHERE id = sqlc.arg(id)
RETURNING *;