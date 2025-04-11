package config

import (
	"errors"
	"log/slog"
	"strconv"
)

type Config struct {
	BotToken string
	ChatID   int64
}

func LoadConfig(logger *slog.Logger) (*Config, error) {
	botToken := "7852675959:AAFTp2u66blhZrZ_AcUX1Zy-_RlXN8n4yug"
	chatIDStr := "-1002548885029"

	if botToken == "" || chatIDStr == "" {
		logger.Error("Необходимые переменные окружения отсутствуют", "BOT_TOKEN", botToken, "CHAT_ID", chatIDStr)
		return nil, errors.New("необходимые переменные окружения отсутствуют")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		logger.Error("Ошибка преобразования CHAT_ID в int64", "error", err)
		return nil, err
	}

	return &Config{
		BotToken: botToken,
		ChatID:   chatID,
	}, nil
}
