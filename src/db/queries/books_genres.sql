-- name: GetBookGenre :one
SELECT *
FROM books_genres
WHERE id = $1
LIMIT 1;
-- name: ListBooksGenresByGenreID :many
SELECT *
FROM books_genres
WHERE genres_id = $1
ORDER BY id;
-- name: ListBooksGenresByBookID :many
SELECT *
FROM books_genres
WHERE books_id = $1
ORDER BY id;
-- name: ListBooksGenresIDByBookID :many
SELECT genres_id
FROM books_genres
WHERE books_id = $1
ORDER BY id;
-- name: CreateBookGenre :one
INSERT INTO books_genres (books_id, genres_id)
VALUES ($1, $2)
RETURNING *;
-- name: DeleteBookGenreByBooksID :exec
DELETE FROM books_genres
WHERE books_id = $1;