-- name: GetReview :one
SELECT * FROM reviews
WHERE id = $1 LIMIT 1;

-- name: GetReviewsByBookID :many
SELECT * FROM reviews
WHERE books_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: CreateReview :one
INSERT INTO reviews (
    users_id,
    books_id,
    comments,
    rating
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;