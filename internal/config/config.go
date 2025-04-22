package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	BotToken     string
	ChatID       int64
	TNASIp       string
	TNASPort     int
	TNASUser     string
	TNASPassword string
}

func LoadConfig(logger *slog.Logger) (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")
	tnasIP := os.Getenv("TNAS_IP")
	tnasPortStr := os.Getenv("TNAS_PORT")
	tnasUser := os.Getenv("TNAS_USER")
	tnasPassword := os.Getenv("TNAS_PASSWORD")

	if botToken == "" || chatIDStr == "" || tnasIP == "" || tnasPortStr == "" || tnasUser == "" || tnasPassword == "" {
		logger.Error("Missing required environment variables")
		return nil, errors.New("missing required environment variables")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		logger.Error("Error converting CHAT_ID", "error", err)
		return nil, err
	}

	tnasPort, err := strconv.Atoi(tnasPortStr)
	if err != nil {
		logger.Error("Error converting TNAS_PORT", "error", err)
		return nil, err
	}

	return &Config{
		BotToken:     botToken,
		ChatID:       chatID,
		TNASIp:       tnasIP,
		TNASPort:     tnasPort,
		TNASUser:     tnasUser,
		TNASPassword: tnasPassword,
	}, nil
}
