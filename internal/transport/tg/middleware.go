package tg

import (
	"context"
	"iziumov/tv-v-x/internal/dto"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (t *TelegramClient) loggingMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil && update.Message.From != nil {
			t.logger.Info("incoming message",
				"user_id", update.Message.From.ID,
				"username", update.Message.From.Username,
				"text", update.Message.Text,
			)
		}
		next(ctx, b, update)
	}
}

func (t *TelegramClient) authMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil || update.Message.From == nil {
			next(ctx, b, update)
			return
		}

		from := update.Message.From

		user, err := t.userSrv.GetByTelegramID(ctx, from.ID)
		if err != nil {
			t.logger.Error("failed to get user", "error", err)
			next(ctx, b, update)
			return
		}

		if user == nil {
			err := t.userSrv.Register(ctx, dto.CreateUser{
				TGID:      from.ID,
				Username:  from.Username,
				FirstName: from.FirstName,
				LastName:  from.LastName,
			})
			if err != nil {
				t.logger.Error("failed to create user", "error", err)
			}

			next(ctx, b, update)
			return
		}

		if user.IsBanned {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "u are banned",
			})
			return
		}

		next(ctx, b, update)
	}
}

func (t *TelegramClient) adminOnly(handler bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var userID int64
		var chatID int64
		var username, firstName, lastName string

		if update.Message != nil && update.Message.From != nil {
			userID = update.Message.From.ID
			chatID = update.Message.Chat.ID
			username = update.Message.From.Username
			firstName = update.Message.From.FirstName
			lastName = update.Message.From.LastName
		} else if update.CallbackQuery != nil && update.CallbackQuery.From.ID != 0 {
			userID = update.CallbackQuery.From.ID
			username = update.CallbackQuery.From.Username
			firstName = update.CallbackQuery.From.FirstName
			lastName = update.CallbackQuery.From.LastName
			if update.CallbackQuery.Message.Type == models.MaybeInaccessibleMessageTypeMessage {
				if update.CallbackQuery.Message.Message != nil {
					chatID = update.CallbackQuery.Message.Message.Chat.ID
				}
			} else {
				if update.CallbackQuery.Message.InaccessibleMessage != nil {
					chatID = update.CallbackQuery.Message.InaccessibleMessage.Chat.ID
				}
			}
		} else {
			return
		}

		if !t.IsAdmin(userID) {
			if chatID != 0 {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "u dont have permissions",
				})
			}
			return
		}

		user, err := t.userSrv.GetByTelegramID(ctx, userID)
		if err != nil || user == nil {
			err := t.userSrv.Register(ctx, dto.CreateUser{
				TGID:      userID,
				Username:  username,
				FirstName: firstName,
				LastName:  lastName,
			})
			if err != nil {
				t.logger.Error("failed to register admin user", "error", err)
			}
		}

		handler(ctx, b, update)
	}
}

func (t *TelegramClient) IsAdmin(userID int64) bool {
	for _, id := range t.adminIDs {
		if id == userID {
			return true
		}
	}

	return false
}
