-- name: GetBookWishlist :one
SELECT * FROM books_wishlists
WHERE id = $1 LIMIT 1;

-- name: ListBooksWishlists :many
SELECT * FROM books_wishlists
WHERE wishlists_id = $1
ORDER BY id;

-- name: CreateBookWishlist :one
INSERT INTO books_wishlists (
  books_id,
  wishlists_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteBookWishlist :exec
DELETE FROM books_wishlists
WHERE id = $1;