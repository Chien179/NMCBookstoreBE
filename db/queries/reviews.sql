-- name: GetReview :one
SELECT * FROM reviews
WHERE id = $1 LIMIT 1;

-- name: ListReviewsByBookID :one
SELECT
  (SELECT (COUNT(*)/sqlc.arg('limit'))
     FROM reviews
     WHERE reviews.books_id = $1) 
     as total_page, 
    (SELECT JSON_AGG(t.*) FROM (
      SELECT * FROM reviews
      WHERE reviews.books_id = $1
      ORDER BY id
      LIMIT sqlc.arg('limit')
      OFFSET sqlc.arg('offset')
    )AS t) AS reviews;

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