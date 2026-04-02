package tg

import (
	"github.com/go-telegram/bot"
)

func (t *TelegramClient) RegisterHandlers() {
	t.Bot.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, t.start)
	t.Bot.RegisterHandler(bot.HandlerTypeMessageText, "download", bot.MatchTypeCommand, t.download)
	t.Bot.RegisterHandler(bot.HandlerTypeMessageText, "history", bot.MatchTypeCommand, t.history)
	t.Bot.RegisterHandler(bot.HandlerTypeMessageText, "stats", bot.MatchTypeCommand, t.stats)
	t.Bot.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommand, t.help)

	//admin

	t.Bot.RegisterHandler(bot.HandlerTypeMessageText, "ban", bot.MatchTypeCommand, t.adminOnly(t.ban))
	t.Bot.RegisterHandler(bot.HandlerTypeMessageText, "unban", bot.MatchTypeCommand, t.adminOnly(t.unban))
	t.Bot.RegisterHandler(bot.HandlerTypeMessageText, "user_list", bot.MatchTypeCommand, t.adminOnly(t.userList))
	t.Bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, "user_list:", bot.MatchTypePrefix, t.adminOnly(t.userListCallback))
	t.Bot.RegisterHandler(bot.HandlerTypeMessageText, "global_stats", bot.MatchTypeCommand, t.adminOnly(t.globalStats))
}
