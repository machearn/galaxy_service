package worker

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/mail"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerificationEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	sender mail.EmailSender
}

func NewRedisTaskProcessor(store db.Store, sender mail.EmailSender, opts asynq.RedisClientOpt) TaskProcessor {
	server := asynq.NewServer(opts, asynq.Config{
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Printf("failed to process task: type: %s, payload: %s, error: %s", task.Type(), task.Payload(), err.Error())
		}),
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		sender: sender,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendEmail, processor.ProcessTaskSendVerificationEmail)
	return processor.server.Start(mux)
}
