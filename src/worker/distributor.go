package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error

	DistributeTaskSendResetPassword(
		ctx context.Context,
		payload *PayloadSendResetPassword,
		opts ...asynq.Option,
	) error

	DistributeTaskSendReportReview(
		ctx context.Context,
		payload *PayloadSendReportReview,
		opts ...asynq.Option,
	) error

	DistributeTaskSendOrderSuccess(
		ctx context.Context,
		payload *PayloadSendOrderSuccess,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
