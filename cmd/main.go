package main

import (
	"context"
	"iziumov/tv-v-x/config"
	"iziumov/tv-v-x/internal/infra/redis"
	"iziumov/tv-v-x/internal/transport/tg"
	"iziumov/tv-v-x/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	postgresDB "iziumov/tv-v-x/internal/infra/postgres"
	"iziumov/tv-v-x/internal/service"
	"iziumov/tv-v-x/internal/worker"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		return
	}

	level := logger.ParseLevel(cfg.Logger)
	log := logger.NewLogger(level)

	db, err := postgresDB.NewDB(cfg.DB)
	if err != nil {
		log.Error("failed connect to db", "error", err)
		return
	}
	if err := db.Migrate(cfg.DB); err != nil {
		log.Error("failed to migrate db", "error", err)
		return
	}

	redisCli, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Error("failed connect to redis", "error", err)
		return
	}

	userRepo := postgresDB.NewUserRepo(db)
	statRepo := postgresDB.NewStatRepo(db)
	jobRepo := postgresDB.NewJobRepo(db)

	userService := service.NewUserService(userRepo)
	statService := service.NewStatService(statRepo)
	jobService := service.NewJobService(jobRepo)
	downloaderService := service.NewDownloaderService(cfg.Ytdlp)
	videoService := service.NewVideoService(jobService, redisCli, downloaderService, log)

	tgClient, err := tg.NewTelegramClient(cfg.TG, userService, jobService, statService, videoService, log)
	if err != nil {
		log.Error("failed to create telegram bot", "error", err)
		return
	}
	tgClient.RegisterHandlers()
	if err := tgClient.SetCommands(ctx); err != nil {
		log.Error("failed to set commands", "error", err)
	}

	pool := worker.NewWorkerPool(cfg.Workers, redisCli, downloaderService, jobService, statService, tgClient.Bot, log)
	pool.Start(ctx)

	go tgClient.Start(ctx)

	<-ctx.Done()
	log.Info("shutting down gracefully....")
}
