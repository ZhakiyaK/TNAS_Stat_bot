package usecases

import (
	"fmt"
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
