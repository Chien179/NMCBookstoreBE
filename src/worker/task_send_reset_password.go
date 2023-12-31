package worker

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendResetPassword = "task:send_reset_password"

type PayloadSendResetPassword struct {
	Email string `json:"email"`
}

func (distributior *RedisTaskDistributor) DistributeTaskSendResetPassword(
	ctx context.Context,
	payload *PayloadSendResetPassword,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendResetPassword, jsonPayload, opts...)
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

func (processor *RedisTaskProcessor) ProcessTaskSendResetPassword(ctx context.Context, task *asynq.Task) error {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	var payload PayloadSendResetPassword
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err, asynq.SkipRetry)
	}

	user, err := processor.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	resetPassword, err := processor.store.CreateResetPassword(ctx, db.CreateResetPasswordParams{
		Username:  user.Username,
		ResetCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Reset password for NMC Bookstore account"
	resetPasswordUrl := fmt.Sprintf(config.CLIENT_HOST+"/reset_password?id=%d&reset_code=%s",
		resetPassword.ID, resetPassword.ResetCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	We received a request to reset your password!<br/>
	Please use the link below to set up a new password for your account. If you did not request to reset your password, ignore this email.</br>
	<a href="%s">SET NEW PASSWORD</a>
	`, user.FullName, resetPasswordUrl)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")
	return nil
}
