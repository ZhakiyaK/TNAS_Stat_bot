package bot

import (
	"context"
	"fmt"
	"log/slog"

	"tg_bot/internal/adapters"
	"tg_bot/internal/config"
	"tg_bot/internal/usecases"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	cfg         *config.Config
	logger      *slog.Logger
	tgAdapter   *adapters.TelegramAdapter
	statsClient *usecases.StatsService
}

func NewBot(cfg *config.Config, logger *slog.Logger) *Bot {
	tgAdapter, err := adapters.NewTelegramAdapter(cfg.BotToken, cfg.ChatID)
	if err != nil {
		fmt.Println(err)
	}
	statsService := usecases.NewStatsService(cfg, logger)
	return &Bot{
		cfg:         cfg,
		logger:      logger,
		tgAdapter:   tgAdapter,
		statsClient: statsService,
	}
}

func (b *Bot) Run(ctx context.Context) {
	b.sendStartupMessage(ctx)
	defer b.sendShutdownMessage(ctx)

	updates := b.tgAdapter.GetUpdatesChan()
	b.logger.Info("Bot started and waiting for commands...")

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			b.handleUpdate(ctx, update)
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	if update.Message.Chat.ID != b.cfg.ChatID {
		b.logger.Warn("Unauthorized chat access attempt", "chatID", update.Message.Chat.ID)
		return
	}

	switch update.Message.Command() {
	case "stats":
		b.handleStatsCommand(ctx)
	}
}

func (b *Bot) handleStatsCommand(ctx context.Context) {
	status, avail := b.statsClient.GetStorageStatus()
	message := usecases.GenerateStatusMessage(status, avail)

	if err := b.tgAdapter.SendMessage(ctx, message); err != nil {
		b.logger.Error("Failed to send stats", "error", err)
	}
}

func (b *Bot) sendStartupMessage(ctx context.Context) {
	status, avail := b.statsClient.GetStorageStatus()
	message := usecases.GenerateStatusMessage(status, avail)

	if err := b.tgAdapter.SendMessage(ctx, message); err != nil {
		b.logger.Error("Failed to send startup status", "error", err)
	}
}

func (b *Bot) sendShutdownMessage(ctx context.Context) {
	message := "🔴 Бот TNAS Stat Bot остановлен"
	if err := b.tgAdapter.SendMessage(ctx, message); err != nil {
		b.logger.Error("Failed to send shutdown message", "error", err)
	}
	b.logger.Info("Bot stopped")
}
