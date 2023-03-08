-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
  username,
  full_name,
  email,
  password,
  image,
  phone_number,
  role
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;

-- name: UpdateUser :one
UPDATE users
SET full_name = $2,
    email = $3,
    image = $4,
    phone_number = $5
WHERE username = $1
RETURNING *;

-- name: UpdatePassword :one
UPDATE users
SET password = $2
WHERE username = $1
RETURNING *;