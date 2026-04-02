package service

import (
	"context"
	"iziumov/tv-v-x/internal/domain"
	"iziumov/tv-v-x/internal/dto"
)

type JobService struct {
	repo domain.JobRepository
}

func NewJobService(repo domain.JobRepository) *JobService {
	return &JobService{
		repo: repo,
	}
}

func (s *JobService) CreateJob(ctx context.Context, job dto.CreateJob) (int, error) {
	return s.repo.CreateJob(ctx, job)
}

func (s *JobService) MarkDone(ctx context.Context, jobID int) error {
	return s.repo.UpdateFinalStatus(ctx, dto.UpdateJob{
		ID:     jobID,
		Status: "done",
	})
}

func (s *JobService) MarkFailed(ctx context.Context, jobID int, errMsg string) error {
	return s.repo.UpdateFinalStatus(ctx, dto.UpdateJob{
		ID:       jobID,
		Status:   "failed",
		ErrorMsg: &errMsg,
	})
}

func (s *JobService) GetFailedJobs(ctx context.Context) ([]*domain.Job, error) {
	return s.repo.GetFailedJobs(ctx)
}

func (s *JobService) GetHistory(ctx context.Context, telegram_id int64) ([]*domain.JobRecord, error) {
	return s.repo.GetHistoryByTgId(ctx, telegram_id)
}
