-- name: GetDislike :one
SELECT *
FROM "dislike"
WHERE username = sqlc.arg(username)
    AND review_id = sqlc.arg(review_id)
LIMIT 1;
-- name: UpdateDislike :one
UPDATE "dislike"
SET is_dislike = sqlc.arg(is_dislike)
WHERE username = sqlc.arg(username)
    AND review_id = sqlc.arg(review_id)
RETURNING *;
-- name: CreatedDislike :one
INSERT INTO "dislike" (username, review_id, is_dislike)
VALUES ($1, $2, $3)
RETURNING *;