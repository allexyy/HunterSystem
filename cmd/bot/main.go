package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/yourname/hunter-system/internal/infrastructure/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	ctx := context.Background()

	pool, err := database.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("database connection established")

	// TODO: инициализация Telegram Bot API клиента и хендлеров
	log.Println("bot starting... (handlers not implemented yet)")

	select {}
}
