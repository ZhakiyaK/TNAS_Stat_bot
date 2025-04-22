package usecases

import (
	"strconv"
	"strings"
)

func ParseMemory(avail string) int {
	if len(avail) < 1 {
		return 0
	}

	// Удаляем все нечисловые символы в начале
	clean := strings.TrimLeft(avail, "0123456789")
	unit := strings.ToUpper(clean[len(clean)-1:])
	valueStr := strings.TrimRight(avail, "GMgmkKbB")

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0
	}

	switch unit {
	case "G":
		return value * 1024
	case "M":
		return value
	default:
		return 0
	}
}
