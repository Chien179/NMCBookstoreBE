// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: genres.sql

package db

import (
	"context"
)

const createGenre = `-- name: CreateGenre :one
INSERT INTO genres (name)
VALUES ($1)
RETURNING id, name, is_deleted, created_at
`

func (q *Queries) CreateGenre(ctx context.Context, name string) (Genre, error) {
	row := q.db.QueryRowContext(ctx, createGenre, name)
	var i Genre
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}

const deleteGenre = `-- name: DeleteGenre :exec
DELETE FROM genres
WHERE id = $1
`

func (q *Queries) DeleteGenre(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteGenre, id)
	return err
}

const getGenre = `-- name: GetGenre :one
SELECT id, name, is_deleted, created_at
FROM genres
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetGenre(ctx context.Context, id int64) (Genre, error) {
	row := q.db.QueryRowContext(ctx, getGenre, id)
	var i Genre
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}

const listGenres = `-- name: ListGenres :many
SELECT id, name, is_deleted, created_at
FROM genres
ORDER BY id
`

func (q *Queries) ListGenres(ctx context.Context) ([]Genre, error) {
	rows, err := q.db.QueryContext(ctx, listGenres)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Genre{}
	for rows.Next() {
		var i Genre
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IsDeleted,
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

const softDeleteGenre = `-- name: SoftDeleteGenre :one
UPDATE genres
SET is_deleted = true
WHERE id = $1
RETURNING id, name, is_deleted, created_at
`

func (q *Queries) SoftDeleteGenre(ctx context.Context, id int64) (Genre, error) {
	row := q.db.QueryRowContext(ctx, softDeleteGenre, id)
	var i Genre
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}

const updateGenre = `-- name: UpdateGenre :one
UPDATE genres
SET name = $2
WHERE id = $1
RETURNING id, name, is_deleted, created_at
`

type UpdateGenreParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateGenre(ctx context.Context, arg UpdateGenreParams) (Genre, error) {
	row := q.db.QueryRowContext(ctx, updateGenre, arg.ID, arg.Name)
	var i Genre
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}