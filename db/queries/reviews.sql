-- name: GetReview :one
SELECT * FROM reviews
WHERE id = $1 LIMIT 1;

-- name: ListReviewsByBookID :many
SELECT * FROM reviews
WHERE books_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: CreateReview :one
INSERT INTO reviews (
    username,
    books_id,
    comments,
    rating
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;