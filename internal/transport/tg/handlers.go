package tg

import (
	"context"
	"fmt"
	"iziumov/tv-v-x/internal/dto"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (t *TelegramClient) start(ctx context.Context, b *bot.Bot, update *models.Update) {
	from := update.Message.From
	tgID := from.ID
	chatID := update.Message.Chat.ID

	err := t.userSrv.Register(ctx, dto.CreateUser{
		TGID:      tgID,
		Username:  from.Username,
		FirstName: from.FirstName,
		LastName:  from.LastName,
	})
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Error register, try later",
		})
		return
	}

	err = t.statSrv.CreateStat(ctx, tgID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Error creating stat, try later",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("Hello %s!\n\n", from.FirstName),
	})
}

func (t *TelegramClient) download(ctx context.Context, b *bot.Bot, update *models.Update) {
	from := update.Message.From
	tgID := from.ID
	chatID := update.Message.Chat.ID

	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "You forgot to add link",
		})
		return
	}

	link := strings.TrimSpace(parts[1])

	err := t.videoSrv.SubmitDownload(ctx, tgID, chatID, link)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Something went wrong, try later",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Download started",
	})
}

func (t *TelegramClient) history(ctx context.Context, b *bot.Bot, update *models.Update) {
	tgID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	histories, err := t.jobSrv.GetHistory(ctx, tgID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Error getting history, try later",
		})
		return
	}

	if len(histories) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Your history is empty",
		})
		return
	}

	var sb strings.Builder
	sb.WriteString("Your download history:\n")
	for i, h := range histories {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, h.Url))
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   sb.String(),
	})
}

func (t *TelegramClient) stats(ctx context.Context, b *bot.Bot, update *models.Update) {
	tgID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	stats, err := t.statSrv.GetStats(ctx, tgID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Error getting stats, try later",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("Your stats:\nTotal Jobs: %d\nSuccess: %d\nFailed: %d\nTotal Bytes: %d", stats.TotalJobs, stats.SuccessJobs, stats.FailedJobs, stats.TotalBytes),
	})
}

func (t *TelegramClient) help(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Available commands:\n/start - Register and start\n/download <link> - Download video\n/history - Show history\n/stats - Show stats\n/help - Show this help",
	})
}

func (t *TelegramClient) ban(ctx context.Context, b *bot.Bot, update *models.Update) {
	parts := strings.SplitN(update.Message.Text, " ", 2)
	chatID := update.Message.Chat.ID

	if len(parts) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Usage: /ban <user_id>"})
		return
	}

	userIDToBan, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Invalid user ID"})
		return
	}

	err = t.userSrv.Ban(ctx, userIDToBan)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error banning user"})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: fmt.Sprintf("User %d banned", userIDToBan)})
}

func (t *TelegramClient) unban(ctx context.Context, b *bot.Bot, update *models.Update) {
	parts := strings.SplitN(update.Message.Text, " ", 2)
	chatID := update.Message.Chat.ID

	if len(parts) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Usage: /unban <user_id>"})
		return
	}

	userIDToUnban, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Invalid user ID"})
		return
	}

	err = t.userSrv.Unban(ctx, userIDToUnban)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error unbanning user"})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: fmt.Sprintf("User %d unbanned", userIDToUnban)})
}

func (t *TelegramClient) sendUserList(ctx context.Context, b *bot.Bot, chatID int64, offset, limit int, editMessageID int) {
	users, err := t.userSrv.GetAll(ctx, limit, offset)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error fetching users"})
		return
	}

	if len(users) == 0 && offset == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "No users found"})
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Users (offset %d):\n", offset))
	for _, u := range users {
		bannedStatus := ""
		if u.IsBanned {
			bannedStatus = "[BANNED]"
		}
		sb.WriteString(fmt.Sprintf("%d - %s (@%s) %s\n", u.ID, u.FirstName, u.Username, bannedStatus))
	}

	var keyboard [][]models.InlineKeyboardButton
	var row []models.InlineKeyboardButton

	if offset >= limit {
		row = append(row, models.InlineKeyboardButton{
			Text:         "<- Prev",
			CallbackData: fmt.Sprintf("user_list:%d", offset-limit),
		})
	}
	if len(users) == limit {
		row = append(row, models.InlineKeyboardButton{
			Text:         "Next ->",
			CallbackData: fmt.Sprintf("user_list:%d", offset+limit),
		})
	}

	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}

	replyMarkup := &models.InlineKeyboardMarkup{InlineKeyboard: keyboard}

	if editMessageID != 0 {
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatID,
			MessageID:   editMessageID,
			Text:        sb.String(),
			ReplyMarkup: replyMarkup,
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatID,
			Text:        sb.String(),
			ReplyMarkup: replyMarkup,
		})
	}
}

func (t *TelegramClient) userList(ctx context.Context, b *bot.Bot, update *models.Update) {
	t.sendUserList(ctx, b, update.Message.Chat.ID, 0, 5, 0)
}

func (t *TelegramClient) userListCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	data := update.CallbackQuery.Data
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return
	}

	offset, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	var chatID int64
	var messageID int

	if update.CallbackQuery.Message.Type == models.MaybeInaccessibleMessageTypeMessage {
		if update.CallbackQuery.Message.Message != nil {
			chatID = update.CallbackQuery.Message.Message.Chat.ID
			messageID = update.CallbackQuery.Message.Message.ID
		}
	} else {
		if update.CallbackQuery.Message.InaccessibleMessage != nil {
			chatID = update.CallbackQuery.Message.InaccessibleMessage.Chat.ID
			messageID = update.CallbackQuery.Message.InaccessibleMessage.MessageID
		}
	}

	if chatID != 0 {
		t.sendUserList(ctx, b, chatID, offset, 5, messageID)
	}
}

func (t *TelegramClient) globalStats(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	stats, err := t.statSrv.GetGlobalStats(ctx)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Error fetching global stats"})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("Global stats:\nTotal Jobs: %d\nSuccess: %d\nFailed: %d\nTotal Bytes: %d", stats.TotalJobs, stats.SuccessJobs, stats.FailedJobs, stats.TotalBytes),
	})
}
