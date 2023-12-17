-- name: GetLike :one
SELECT *
FROM "like"
WHERE username = sqlc.arg(username)
    AND review_id = sqlc.arg(review_id)
LIMIT 1;
-- name: UpdateLike :one
UPDATE "like"
SET is_like = sqlc.arg(is_like)
WHERE username = sqlc.arg(username)
    AND review_id = sqlc.arg(review_id)
RETURNING *;
-- name: CreateLike :one
INSERT INTO "like" (username, review_id, is_like)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetCountLikeByUser :one
SELECT COUNT(*) AS votes
FROM "like" as l
WHERE l.username = $1;