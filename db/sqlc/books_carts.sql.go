// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: books_carts.sql

package db

import (
	"context"
)

const createBookCart = `-- name: CreateBookCart :one
INSERT INTO books_carts (
  books_id,
  carts_id
) VALUES (
  $1, $2
)
RETURNING id, books_id, carts_id, created_at
`

type CreateBookCartParams struct {
	BooksID int64 `json:"books_id"`
	CartsID int64 `json:"carts_id"`
}

func (q *Queries) CreateBookCart(ctx context.Context, arg CreateBookCartParams) (BooksCart, error) {
	row := q.db.QueryRowContext(ctx, createBookCart, arg.BooksID, arg.CartsID)
	var i BooksCart
	err := row.Scan(
		&i.ID,
		&i.BooksID,
		&i.CartsID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteBookCart = `-- name: DeleteBookCart :exec
DELETE FROM books_carts
WHERE id = $1
`

func (q *Queries) DeleteBookCart(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteBookCart, id)
	return err
}

const getBookCart = `-- name: GetBookCart :one
SELECT id, books_id, carts_id, created_at FROM books_carts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetBookCart(ctx context.Context, id int64) (BooksCart, error) {
	row := q.db.QueryRowContext(ctx, getBookCart, id)
	var i BooksCart
	err := row.Scan(
		&i.ID,
		&i.BooksID,
		&i.CartsID,
		&i.CreatedAt,
	)
	return i, err
}

const listBooksCarts = `-- name: ListBooksCarts :many
SELECT id, books_id, carts_id, created_at FROM books_carts
WHERE carts_id = $1
ORDER BY id
`

func (q *Queries) ListBooksCarts(ctx context.Context, cartsID int64) ([]BooksCart, error) {
	rows, err := q.db.QueryContext(ctx, listBooksCarts, cartsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []BooksCart{}
	for rows.Next() {
		var i BooksCart
		if err := rows.Scan(
			&i.ID,
			&i.BooksID,
			&i.CartsID,
			&i.CreatedAt,
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
