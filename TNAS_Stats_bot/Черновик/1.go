package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

func parseAvail(output string) (string, error) {
	lines := strings.Split(output, "\n")

	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		// Проверяем что это нужная точка монтирования
		if fields[5] == "/Volume3" {
			return fields[3], nil
		}
	}

	return "", fmt.Errorf("раздел /Volume3 не найден")
}

func SSHClient(ip string, port int, user, password string) (string, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Подключение
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err != nil {
		return "", fmt.Errorf("ошибка подключения: %v", err)
	}
	defer client.Close()

	// Выполнение команды
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("ошибка сессии: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput("df -h /Volume3")
	if err != nil {
		return "", fmt.Errorf("ошибка команды: %v", err)
	}

	// Парсинг вывода
	avail, err := parseAvail(string(output))
	if err != nil {
		return "", err
	}

	return avail, nil
}

func main() {
	ip := "192.168.10.15"
	port := 9222
	user := "TNAS-12E5"
	password := "Zaq12wsx"

	avail, err := SSHClient(ip, port, user, password)

	// Формируем статус
	status := "Не подключен"
	if err == nil {
		status = "Подключен"
	} else {
		log.Printf("Детали ошибки: %v", err) // Логируем ошибку для отладки
	}

	// Форматируем вывод
	fmt.Printf("Статус: %s\n", status)
	if avail != "" {
		fmt.Printf("Осталось места: %s\n", avail)
	} else {
		fmt.Println("Осталось места: недоступно")
	}
}
