// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: subgenres.sql

package db

import (
	"context"
	"database/sql"
)

const createSubgenre = `-- name: CreateSubgenre :one
INSERT INTO subgenres (genres_id, name)
VALUES ($1, $2)
RETURNING id, genres_id, name, is_deleted, created_at
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
		&i.IsDeleted,
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
SELECT id, genres_id, name, is_deleted, created_at
FROM subgenres
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetSubgenre(ctx context.Context, id int64) (Subgenre, error) {
	row := q.db.QueryRowContext(ctx, getSubgenre, id)
	var i Subgenre
	err := row.Scan(
		&i.ID,
		&i.GenresID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}

const listAllSubgenres = `-- name: ListAllSubgenres :many
SELECT id, genres_id, name, is_deleted, created_at
FROM subgenres
ORDER BY id
`

func (q *Queries) ListAllSubgenres(ctx context.Context) ([]Subgenre, error) {
	rows, err := q.db.QueryContext(ctx, listAllSubgenres)
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

const listSubgenres = `-- name: ListSubgenres :many
SELECT id, genres_id, name, is_deleted, created_at
FROM subgenres
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

const listSubgenresNoticeable = `-- name: ListSubgenresNoticeable :many
SELECT subgenres.id,
  subgenres.name,
  subgenres.genres_id
FROM subgenres
  INNER JOIN books_subgenres ON subgenres.id = books_subgenres.subgenres_id
GROUP BY subgenres.id
LIMIT 6
`

type ListSubgenresNoticeableRow struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	GenresID int64  `json:"genres_id"`
}

func (q *Queries) ListSubgenresNoticeable(ctx context.Context) ([]ListSubgenresNoticeableRow, error) {
	rows, err := q.db.QueryContext(ctx, listSubgenresNoticeable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListSubgenresNoticeableRow{}
	for rows.Next() {
		var i ListSubgenresNoticeableRow
		if err := rows.Scan(&i.ID, &i.Name, &i.GenresID); err != nil {
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

const softDeleteSubgenre = `-- name: SoftDeleteSubgenre :one
UPDATE subgenres
SET is_deleted = true
WHERE id = $1
RETURNING id, genres_id, name, is_deleted, created_at
`

func (q *Queries) SoftDeleteSubgenre(ctx context.Context, id int64) (Subgenre, error) {
	row := q.db.QueryRowContext(ctx, softDeleteSubgenre, id)
	var i Subgenre
	err := row.Scan(
		&i.ID,
		&i.GenresID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}

const updateSubgenre = `-- name: UpdateSubgenre :one
UPDATE subgenres
SET genres_id = COALESCE($1, genres_id),
  name = COALESCE($2, name)
WHERE id = $3
RETURNING id, genres_id, name, is_deleted, created_at
`

type UpdateSubgenreParams struct {
	GenresID sql.NullInt64  `json:"genres_id"`
	Name     sql.NullString `json:"name"`
	ID       int64          `json:"id"`
}

func (q *Queries) UpdateSubgenre(ctx context.Context, arg UpdateSubgenreParams) (Subgenre, error) {
	row := q.db.QueryRowContext(ctx, updateSubgenre, arg.GenresID, arg.Name, arg.ID)
	var i Subgenre
	err := row.Scan(
		&i.ID,
		&i.GenresID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}