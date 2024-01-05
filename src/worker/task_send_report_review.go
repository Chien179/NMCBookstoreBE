package worker

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendReportReview = "task:send_report_review"

type PayloadSendReportReview struct {
	Review db.Review `json:"review"`
}

func (distributior *RedisTaskDistributor) DistributeTaskSendReportReview(
	ctx context.Context,
	payload *PayloadSendReportReview,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendReportReview, jsonPayload, opts...)
	info, err := distributior.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendReportReview(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendReportReview
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err, asynq.SkipRetry)
	}

	book, err := processor.store.GetBook(ctx, payload.Review.BooksID)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	user, err := processor.store.GetUser(ctx, payload.Review.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	subject := "Review report"
	content := fmt.Sprintf(`Hello %s,<br/>
	Your review "%s" on book "%s" have been reported</br>
	`, user.Username, payload.Review.Comments, book.Name)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send report review email %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")
	return nil
}
