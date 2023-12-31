package db

import "context"

type ReportReviewTxParams struct {
	UpdateReviewParams
	AfterCreate func(id int64) error
}

func (store *SQLStore) ReportReviewTx(ctx context.Context, arg ReportReviewTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.UpdateReview(ctx, arg.UpdateReviewParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(arg.ID)
	})

	return err
}
