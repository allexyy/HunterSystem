package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseURL        string
	APIPort            string
	TelegramBotToken   string
	TelegramBotHookUrl string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	if cfg.TelegramBotToken == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is not set")
	}
	return cfg, nil
}
