package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	task_send_verify_email = "task::send_verify_email"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

// RedisTaskDistributor implements TaskDistributor interface
func (distributor *RedisTaskDistributor) DistributeTaskVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf(">>> fail to marshal payload: %w", err)
	}

	task := asynq.NewTask(task_send_verify_email, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf(">>> fail to enqueue task: %w", err)
	}
	log.Logger.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("Enqueued task")
	return nil
}

// RedisTaskProcessor implements TaskProcessor interface
func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf(">>> fail to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf(">>> the user does not exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf(">>> fail to get the user: %w", err)
	}

	// TODO: send email

	log.Logger.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("Processed task")
	return nil
}
