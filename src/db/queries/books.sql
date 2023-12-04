-- name: GetBook :one
SELECT *
FROM books
WHERE id = $1
LIMIT 1;
-- name: ListBooks :one
SELECT t.total_page,
  JSON_AGG(
    json_build_object (
      'id',
      t.id,
      'name',
      t.name,
      'price',
      t.price,
      'image',
      t.image,
      'description',
      t.description,
      'author',
      t.author,
      'publisher',
      t.publisher,
      'quantity',
      t.quantity,
      'rating',
      t.rating,
      'created_at',
      t.created_at
    )
  ) AS books
FROM (
    SELECT CEILING(
        CAST(COUNT(id) OVER () AS FLOAT) / sqlc.arg('limit')
      ) AS total_page,
      *
    FROM books
    ORDER BY id
    LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset')
  ) AS t
GROUP BY t.total_page;
-- name: ListAllBooks :many
SELECT *
FROM books
ORDER BY id;
-- name: ListTheBestBooks :many
SELECT *
FROM books
ORDER BY rating DESC
LIMIT 20;
-- name: ListNewestBooks :many
SELECT *
FROM books
ORDER BY created_at DESC
LIMIT 20;
-- name: CreateBook :one
INSERT INTO books (
    name,
    price,
    image,
    description,
    author,
    publisher,
    quantity
  )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;
-- name: SoftDeleteBook :one
UPDATE books
SET is_deleted = true
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: UpdateBook :one
UPDATE books
SET name = COALESCE(sqlc.narg(name), name),
  price = COALESCE(sqlc.narg(price), price),
  image = COALESCE(sqlc.narg(image), image),
  description = COALESCE(sqlc.narg(description), description),
  author = COALESCE(sqlc.narg(author), author),
  publisher = COALESCE(sqlc.narg(publisher), publisher),
  quantity = COALESCE(sqlc.narg(quantity), quantity)
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: GetBestBookByUser :one
select b.id,
  b.name,
  r.rating
from reviews as r
  inner join books as b on r.books_id = b.id
where r.username = $1
group by b.id,
  b.name,
  r.rating
order by r.rating desc
limit 1;