package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	db "github.com/machearn/galaxy_service/db/sqlc"
)

const (
	TaskSendEmail = "task:send_verification_email"
)

type EmailPayload struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerificationEmail(ctx context.Context, payload EmailPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal payload: %s", err.Error())
		return err
	}

	task := asynq.NewTask(TaskSendEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		log.Printf("failed to enqueue task: %s", err.Error())
		return err
	}

	log.Printf("enqueued task: type: %s, payload: %s, queue: %s, max_retry: %d", task.Type(), task.Payload(), info.Queue, info.MaxRetry)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerificationEmail(ctx context.Context, task *asynq.Task) error {
	var payload EmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Printf("failed to unmarshal payload: %s", err.Error())
		return err
	}

	user, err := processor.store.GetMemberByName(ctx, payload.Username)
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Printf("failed to get user: %s", pqErr.Error())
		return pqErr
	}

	email, err := processor.store.CreateVerificationEmail(ctx, db.CreateVerificationEmailParams{
		MemberID:   user.ID,
		Email:      user.Email,
		SecretCode: uuid.New().String(),
	})
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Printf("failed to create verification email: %s", pqErr.Error())
		return pqErr
	}

	to := []string{user.Email}
	subject := "Verify your email"
	body := fmt.Sprintf("Please Verify your email by clicking <a href=\"http://localhost:8080/verify?email_id=%d&code=%s\">here</a>", email.ID, email.SecretCode)

	if err := processor.sender.SendEmail(to, subject, body); err != nil {
		log.Printf("failed to send email: %s", err.Error())
		return err
	}

	log.Printf("processed task: type: %s, payload: %s", task.Type(), task.Payload())
	return nil
}
