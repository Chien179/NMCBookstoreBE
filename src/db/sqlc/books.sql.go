// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: books.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
)

const createBook = `-- name: CreateBook :one
INSERT INTO books (
    name,
    price,
    sale,
    image,
    description,
    author,
    publisher,
    quantity
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, name, price, image, description, author, publisher, sale, quantity, is_deleted, created_at, rating
`

type CreateBookParams struct {
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Sale        float64  `json:"sale"`
	Image       []string `json:"image"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Publisher   string   `json:"publisher"`
	Quantity    int32    `json:"quantity"`
}

func (q *Queries) CreateBook(ctx context.Context, arg CreateBookParams) (Book, error) {
	row := q.db.QueryRowContext(ctx, createBook,
		arg.Name,
		arg.Price,
		arg.Sale,
		pq.Array(arg.Image),
		arg.Description,
		arg.Author,
		arg.Publisher,
		arg.Quantity,
	)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Price,
		pq.Array(&i.Image),
		&i.Description,
		&i.Author,
		&i.Publisher,
		&i.Sale,
		&i.Quantity,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.Rating,
	)
	return i, err
}

const deleteBook = `-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1
`

func (q *Queries) DeleteBook(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteBook, id)
	return err
}

const getBestBookByUser = `-- name: GetBestBookByUser :one
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
limit 1
`

type GetBestBookByUserRow struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Rating int32  `json:"rating"`
}

func (q *Queries) GetBestBookByUser(ctx context.Context, username string) (GetBestBookByUserRow, error) {
	row := q.db.QueryRowContext(ctx, getBestBookByUser, username)
	var i GetBestBookByUserRow
	err := row.Scan(&i.ID, &i.Name, &i.Rating)
	return i, err
}

const getBook = `-- name: GetBook :one
SELECT id, name, price, image, description, author, publisher, sale, quantity, is_deleted, created_at, rating
FROM books
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetBook(ctx context.Context, id int64) (Book, error) {
	row := q.db.QueryRowContext(ctx, getBook, id)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Price,
		pq.Array(&i.Image),
		&i.Description,
		&i.Author,
		&i.Publisher,
		&i.Sale,
		&i.Quantity,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.Rating,
	)
	return i, err
}

const listAllBooks = `-- name: ListAllBooks :many
SELECT id, name, price, image, description, author, publisher, sale, quantity, is_deleted, created_at, rating
FROM books
ORDER BY id
`

func (q *Queries) ListAllBooks(ctx context.Context) ([]Book, error) {
	rows, err := q.db.QueryContext(ctx, listAllBooks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Book{}
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Price,
			pq.Array(&i.Image),
			&i.Description,
			&i.Author,
			&i.Publisher,
			&i.Sale,
			&i.Quantity,
			&i.IsDeleted,
			&i.CreatedAt,
			&i.Rating,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBookFollowGenre = `-- name: ListBookFollowGenre :many
SELECT DISTINCT b.id,
  b."name",
  b.description,
  b.price,
  b.sale,
  b.sale,
  b.rating
FROM books AS b
  INNER JOIN books_genres AS bg ON b.id = bg.books_id
  INNER JOIN genres AS g ON bg.genres_id = g.id
WHERE g.id = $1
ORDER BY b.id,
  b."name",
  b.description,
  b.price,
  b.sale,
  b.rating
LIMIT $2
`

type ListBookFollowGenreParams struct {
	ID    int64 `json:"id"`
	Limit int32 `json:"limit"`
}

type ListBookFollowGenreRow struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Sale        float64 `json:"sale"`
	Sale_2      float64 `json:"sale_2"`
	Rating      float64 `json:"rating"`
}

func (q *Queries) ListBookFollowGenre(ctx context.Context, arg ListBookFollowGenreParams) ([]ListBookFollowGenreRow, error) {
	rows, err := q.db.QueryContext(ctx, listBookFollowGenre, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListBookFollowGenreRow{}
	for rows.Next() {
		var i ListBookFollowGenreRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Sale,
			&i.Sale_2,
			&i.Rating,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBooks = `-- name: ListBooks :one
SELECT t.total_page,
  JSON_AGG(
    json_build_object (
      'id',
      t.id,
      'name',
      t.name,
      'price',
      t.price,
      'sale',
      t.sale,
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
      'is_deleted',
      t.is_deleted,
      'created_at',
      t.created_at
    )
  ) AS books
FROM (
    SELECT CEILING(
        CAST(COUNT(id) OVER () AS FLOAT) / $1
      ) AS total_page,
      id, name, price, image, description, author, publisher, sale, quantity, is_deleted, created_at, rating
    FROM books
    ORDER BY id
    LIMIT $1 OFFSET $2
  ) AS t
GROUP BY t.total_page
`

type ListBooksParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListBooksRow struct {
	TotalPage float64         `json:"total_page"`
	Books     json.RawMessage `json:"books"`
}

func (q *Queries) ListBooks(ctx context.Context, arg ListBooksParams) (ListBooksRow, error) {
	row := q.db.QueryRowContext(ctx, listBooks, arg.Limit, arg.Offset)
	var i ListBooksRow
	err := row.Scan(&i.TotalPage, &i.Books)
	return i, err
}

const listNewestBooks = `-- name: ListNewestBooks :many
SELECT id, name, price, image, description, author, publisher, sale, quantity, is_deleted, created_at, rating
FROM books
ORDER BY created_at DESC
LIMIT 20
`

func (q *Queries) ListNewestBooks(ctx context.Context) ([]Book, error) {
	rows, err := q.db.QueryContext(ctx, listNewestBooks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Book{}
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Price,
			pq.Array(&i.Image),
			&i.Description,
			&i.Author,
			&i.Publisher,
			&i.Sale,
			&i.Quantity,
			&i.IsDeleted,
			&i.CreatedAt,
			&i.Rating,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTheBestBooks = `-- name: ListTheBestBooks :many
SELECT id, name, price, image, description, author, publisher, sale, quantity, is_deleted, created_at, rating
FROM books
ORDER BY rating DESC
LIMIT 20
`

func (q *Queries) ListTheBestBooks(ctx context.Context) ([]Book, error) {
	rows, err := q.db.QueryContext(ctx, listTheBestBooks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Book{}
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Price,
			pq.Array(&i.Image),
			&i.Description,
			&i.Author,
			&i.Publisher,
			&i.Sale,
			&i.Quantity,
			&i.IsDeleted,
			&i.CreatedAt,
			&i.Rating,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const softDeleteBook = `-- name: SoftDeleteBook :one
UPDATE books
SET is_deleted = true
WHERE id = $1
RETURNING id, name, price, image, description, author, publisher, sale, quantity, is_deleted, created_at, rating
`

func (q *Queries) SoftDeleteBook(ctx context.Context, id int64) (Book, error) {
	row := q.db.QueryRowContext(ctx, softDeleteBook, id)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Price,
		pq.Array(&i.Image),
		&i.Description,
		&i.Author,
		&i.Publisher,
		&i.Sale,
		&i.Quantity,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.Rating,
	)
	return i, err
}

const updateBook = `-- name: UpdateBook :one
UPDATE books
SET name = COALESCE($1, name),
  price = COALESCE($2, price),
  sale = COALESCE($3, sale),
  image = COALESCE($4, image),
  description = COALESCE($5, description),
  author = COALESCE($6, author),
  publisher = COALESCE($7, publisher),
  quantity = COALESCE($8, quantity)
WHERE id = $9
RETURNING id, name, price, image, description, author, publisher, sale, quantity, is_deleted, created_at, rating
`

type UpdateBookParams struct {
	Name        sql.NullString  `json:"name"`
	Price       sql.NullFloat64 `json:"price"`
	Sale        sql.NullFloat64 `json:"sale"`
	Image       []string        `json:"image"`
	Description sql.NullString  `json:"description"`
	Author      sql.NullString  `json:"author"`
	Publisher   sql.NullString  `json:"publisher"`
	Quantity    sql.NullInt32   `json:"quantity"`
	ID          int64           `json:"id"`
}

func (q *Queries) UpdateBook(ctx context.Context, arg UpdateBookParams) (Book, error) {
	row := q.db.QueryRowContext(ctx, updateBook,
		arg.Name,
		arg.Price,
		arg.Sale,
		pq.Array(arg.Image),
		arg.Description,
		arg.Author,
		arg.Publisher,
		arg.Quantity,
		arg.ID,
	)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Price,
		pq.Array(&i.Image),
		&i.Description,
		&i.Author,
		&i.Publisher,
		&i.Sale,
		&i.Quantity,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.Rating,
	)
	return i, err
}