package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) handleStart(ctx context.Context, bot2 *bot.Bot, update *models.Update) {
	tgID := update.Message.From.ID
	username := update.Message.From.Username

	_, err := b.userService.RegisterUser(ctx, tgID, username)
	if err != nil {
		b.api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Something went wrong, please try again later."})
		return
	}
	b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hunter System initialized.\n\nWelcome, Hunter.\n\nLevel 1\n\nXP: 0\n\nGold: 0\n\nRank: E"})
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
