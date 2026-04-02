package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const queryKey = "jobs:pending"

type JobMessage struct {
	JobID    int    `json:"job_id"`
	UserID   int64  `json:"user_id"`
	ChatID   int64  `json:"chat_id"`
	URL      string `json:"url"`
	Platform string `json:"platform"`
}

func (r *RedisClient) Enqueue(ctx context.Context, job JobMessage) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	return r.LPush(ctx, queryKey, data).Err()
}

func (r *RedisClient) Dequeue(ctx context.Context, timeout time.Duration) (*JobMessage, error) {
	result, err := r.BRPop(ctx, timeout, queryKey).Result()
	if err != nil {
		return nil, err
	}

	var job JobMessage
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}
