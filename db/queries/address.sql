-- name: GetAddres :one
SELECT * FROM address
WHERE id = $1 LIMIT 1;

-- name: ListAddress :many
SELECT * FROM address
ORDER BY id;

-- name: CreateAddres :one
INSERT INTO address (
  users_id,
  address,
  district,
  city
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteAddres :exec
DELETE FROM address
WHERE id = $1;

-- name: UpdateAddres :one
UPDATE address
SET  address = $2,
  district = $3,
  city = $4
WHERE id = $1
RETURNING *;