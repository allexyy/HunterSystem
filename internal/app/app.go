package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourname/hunter-system/internal/config"
	"github.com/yourname/hunter-system/internal/db"
	"github.com/yourname/hunter-system/internal/infrastructure/database"
	"github.com/yourname/hunter-system/internal/transport/telegram"
	"github.com/yourname/hunter-system/internal/user"
)

type App struct {
	cfg  *config.Config
	pool *pgxpool.Pool
	bot  *telegram.Bot
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	pool, err := database.New(ctx, cfg.DatabaseURL)
	queries := db.New(pool)
	txManager := database.NewTxManager(pool)
	userService := user.NewService(queries, txManager)
	bot, err := telegram.New(cfg.TelegramBotToken, userService)

	return &App{cfg: cfg, pool: pool, bot: bot}, err
}

func (a *App) Run(ctx context.Context) {
	a.bot.Run(ctx)
}

func (a *App) Close() {
	defer a.pool.Close()
}
