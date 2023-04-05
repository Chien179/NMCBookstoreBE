package db

import (
	"context"
	"database/sql"
	"time"
)

type ResetPasswordTxParams struct {
	Id        int64
	ResetCode string
	Password  string
}

type ResetPasswordTxResult struct {
	User          User
	ResetPassword ResetPassword
}

func (store *SQLStore) ResetPasswordTx(ctx context.Context, arg ResetPasswordTxParams) (ResetPasswordTxResult, error) {
	var result ResetPasswordTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.ResetPassword, err = q.UpdateResetPassword(ctx, UpdateResetPasswordParams{
			ID:        arg.Id,
			ResetCode: arg.ResetCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.ResetPassword.Username,
			Password: sql.NullString{
				String: arg.Password,
				Valid:  true,
			},
			PasswordChangedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})
		return err
	})

	return result, err
}
