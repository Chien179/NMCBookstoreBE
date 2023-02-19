-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE role = &1
ORDER BY id
LIMIT $2
OFFSET $3;;

-- name: CreateUser :one
INSERT INTO users (
  username,
  full_name,
  email,
  password,
  image,
  phone_number,
  role,
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET full_name = $2,
    email = $3,
    password = $4,
    image = $5,
    phone_number = $6,
WHERE id = $1
RETURNING *;