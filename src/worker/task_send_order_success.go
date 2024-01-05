package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendOrderSuccess = "task:send_order_success"

type PayloadSendOrderSuccess struct {
	Order models.OrderReponse `json:"order"`
}

func (distributior *RedisTaskDistributor) DistributeTaskSendOrderSuccess(
	ctx context.Context,
	payload *PayloadSendOrderSuccess,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendOrderSuccess, jsonPayload, opts...)
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

func (processor *RedisTaskProcessor) ProcessTaskSendOrderSuccess(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendOrderSuccess
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err, asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Order.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	subject := "Order Confirmation"
	content := fmt.Sprintf(`Hello, %s<br/>
	Thank you for your order<br/>
	</br>
	Here's the details information:<br/>
	</br>
	Order confirmation: %d<br/>
	Total: %.2f$</br>
	Delivery address: %s<br/>`, user.Username, payload.Order.ID, payload.Order.SubTotal, payload.Order.ToAddress)
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
