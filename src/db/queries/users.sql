-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;
-- name: ListUsers :many
SELECT *
FROM users
ORDER BY username;
-- name: CreateUser :one
INSERT INTO users (
    username,
    full_name,
    email,
    password,
    age,
    sex,
    image,
    phone_number
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;
-- name: SoftDeleteUser :one
UPDATE users
SET is_deleted = true
WHERE username = sqlc.arg(id)
RETURNING *;
-- name: UpdateUser :one
UPDATE users
SET full_name = COALESCE(sqlc.narg(full_name), full_name),
  email = COALESCE(sqlc.narg(email), email),
  image = COALESCE(sqlc.narg(image), image),
  phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
  age = COALESCE(sqlc.narg(age), age),
  sex = COALESCE(sqlc.narg(sex), sex),
  password = COALESCE(sqlc.narg(password), password),
  password_changed_at = COALESCE(
    sqlc.narg(password_changed_at),
    password_changed_at
  ),
  is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified)
WHERE username = sqlc.arg(username)
RETURNING *;