package service

import (
	"context"
	"iziumov/tv-v-x/internal/domain"
)

type StatService struct {
	repo domain.StatRepository
}

func NewStatService(repo domain.StatRepository) *StatService {
	return &StatService{
		repo: repo,
	}
}

func (s *StatService) CreateStat(ctx context.Context, telegram_id int64) error {
	return s.repo.CreateStat(ctx, telegram_id)
}

func (s *StatService) GetStats(ctx context.Context, telegram_id int64) (*domain.Stats, error) {
	return s.repo.GetUserStats(ctx, telegram_id)
}

func (s *StatService) RecordSuccessStat(ctx context.Context, telegram_id, bytes int64) error {
	return s.repo.IncrementSuccess(ctx, telegram_id, bytes)
}

func (s *StatService) RecordFailureStat(ctx context.Context, telegram_id int64) error {
	return s.repo.IncrementFailed(ctx, telegram_id)
}

func (s *StatService) GetGlobalStats(ctx context.Context) (*domain.Stats, error) {
	return s.repo.GetGlobalStats(ctx)
}
