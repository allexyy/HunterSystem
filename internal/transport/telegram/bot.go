package telegram

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/yourname/hunter-system/internal/habit"
	"github.com/yourname/hunter-system/internal/quest"
	"github.com/yourname/hunter-system/internal/user"
)

type Bot struct {
	api          *bot.Bot
	userService  *user.Service
	habitService *habit.Service
	questService *quest.Service
}

func New(token string, userService *user.Service, habitService *habit.Service, questService *quest.Service) (*Bot, error) {
	b, err := bot.New(token,
		bot.WithDefaultHandler(defaultHandler),
	)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	wrapper := &Bot{api: b, userService: userService, habitService: habitService, questService: questService}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, wrapper.handleStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/newhabit", bot.MatchTypePrefix, wrapper.handleAddHabit)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/habits", bot.MatchTypePrefix, wrapper.handleHabitsList)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/today", bot.MatchTypePrefix, wrapper.handleQuestList)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/newquest", bot.MatchTypePrefix, defaultHandler)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "done:", bot.MatchTypePrefix, wrapper.handleDoneCallback)

	return wrapper, nil
}

func (b *Bot) Run(ctx context.Context) {
	b.api.Start(ctx)
}
