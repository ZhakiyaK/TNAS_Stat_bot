package usecases

import (
	"fmt"
	"time"
)

func GenerateStatusMessage(status string, avail string) string {
	now := time.Now()
	message := fmt.Sprintf(
		"Информация по TNAS:\n\nДата: %s\nВремя: %s\nСтатус: %s\nОсталось места: %s",
		now.Format("02.01.2006"),
		now.Format("15:04"),
		status,
		avail,
	)

	// Добавляем предупреждение если меньше 100G
	if ParseMemory(avail) < 100*1024 { // 100G в мегабайтах
		message += "\n\n❗️❗️❗️Осталось мало места. Поменяйте диск"
	}

	return message
}
