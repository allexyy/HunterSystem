package main

import (
	"context"
	"github.com/yourname/hunter-system/internal/app"
	"github.com/yourname/hunter-system/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to build app: %v", err)
	}
	defer a.Close()

	a.Run(ctx)
	log.Println("shutdown complete")
}
