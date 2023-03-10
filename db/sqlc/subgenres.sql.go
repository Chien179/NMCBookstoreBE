// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: subgenres.sql

package db

import (
	"context"
)

const createSubgenre = `-- name: CreateSubgenre :one
INSERT INTO subgenres (
    genres_id,
    name
) VALUES (
  $1, $2
)
RETURNING id, genres_id, name, created_at
`

type CreateSubgenreParams struct {
	GenresID int64  `json:"genres_id"`
	Name     string `json:"name"`
}

func (q *Queries) CreateSubgenre(ctx context.Context, arg CreateSubgenreParams) (Subgenre, error) {
	row := q.db.QueryRowContext(ctx, createSubgenre, arg.GenresID, arg.Name)
	var i Subgenre
	err := row.Scan(
		&i.ID,
		&i.GenresID,
		&i.Name,
		&i.CreatedAt,
	)
	return i, err
}

const deleteSubgenre = `-- name: DeleteSubgenre :exec
DELETE FROM subgenres
WHERE id = $1
`

func (q *Queries) DeleteSubgenre(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteSubgenre, id)
	return err
}

const getSubgenre = `-- name: GetSubgenre :one
SELECT id, genres_id, name, created_at FROM subgenres
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSubgenre(ctx context.Context, id int64) (Subgenre, error) {
	row := q.db.QueryRowContext(ctx, getSubgenre, id)
	var i Subgenre
	err := row.Scan(
		&i.ID,
		&i.GenresID,
		&i.Name,
		&i.CreatedAt,
	)
	return i, err
}

const listSubgenres = `-- name: ListSubgenres :many
SELECT id, genres_id, name, created_at FROM subgenres
WHERE genres_id = $1
ORDER BY id
`

func (q *Queries) ListSubgenres(ctx context.Context, genresID int64) ([]Subgenre, error) {
	rows, err := q.db.QueryContext(ctx, listSubgenres, genresID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Subgenre{}
	for rows.Next() {
		var i Subgenre
		if err := rows.Scan(
			&i.ID,
			&i.GenresID,
			&i.Name,
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

const updateSubgenre = `-- name: UpdateSubgenre :one
UPDATE subgenres
SET 
  genres_id = $2,
  name = $3
WHERE id = $1
RETURNING id, genres_id, name, created_at
`

type UpdateSubgenreParams struct {
	ID       int64  `json:"id"`
	GenresID int64  `json:"genres_id"`
	Name     string `json:"name"`
}

func (q *Queries) UpdateSubgenre(ctx context.Context, arg UpdateSubgenreParams) (Subgenre, error) {
	row := q.db.QueryRowContext(ctx, updateSubgenre, arg.ID, arg.GenresID, arg.Name)
	var i Subgenre
	err := row.Scan(
		&i.ID,
		&i.GenresID,
		&i.Name,
		&i.CreatedAt,
	)
	return i, err
}
