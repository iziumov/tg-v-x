package tg

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (t *TelegramClient) SetCommands(ctx context.Context) error {
	userCommands := []models.BotCommand{
		{Command: "start", Description: "Start bot"},
		{Command: "download", Description: "Start download from link"},
		{Command: "history", Description: "Your download history"},
		{Command: "stats", Description: "My stats"},
		{Command: "help", Description: "Help"},
	}

	_, err := t.Bot.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: userCommands,
	})
	if err != nil {
		return fmt.Errorf("failed to set user commands: %w", err)
	}

	adminCommands := append(userCommands,
		models.BotCommand{Command: "ban", Description: "Ban user by id"},
		models.BotCommand{Command: "unban", Description: "Unban user by id"},
		models.BotCommand{Command: "user_list", Description: "User lists"},
		models.BotCommand{Command: "global_stats", Description: "Global stats"},
	)

	for _, adminID := range t.adminIDs {
		_, err := t.Bot.SetMyCommands(ctx, &bot.SetMyCommandsParams{
			Commands: adminCommands,
			Scope:    &models.BotCommandScopeChat{ChatID: adminID},
		})

		if err != nil {
			return fmt.Errorf("failed to set admin commands for %d: %w", adminID, err)
		}
	}

	return nil
}
