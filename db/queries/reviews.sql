-- name: GetReviewsByBookID :many
SELECT * FROM reviews
WHERE books_id = $1
ORDER BY id;

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