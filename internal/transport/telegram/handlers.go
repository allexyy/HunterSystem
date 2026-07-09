package telegram

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yourname/hunter-system/internal/db"
	"log"
	"strconv"
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

func (b *Bot) handleQuestList(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	tgID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	user, err := b.userService.GetUser(ctx, tgID)
	if err != nil {
		log.Printf("handleQuestList: get user tg=%d: %v", tgID, err)
		b.reply(ctx, chatID, "Что-то пошло не так, попробуй позже")
		return
	}

	quests, err := b.questService.GetQuestList(ctx, user.ID)
	if err != nil {
		log.Printf("handleQuestList: ensure quests user=%d: %v", user.ID, err)
		b.reply(ctx, chatID, "Не удалось получить квесты, попробуй позже")
		return
	}

	if len(quests) == 0 {
		b.reply(ctx, chatID, "На сегодня квестов нет. Создай привычку: /newhabit")
		return
	}

	var sb strings.Builder
	sb.WriteString("⚔️ Квесты на сегодня:\n\n")

	rows := make([][]models.InlineKeyboardButton, 0, len(quests))
	for _, quest := range quests {
		mark := "⬜"
		if quest.Status == db.QuestStatusCompleted {
			mark = "✅"
		}
		fmt.Fprintf(&sb, "%s %s (+%d XP, +%d 💰)\n", mark, quest.Title, quest.XpReward, quest.GoldReward)
		if quest.Description.Valid && quest.Description.String != "" {
			fmt.Fprintf(&sb, "   %s\n", quest.Description.String)
		}

		if quest.Status != db.QuestStatusPending {
			continue
		}
		rows = append(rows, []models.InlineKeyboardButton{
			{
				Text:         "✅ " + quest.Title,
				CallbackData: fmt.Sprintf("done:%d", quest.ID),
			},
		})
	}

	_, err = b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   sb.String(),
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: rows,
		},
	})
	if err != nil {
		log.Printf("handleQuestList: send reply chat=%d: %v", chatID, err)
	}
}

func (b *Bot) handleDoneCallback(ctx context.Context, bot2 *bot.Bot, update *models.Update) {
	cb := update.CallbackQuery
	if cb == nil {
		return
	}

	_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
	})

	idStr := strings.TrimPrefix(cb.Data, "done:")
	questID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("bad callback data %q: %v", cb.Data, err)
		return
	}

	b.questService.CompleteQuest(ctx, questID)
	b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		Text: "Квест выполнен! +20 XP",
	})
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

func (b *Bot) reply(ctx context.Context, chatID int64, text string) {
	if _, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		log.Printf("send message chat=%d: %v", chatID, err)
	}
}
