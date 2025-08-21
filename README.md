# TNAS Stat Bot

Telegram бот для мониторинга состояния сетевого хранилища TNAS через SSH-соединение с отправкой уведомлений и статистики.

## Основные функции

- 📊 Мониторинг доступности сетевого хранилища
- 💾 Контроль свободного места на дисках
- 🔔 Автоматические уведомления при нехватке места
- ⚡ Команда `/stats` для получения текущей статистики
- 📝 Логирование в формате JSON
- 🔒 Безопасное хранение конфигурации через переменные окружения

## Технологический стек

- **Язык**: Go 1.21+
- **Telegram API**: [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
- **SSH-клиент**: `golang.org/x/crypto/ssh`
- **Конфигурация**: environment variables через `github.com/joho/godotenv`

## Установка и настройка


1. Клонируйте репозиторий:
```bash
git clone https://github.com/ZhakiyaK/TNAS_Stat_bot
cd TNAS_Stat_Bot
```
2. Установите зависимости:
```bash
go mod download
```
3. Создайте файл .env в корне проекта:
```env
BOT_TOKEN=your_telegram_bot_token
CHAT_ID=your_chat_id
USER=ssh_username
PASSWORD=ssh_password
```
