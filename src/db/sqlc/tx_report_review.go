package db

import "context"

type ReportReviewTxParams struct {
	ID          int64
	AfterCreate func(review Review) error
}

func (store *SQLStore) ReportReviewTx(ctx context.Context, arg ReportReviewTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result, err := q.SoftDeleteReview(ctx, arg.ID)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result)
	})

	return err
}
