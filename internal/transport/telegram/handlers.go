package telegram

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strings"
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

func (b *Bot) handleHabitsList(ctx context.Context, bot2 *bot.Bot, update *models.Update) {
	tgID := update.Message.From.ID
	user, err := b.userService.GetUser(ctx, tgID)
	if err != nil {
		fmt.Errorf("User not found")
		return
	}
	habits, err := b.habitService.GetHabits(ctx, user.ID)
	if err != nil {
		fmt.Errorf("Habits not found")
		return
	}

	habitList := "Habit list :\n"
	for _, habit := range habits {
		habitList += fmt.Sprintf(
			"%d) %s\n%s",
			habit.ID, habit.Title, habit.Description.String)
	}

	b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   habitList})
}

func (b *Bot) handleAddHabit(ctx context.Context, bot2 *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	tgID := update.Message.From.ID
	user, err := b.userService.GetUser(ctx, tgID)
	if err != nil {
		fmt.Errorf("User not found")
		return
	}

	habitText := strings.Split(update.Message.Text, "|")
	//TODO: Add mapper and struct to income text
	h, err := b.habitService.CreateHabit(ctx, user.ID, habitText[0], habitText[1], habitText[2], habitText[3:])
	if err != nil {
		fmt.Errorf("Error creating habit: %v", err)
	}
	b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Create new habit: %s", h.Title),
	})
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "unknown command",
	})
}
