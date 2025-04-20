package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tg_bot/internal/adapters"
	"tg_bot/internal/config"
	"tg_bot/internal/entities"
	"tg_bot/internal/usecases"

	"github.com/joho/godotenv"
)

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger
}

func main() {
	// Инициализация логгера
	logger := setupLogger()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

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

	// Настройка режима отладки
	tgAdapter.SetDebug(true)

	ip := "192.168.10.15"
	port := 9222
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	avail, err := usecases.SSHClient(ip, port, user, password)

	// Формируем статус
	status := "Не подключен⛔️"
	if err == nil {
		status = "Подключен✅"
	} else {
		log.Printf("Детали ошибки: %v", err) // Логируем ошибку для отладки
	}

	// Форматируем вывод
	var output, output1 string
	date := time.Now().Format("02.01.2006")
	time := time.Now().Format("15:04")
	min := "100G"

	//availMemory, err := strconv.Atoi(avail)
	if min > avail {
		output1 += fmt.Sprintf("\n❗️❗️❗️Осталось мало место. Поменяйте диск")
	}

	output = fmt.Sprintf("Информация по TNAS:\n\nДата: %s\nВремя: %s\nСтатус: %s\n", date, time, status)
	if avail != "" {
		output += fmt.Sprintf("Осталось места: %s\n%s\n", avail, output1)
	} else {
		fmt.Println("Осталось места: недоступно")
	}

	// Отправка сообщения при запуске
	ctx := context.Background()
	if err := tgAdapter.SendMessage(ctx, output); err != nil {
		logger.Error("Ошибка отправки стартового сообщения", "error", err)
	} else {
		logger.Info("Стартовое сообщение отправлено")
	}

	// Инициализация сервисов
	sendStatsService := usecases.NewSendStatsService(tgAdapter)

	// Получение канала обновлений через адаптер
	updates := tgAdapter.GetUpdatesChan()

	// Обработка сигналов завершения
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("Бот запущен и ожидает команд...")

	for {
		select {
		case <-ctx.Done():
			stopMessage := "🔴 Бот TNAS Stat Bot остановлен"
			if err := tgAdapter.SendMessage(context.Background(), stopMessage); err != nil {
				logger.Error("Ошибка отправки сообщения об остановке", "error", err)
			}
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
		user     = "root"
		password = "Qwerty"
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
	if len(avail) < 1 {
		return 0
	}

	// Удаляем последний символ и преобразуем в число
	var value int
	_, err := fmt.Sscanf(avail, "%d", &value)
	if err != nil {
		return 0
	}
	return value
}
