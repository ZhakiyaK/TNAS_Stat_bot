package usecases

import (
	"fmt"
	"time"
)

func GenerateStartupMessage() string {
	return "🟢 TNAS Stat Bot инициализирован"
}

func GenerateStatusMessage(status string, avail string) string {
	now := time.Now()
	return fmt.Sprintf(
		"Информация по TNAS:\n\nДата: %s\nВремя: %s\nСтатус: %s\nОсталось места: %s",
		now.Format("02.01.2006"),
		now.Format("15:04"),
		status,
		avail,
	)
}
