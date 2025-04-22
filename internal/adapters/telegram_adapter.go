package adapters

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramAdapter struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramAdapter(botToken string, chatID int64) (*TelegramAdapter, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	return &TelegramAdapter{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (t *TelegramAdapter) SendMessage(ctx context.Context, message string) error {
	msg := tgbotapi.NewMessage(t.chatID, message)
	_, err := t.bot.Send(msg)
	return err
}

// Новый метод для получения канала обновлений
func (t *TelegramAdapter) GetUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return t.bot.GetUpdatesChan(u)
}

// Новый метод для доступа к конфигурации бота
func (t *TelegramAdapter) SetDebug(debug bool) {
	t.bot.Debug = debug
}
