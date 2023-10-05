package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	db "github.com/kvnyijia/bank-app/db/sqlc"
	"github.com/kvnyijia/bank-app/util"
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
		// Skip retry, comment out the following if would like to retry
		// if err == sql.ErrNoRows {
		// 	return fmt.Errorf(">>> the user does not exist: %w", asynq.SkipRetry)
		// }
		return fmt.Errorf(">>> fail to get the user: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf(">>> failed to create verify email: %w", err)
	}
	// TODO: send email
	subject := "Welcome to Bank App"
	verifyUrl := fmt.Sprintf("http://bank-app.org/verify_email?id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s, <br/>
		Thank you for registering with us! <br/>
		Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.FullName, verifyUrl)
	to := []string{user.Email}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf(">>> failed to send verift email: %w", err)
	}

	log.Logger.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("Processed task")
	return nil
}
