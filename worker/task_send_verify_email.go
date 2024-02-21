package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/template/email_tmpl"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	// TaskSendVerifyEmail represents the name of the task that sends the email verification email.
	TaskSendVerifyEmail = "task:send_verify_email"
)

// PayloadSendVerifyEmail provides the userID.
type PayloadSendVerifyEmail struct {
	UserID    int64  `json:"user_id"`
	ClientIp  string `json:"client_ip"`
}

// DistributeTaskSendVerifyEmail enqueues the given task to be processed by a worker. It returns an error if the task could
// not be enqueued.
func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")

	return nil
}

// ProcessTaskSendVerifyEmail processes a 'TaskSendVerifyEmail' task.
func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.dbStore.GetUserByID(ctx, payload.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	token, tokenpayload, err := processor.tokenMaker.CreateToken(user.ID, user.AccountID, 25*time.Minute, nil)
	if err != nil {
		return err
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	verifyEmailSession, err := processor.dbStore.CreateSession(ctx, db.CreateSessionParams{
		ID:        uuid,
		UserID:    user.ID,
		Token:     util.Extract(token),
		Scope:     "verify_email",
		ClientIp:  payload.ClientIp,
		IsBlocked: false,
		ExpiresAt: tokenpayload.ExpiredAt,
	})
	if err != nil {
		werr := fmt.Errorf("failed to create verify email session: %s", err.Error())
		return werr
	}

	subject := "Welcome to StoreHub"
	content := email_tmpl.WelcomeTmpl(user.FirstName, fmt.Sprintf(
		"http://store-hub-frontend.vercel.app/auth/verify-email?email=%s&secret_code=%s",
		user.Email,
		verifyEmailSession.Token,
	))
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")

	return nil
}
