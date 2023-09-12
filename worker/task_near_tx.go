package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OCD-Labs/store-hub/near"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	// TaskNEARTx represents the name of the task for NEAR transactions.
	TaskNEARTx = "task:near_tx"
)

// PayloadNEARTx provides the NEAR transaction args.
type PayloadNEARTx struct {
	Args []string `json:"args"`
}

// DistributeTaskNEARTx enqueues the given near task to be processed by a worker.
func (distributor *RedisTaskDistributor) DistributeTaskNEARTx(
	ctx context.Context,
	args *PayloadNEARTx,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskNEARTx, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("task enqueued")

	return nil
}

// ProcessTaskNEARTx processes a 'TaskNEARTx' task.
func (processor *RedisTaskProcessor) ProcessTaskNEARTx(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadNEARTx
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	if err := near.RunNearCLICommand(payload.Args...); err != nil {
		return err
	}

	log.Info().Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Msg("processed task")

	return nil
}
