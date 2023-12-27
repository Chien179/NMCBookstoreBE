// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: reset_password.sql

package db

import (
	"context"
)

const createResetPassword = `-- name: CreateResetPassword :one
INSERT INTO reset_passwords (username, reset_code)
VALUES ($1, $2)
RETURNING id, username, reset_code, is_used, created_at, expired_at
`

type CreateResetPasswordParams struct {
	Username  string `json:"username"`
	ResetCode string `json:"reset_code"`
}

func (q *Queries) CreateResetPassword(ctx context.Context, arg CreateResetPasswordParams) (ResetPassword, error) {
	row := q.db.QueryRowContext(ctx, createResetPassword, arg.Username, arg.ResetCode)
	var i ResetPassword
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.ResetCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

const updateResetPassword = `-- name: UpdateResetPassword :one
UPDATE reset_passwords
SET is_used = TRUE
WHERE id = $1
    AND reset_code = $2
    AND is_used = FALSE
    AND expired_at > now()
RETURNING id, username, reset_code, is_used, created_at, expired_at
`

type UpdateResetPasswordParams struct {
	ID        int64  `json:"id"`
	ResetCode string `json:"reset_code"`
}

func (q *Queries) UpdateResetPassword(ctx context.Context, arg UpdateResetPasswordParams) (ResetPassword, error) {
	row := q.db.QueryRowContext(ctx, updateResetPassword, arg.ID, arg.ResetCode)
	var i ResetPassword
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.ResetCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}