// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: like.sql

package db

import (
	"context"
)

const createLike = `-- name: CreateLike :one
INSERT INTO "like" (username, review_id, is_like)
VALUES ($1, $2, $3)
RETURNING id, username, review_id, is_like
`

type CreateLikeParams struct {
	Username string `json:"username"`
	ReviewID int64  `json:"review_id"`
	IsLike   bool   `json:"is_like"`
}

func (q *Queries) CreateLike(ctx context.Context, arg CreateLikeParams) (Like, error) {
	row := q.db.QueryRowContext(ctx, createLike, arg.Username, arg.ReviewID, arg.IsLike)
	var i Like
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.ReviewID,
		&i.IsLike,
	)
	return i, err
}

const getCountLikeByUser = `-- name: GetCountLikeByUser :one
SELECT COUNT(*) AS votes
FROM "like" as l
WHERE l.username = $1
`

func (q *Queries) GetCountLikeByUser(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getCountLikeByUser, username)
	var votes int64
	err := row.Scan(&votes)
	return votes, err
}

const getLike = `-- name: GetLike :one
SELECT id, username, review_id, is_like
FROM "like"
WHERE username = $1
    AND review_id = $2
LIMIT 1
`

type GetLikeParams struct {
	Username string `json:"username"`
	ReviewID int64  `json:"review_id"`
}

func (q *Queries) GetLike(ctx context.Context, arg GetLikeParams) (Like, error) {
	row := q.db.QueryRowContext(ctx, getLike, arg.Username, arg.ReviewID)
	var i Like
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.ReviewID,
		&i.IsLike,
	)
	return i, err
}

const listLike = `-- name: ListLike :many
SELECT id, username, review_id, is_like
FROM "like"
WHERE username = $1
ORDER BY review_id
`

func (q *Queries) ListLike(ctx context.Context, username string) ([]Like, error) {
	rows, err := q.db.QueryContext(ctx, listLike, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Like{}
	for rows.Next() {
		var i Like
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.ReviewID,
			&i.IsLike,
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

const updateLike = `-- name: UpdateLike :one
UPDATE "like"
SET is_like = $1
WHERE username = $2
    AND review_id = $3
RETURNING id, username, review_id, is_like
`

type UpdateLikeParams struct {
	IsLike   bool   `json:"is_like"`
	Username string `json:"username"`
	ReviewID int64  `json:"review_id"`
}

func (q *Queries) UpdateLike(ctx context.Context, arg UpdateLikeParams) (Like, error) {
	row := q.db.QueryRowContext(ctx, updateLike, arg.IsLike, arg.Username, arg.ReviewID)
	var i Like
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.ReviewID,
		&i.IsLike,
	)
	return i, err
}
