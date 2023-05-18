-- name: GetReview :one
SELECT * FROM reviews
WHERE id = $1 LIMIT 1;

-- name: ListReviewsByBookID :one
SELECT t.total_page, JSON_AGG(json_build_object
    ('id',id,
    'username',username,
    'books_id',books_id,
    'comments',comments,
    'rating',rating,
    'created_at',created_at)
    ) AS reviews
	FROM (
      SELECT 
        CEILING(CAST(COUNT(id) OVER () AS FLOAT)/sqlc.arg('limit')) AS total_page, * 
      FROM reviews
      WHERE reviews.books_id = $1
      ORDER BY id
      LIMIT sqlc.arg('limit')
      OFFSET sqlc.arg('offset')
    ) AS t
    GROUP BY t.total_page;

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