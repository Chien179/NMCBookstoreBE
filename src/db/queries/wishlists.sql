-- name: GetWishlist :one
SELECT *
FROM wishlists
WHERE id = $1
LIMIT 1;
-- name: ListWishlistsByUsername :many
SELECT *
FROM wishlists
WHERE username = $1
ORDER BY id;
-- name: CreateWishlist :one
INSERT INTO wishlists (books_id, username)
VALUES ($1, $2)
RETURNING *;
-- name: DeleteWishlist :exec
DELETE FROM wishlists
WHERE id = $1
  AND username = $2;