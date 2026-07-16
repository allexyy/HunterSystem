package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourname/hunter-system/internal/config"
	"github.com/yourname/hunter-system/internal/db"
	"github.com/yourname/hunter-system/internal/habit"
	"github.com/yourname/hunter-system/internal/infrastructure/database"
	"github.com/yourname/hunter-system/internal/quest"
	"github.com/yourname/hunter-system/internal/scheduler"
	"github.com/yourname/hunter-system/internal/transport/telegram"
	"github.com/yourname/hunter-system/internal/user"
	"sync"
)

type App struct {
	cfg       *config.Config
	pool      *pgxpool.Pool
	bot       *telegram.Bot
	scheduler *scheduler.Scheduler
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	pool, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	queries := db.New(pool)
	txManager := database.NewTxManager(pool)
	userService := user.NewService(queries, txManager)
	habitService := habit.NewService(queries, txManager)
	questService := quest.NewService(queries, txManager)

	bot, err := telegram.New(cfg.TelegramBotToken, userService, habitService, questService)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}
	sch := scheduler.NewScheduler(userService, questService)

	return &App{cfg: cfg, pool: pool, bot: bot, scheduler: sch}, err
}

func (a *App) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.scheduler.Run(ctx)
	}()
	a.bot.Run(ctx)
	wg.Wait()
}

func (a *App) Close() {
	defer a.pool.Close()
}
