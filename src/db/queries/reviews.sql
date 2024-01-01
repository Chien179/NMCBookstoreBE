-- name: GetReview :one
SELECT *
FROM reviews
WHERE id = $1
LIMIT 1;
-- name: ListReviewsByBookID :one
SELECT t.total_page,
  JSON_AGG(
    json_build_object (
      'id',
      id,
      'username',
      username,
      'image',
      image,
      'books_id',
      books_id,
      'comments',
      comments,
      'rating',
      rating,
      'is_deleted',
      is_deleted,
      'created_at',
      created_at
    )
  ) AS reviews
FROM (
    SELECT reviews.id,
      CEILING(
        CAST(COUNT(id) OVER () AS FLOAT) / sqlc.arg('limit')
      ) AS total_page,
      users.username AS username,
      users.image AS image,
      reviews.books_id AS books_id,
      reviews."comments" AS "comments",
      reviews.rating AS rating,
      reviews.is_deleted AS is_deleted,
      reviews.created_at AS created_at
    FROM reviews
      INNER JOIN users ON reviews.username = users.username
    WHERE reviews.books_id = $1
    ORDER BY id
    LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset')
  ) AS t
GROUP BY t.total_page;
-- name: CreateReview :one
INSERT INTO reviews (
    username,
    books_id,
    comments,
    rating
  )
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: UpdateReview :one
UPDATE reviews
SET liked = COALESCE(sqlc.narg(liked), liked),
  disliked = COALESCE(sqlc.narg(disliked), disliked),
  reported = COALESCE(sqlc.narg(reported), reported)
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;
-- name: SoftDeleteReview :one
UPDATE reviews
SET is_deleted = true
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: GetCountReviewByUser :one
SELECT COUNT(*) AS reviews
FROM reviews AS r
WHERE r.username = $1;
-- name: ListReviews :many
SELECT *
FROM reviews
ORDER BY id;