package tg

import (
	"iziumov/tv-v-x/config"
	"iziumov/tv-v-x/internal/service"
	"log/slog"

	"github.com/go-telegram/bot"
)

type TelegramClient struct {
	*bot.Bot
	adminIDs []int64

	userSrv  *service.UserService
	jobSrv   *service.JobService
	statSrv  *service.StatService
	videoSrv *service.VideoService

	logger *slog.Logger
}

func NewTelegramClient(conf config.TGConfig, userSrv *service.UserService, jobSrv *service.JobService, statSrv *service.StatService, videoSrv *service.VideoService, log *slog.Logger) (*TelegramClient, error) {
	tc := TelegramClient{
		adminIDs: conf.Admins_Ids,
		userSrv:  userSrv,
		jobSrv:   jobSrv,
		statSrv:  statSrv,
		videoSrv: videoSrv,
		logger:   log,
	}

	opts := []bot.Option{
		bot.WithMiddlewares(
			tc.loggingMiddleware,
			tc.authMiddleware,
		),
	}

	b, err := bot.New(conf.Token, opts...)
	if err != nil {
		return nil, err
	}
	tc.Bot = b

	return &tc, nil
}
