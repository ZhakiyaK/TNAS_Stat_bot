package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"tg_bot/cmd/bot"
	"tg_bot/internal/config"

	"github.com/joho/godotenv"
)

func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func main() {
	logger := setupLogger()
	defer handlePanic(logger)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application := bot.NewBot(cfg, logger)
	go application.Run(ctx)

	<-ctx.Done()
}

func handlePanic(logger *slog.Logger) {
	if r := recover(); r != nil {
		logger.Error("Panic occurred", "recover", r)
	}
}
