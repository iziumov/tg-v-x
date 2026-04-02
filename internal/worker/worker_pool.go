package worker

import (
	"context"
	"iziumov/tv-v-x/internal/infra/redis"
	"iziumov/tv-v-x/internal/service"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type WorkerPool struct {
	count      int
	queue      *redis.RedisClient
	downloader *service.DownloaderService
	jobSrv     *service.JobService
	statSrv    *service.StatService
	bot        *bot.Bot
	logger     *slog.Logger
}

func NewWorkerPool(count int,
	queue *redis.RedisClient,
	downloader *service.DownloaderService,
	jobSrv *service.JobService,
	statSrv *service.StatService,
	bot *bot.Bot,
	logger *slog.Logger) *WorkerPool {
	return &WorkerPool{
		count:      count,
		queue:      queue,
		downloader: downloader,
		jobSrv:     jobSrv,
		statSrv:    statSrv,
		bot:        bot,
		logger:     logger,
	}
}

func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.count; i++ {
		go p.run(ctx, i)
	}
	p.logger.Info("worker pool started", "count", p.count)
}

func (p *WorkerPool) run(ctx context.Context, id int) {
	p.logger.Info("worker started", "worker_id", id)

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("worker stopped", "worker_id", id)
			return
		default:
			p.processNext(ctx, id)
		}
	}
}

func (p *WorkerPool) processNext(ctx context.Context, id int) {
	job, err := p.queue.Dequeue(ctx, 5*time.Second)
	if err != nil {
		return
	}

	p.logger.Info("proccessing job", "worker_id", id, "job_id", job.JobID, "url", job.URL)

	jobCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	result, err := p.downloader.Download(jobCtx, job.URL, job.JobID)
	if err != nil {
		p.handleFailure(ctx, job, err)
		return
	}
	defer p.downloader.CleanUp(result.FilePath)

	err = p.sendVideo(ctx, job.ChatID, result.FilePath)
	if err != nil {
		p.handleFailure(ctx, job, err)
		return
	}

	if err := p.jobSrv.MarkDone(ctx, job.JobID); err != nil {
		p.logger.Error("failed to mak job done", "error", err)
	}

	p.logger.Info("recording success stat", "user_id", job.UserID, "file_size", result.FileSize)
	if err := p.statSrv.RecordSuccessStat(ctx, job.UserID, result.FileSize); err != nil {
		p.logger.Error("failed to record success stat", "error", err)
	} else {
		p.logger.Info("success stat recorded", "user_id", job.UserID, "file_size", result.FileSize)
	}

	p.logger.Info("job completed", "worker_id", id, "job_id", job.JobID, "file_size", result.FileSize)
}

func (p *WorkerPool) handleFailure(ctx context.Context, job *redis.JobMessage, err error) error {
	p.logger.Error("job failed", "job_id", job.JobID, "url", job.URL, "error", err)

	if markErr := p.jobSrv.MarkFailed(ctx, job.JobID, err.Error()); markErr != nil {
		p.logger.Error("failed to mark job failed", "error", markErr)
	}

	if statErr := p.statSrv.RecordFailureStat(ctx, job.UserID); statErr != nil {
		p.logger.Error("failed to record stats", "error", statErr)
	}

	p.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: job.ChatID,
		Text:   "cant download video " + err.Error(),
	})

	return nil
}

func (p *WorkerPool) sendVideo(ctx context.Context, chantID int64, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = p.bot.SendVideo(ctx, &bot.SendVideoParams{
		ChatID: chantID,
		Video: &models.InputFileUpload{
			Filename: filepath.Base(filePath),
			Data:     file,
		},
	})

	return err
}
