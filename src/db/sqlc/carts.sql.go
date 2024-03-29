// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: carts.sql

package db

import (
	"context"
)

const createCart = `-- name: CreateCart :one
INSERT INTO carts (books_id, username, amount, total)
VALUES ($1, $2, $3, $4)
RETURNING id, books_id, username, created_at, amount, total
`

type CreateCartParams struct {
	BooksID  int64   `json:"books_id"`
	Username string  `json:"username"`
	Amount   int32   `json:"amount"`
	Total    float64 `json:"total"`
}

func (q *Queries) CreateCart(ctx context.Context, arg CreateCartParams) (Cart, error) {
	row := q.db.QueryRowContext(ctx, createCart,
		arg.BooksID,
		arg.Username,
		arg.Amount,
		arg.Total,
	)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.BooksID,
		&i.Username,
		&i.CreatedAt,
		&i.Amount,
		&i.Total,
	)
	return i, err
}

const deleteCart = `-- name: DeleteCart :exec
DELETE FROM carts
WHERE id = $1
  AND username = $2
`

type DeleteCartParams struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (q *Queries) DeleteCart(ctx context.Context, arg DeleteCartParams) error {
	_, err := q.db.ExecContext(ctx, deleteCart, arg.ID, arg.Username)
	return err
}

const getCart = `-- name: GetCart :one
SELECT id, books_id, username, created_at, amount, total
FROM carts
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetCart(ctx context.Context, id int64) (Cart, error) {
	row := q.db.QueryRowContext(ctx, getCart, id)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.BooksID,
		&i.Username,
		&i.CreatedAt,
		&i.Amount,
		&i.Total,
	)
	return i, err
}

const listCartsByUsername = `-- name: ListCartsByUsername :many
SELECT id, books_id, username, created_at, amount, total
FROM carts
WHERE username = $1
ORDER BY id
`

func (q *Queries) ListCartsByUsername(ctx context.Context, username string) ([]Cart, error) {
	rows, err := q.db.QueryContext(ctx, listCartsByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Cart{}
	for rows.Next() {
		var i Cart
		if err := rows.Scan(
			&i.ID,
			&i.BooksID,
			&i.Username,
			&i.CreatedAt,
			&i.Amount,
			&i.Total,
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

const updateAmount = `-- name: UpdateAmount :one
UPDATE carts
SET amount = $2,
  total = $3
WHERE id = $1
RETURNING id, books_id, username, created_at, amount, total
`

type UpdateAmountParams struct {
	ID     int64   `json:"id"`
	Amount int32   `json:"amount"`
	Total  float64 `json:"total"`
}

func (q *Queries) UpdateAmount(ctx context.Context, arg UpdateAmountParams) (Cart, error) {
	row := q.db.QueryRowContext(ctx, updateAmount, arg.ID, arg.Amount, arg.Total)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.BooksID,
		&i.Username,
		&i.CreatedAt,
		&i.Amount,
		&i.Total,
	)
	return i, err
}
