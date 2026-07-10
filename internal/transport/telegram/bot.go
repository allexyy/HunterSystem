package telegram

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/yourname/hunter-system/internal/user"
)

type Bot struct {
	api         *bot.Bot
	userService *user.Service
}

func New(token string, userService *user.Service) (*Bot, error) {
	b, err := bot.New(token,
		bot.WithDefaultHandler(defaultHandler),
	)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	wrapper := &Bot{api: b, userService: userService}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, wrapper.handleStart)

	return wrapper, nil
}

func (b *Bot) Run(ctx context.Context) {
	b.api.Start(ctx)
}
