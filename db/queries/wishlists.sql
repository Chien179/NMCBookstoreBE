-- name: GetWishlist :one
SELECT * FROM wishlists
WHERE id = $1 LIMIT 1;

-- name: CreateWishlist :one
INSERT INTO wishlists (
    users_id
) VALUES (
  $1
)
RETURNING *;