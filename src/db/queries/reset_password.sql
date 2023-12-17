-- name: CreateResetPassword :one
INSERT INTO reset_passwords (username, reset_code)
VALUES ($1, $2)
RETURNING *;
-- name: UpdateResetPassword :one
UPDATE reset_passwords
SET is_used = TRUE
WHERE id = @id
    AND reset_code = @reset_code
    AND is_used = FALSE
    AND expired_at > now()
RETURNING *;