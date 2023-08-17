package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/kvnyijia/bank-app/db/sqlc"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func newRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{})
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	// Register task handlers
	mux.HandleFunc(task_send_verify_email, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
