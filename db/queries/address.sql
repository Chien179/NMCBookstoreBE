-- name: GetAddress :one
SELECT * FROM address
WHERE id = $1 LIMIT 1;

-- name: ListAddresses :many
SELECT * FROM address
WHERE username = $1
ORDER BY id
LIMIT $2
OFFSET $3;

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
SET  address = $2,
  district = $3,
  city = $4
WHERE id = $1
RETURNING *;