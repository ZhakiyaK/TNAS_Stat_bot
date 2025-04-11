package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"tg_bot/internal/adapters"
	"tg_bot/internal/config"
	"tg_bot/internal/entities"
	"tg_bot/internal/usecases"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Инициализация логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация Telegram адаптера
	tgAdapter, err := adapters.NewTelegramAdapter(cfg.BotToken, cfg.ChatID)
	if err != nil {
		log.Fatalf("Ошибка инициализации Telegram адаптера: %v", err)
	}

	// Инициализация сервисов
	sendStatsService := usecases.NewSendStatsService(tgAdapter)

	// Настройка обработчика команд
	//bot := tgAdapter.Bot
	bot := tgAdapter.Bot
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	//u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Обработка сигналов завершения
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("Бот запущен и ожидает команд...")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Остановка бота...")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			// Обработка только команд из разрешенного чата
			if update.Message.Chat.ID != cfg.ChatID {
				continue
			}

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "stats":
					handleStatsCommand(ctx, sendStatsService, logger)
				}
			}
		}
	}
}

func handleStatsCommand(ctx context.Context, sender *usecases.SendStatsService, logger *slog.Logger) {
	// Параметры подключения к TNAS
	const (
		ip       = "192.168.10.15"
		port     = 22
		user     = "TNAS-12E5"
		password = "Zaq12wsx"
	)

	// Получение статистики
	avail, err := usecases.SSHClient(ip, port, user, password)
	status := "Не подключен"
	if err == nil {
		status = "Подключен"
	}

	// Формирование объекта статистики
	stats := &entities.Stats{
		Status:     status,
		MemoryLeft: parseMemory(avail), // Нужно реализовать преобразование
	}

	// Отправка статистики
	if err := sender.SendStats(ctx, stats); err != nil {
		logger.Error("Ошибка отправки статистики", "error", err)
	}
}

// Вспомогательная функция для преобразования строки в число
func parseMemory(avail string) int {
	// Реализуйте преобразование строки вида "291G" в мегабайты
	// Пример упрощенной реализации:
	if len(avail) < 1 {
		return 0
	}

	// Удаляем последний символ и преобразуем в число
	// В реальной реализации нужно учитывать единицы измерения (G/M)
	var value int
	_, err := fmt.Sscanf(avail, "%d", &value)
	if err != nil {
		return 0
	}
	return value
}
