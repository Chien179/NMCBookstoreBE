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
SET full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email),
    image = COALESCE(sqlc.narg(image), image),
    phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
    password = COALESCE(sqlc.narg(password), password),
    is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified)
WHERE 
  username = sqlc.arg(username)
RETURNING *;