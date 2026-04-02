package service

import (
	"context"
	"iziumov/tv-v-x/internal/dto"
	"iziumov/tv-v-x/internal/infra/redis"
	"log/slog"
	"strings"
)

type VideoService struct {
	jobSrv     *JobService
	queue      *redis.RedisClient
	downloader *DownloaderService
	logger     *slog.Logger
}

func NewVideoService(JobSrv *JobService, queue *redis.RedisClient, downloader *DownloaderService, logger *slog.Logger) *VideoService {
	return &VideoService{
		jobSrv:     JobSrv,
		queue:      queue,
		downloader: downloader,
		logger:     logger,
	}
}

func detectPlatform(url string) string {
	if strings.Contains(url, "twitter.com") || strings.Contains(url, "x.com") {
		return "twitter"
	}

	if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
		return "youtube"
	}

	return "unknown"
}

func (s *VideoService) SubmitDownload(ctx context.Context, userID int64, chatID int64, url string) error {
	platform := detectPlatform(url)

	jobID, err := s.jobSrv.CreateJob(ctx, dto.CreateJob{
		TgId:     userID,
		Url:      url,
		Platform: platform,
	})
	if err != nil {
		return err
	}

	err = s.queue.Enqueue(ctx, redis.JobMessage{
		JobID:    jobID,
		UserID:   userID,
		ChatID:   chatID,
		URL:      url,
		Platform: platform,
	})
	if err != nil {
		return err
	}

	s.logger.Info("job_submitted", "user_id", userID, "url", url)

	return nil
}
