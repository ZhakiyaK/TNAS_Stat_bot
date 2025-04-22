package usecases

import (
	"context"

	"tg_bot/internal/entities"
	"tg_bot/internal/interfaces"
)

type SendStatsService struct {
	telegram interfaces.TelegramSender
}

func NewSendStatsService(telegram interfaces.TelegramSender) *SendStatsService {
	return &SendStatsService{telegram: telegram}
}

func (s *SendStatsService) SendStats(ctx context.Context, stats *entities.Stats) error {
	return s.telegram.SendMessage(ctx, stats.String())
}
