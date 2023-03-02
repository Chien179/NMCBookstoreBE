// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: books_wishlists.sql

package db

import (
	"context"
)

const createBookWishlist = `-- name: CreateBookWishlist :one
INSERT INTO books_wishlists (
  books_id,
  wishlists_id
) VALUES (
  $1, $2
)
RETURNING id, books_id, wishlists_id, created_at
`

type CreateBookWishlistParams struct {
	BooksID     int64 `json:"books_id"`
	WishlistsID int64 `json:"wishlists_id"`
}

func (q *Queries) CreateBookWishlist(ctx context.Context, arg CreateBookWishlistParams) (BooksWishlist, error) {
	row := q.db.QueryRowContext(ctx, createBookWishlist, arg.BooksID, arg.WishlistsID)
	var i BooksWishlist
	err := row.Scan(
		&i.ID,
		&i.BooksID,
		&i.WishlistsID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteBookWishlist = `-- name: DeleteBookWishlist :exec
DELETE FROM books_wishlists
WHERE id = $1
`

func (q *Queries) DeleteBookWishlist(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteBookWishlist, id)
	return err
}

const getBookWishlist = `-- name: GetBookWishlist :one
SELECT id, books_id, wishlists_id, created_at FROM books_wishlists
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetBookWishlist(ctx context.Context, id int64) (BooksWishlist, error) {
	row := q.db.QueryRowContext(ctx, getBookWishlist, id)
	var i BooksWishlist
	err := row.Scan(
		&i.ID,
		&i.BooksID,
		&i.WishlistsID,
		&i.CreatedAt,
	)
	return i, err
}

const listBooksWishlists = `-- name: ListBooksWishlists :many
SELECT id, books_id, wishlists_id, created_at FROM books_wishlists
WHERE wishlists_id = $1
ORDER BY id
`

func (q *Queries) ListBooksWishlists(ctx context.Context, wishlistsID int64) ([]BooksWishlist, error) {
	rows, err := q.db.QueryContext(ctx, listBooksWishlists, wishlistsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []BooksWishlist{}
	for rows.Next() {
		var i BooksWishlist
		if err := rows.Scan(
			&i.ID,
			&i.BooksID,
			&i.WishlistsID,
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
