package usecases

import (
	"log/slog"

	"tg_bot/internal/config"
)

type StatsService struct {
	cfg    *config.Config
	logger *slog.Logger
}

func NewStatsService(cfg *config.Config, logger *slog.Logger) *StatsService {
	return &StatsService{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *StatsService) GetStorageStatus() (string, string) {
	avail, err := SSHClient(
		s.cfg.TNASIp,
		s.cfg.TNASPort,
		s.cfg.TNASUser,
		s.cfg.TNASPassword,
	)

	if err != nil {
		s.logger.Error("Failed to get storage status",
			"error", err,
			"ip", s.cfg.TNASIp,
			"port", s.cfg.TNASPort)
		return "Не подключен⛔️", "недоступно"
	}
	s.logger.Info("Успешно получен статус хранилища", "avail", avail)
	return "Подключен✅", avail
}
