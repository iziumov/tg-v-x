package domain

import (
	"context"
	"iziumov/tv-v-x/internal/dto"
)

type UserRepository interface {
	GetByTgId(ctx context.Context, telegram_id int64) (*User, error)
	CreateUser(ctx context.Context, user dto.CreateUser) error
	Ban(ctx context.Context, telegram_id int64) error
	Unban(ctx context.Context, telegram_id int64) error
	AddAdmin(ctx context.Context, telegram_id int64) error
	RemoveAdmin(ctx context.Context, telegram_id int64) error
	GetAll(ctx context.Context, limit, offset int) ([]*User, error)
}

type JobRepository interface {
	CreateJob(ctx context.Context, job dto.CreateJob) (int, error)
	UpdateFinalStatus(ctx context.Context, job dto.UpdateJob) error
	GetFailedJobs(ctx context.Context) ([]*Job, error)
	GetHistoryByTgId(ctx context.Context, telegram_id int64) ([]*JobRecord, error)
}

type StatRepository interface {
	CreateStat(ctx context.Context, telegram_id int64) error
	IncrementSuccess(ctx context.Context, telegram_id int64, bytes int64) error
	IncrementFailed(ctx context.Context, telegram_id int64) error
	GetUserStats(ctx context.Context, telegram_id int64) (*Stats, error)
	GetGlobalStats(ctx context.Context) (*Stats, error)
}
