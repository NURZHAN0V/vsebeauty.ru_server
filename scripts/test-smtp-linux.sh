#!/bin/bash

# Тест SMTP сервера для Linux

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Использование: ./test-smtp-linux.sh <smtp-host> <email-to>"
    echo "Пример: ./test-smtp-linux.sh localhost test@vsebeauty.ru"
    exit 1
fi

HOST=${1:-localhost}
PORT=${2:-2525}
TO=$3

if [ -z "$TO" ]; then
    echo "Ошибка: не указан email получателя"
    echo "Использование: ./test-smtp-linux.sh <smtp-host> <smtp-port> <email-to>"
    exit 1
fi

echo "Тестирование SMTP сервера $HOST:$PORT"
echo "Получатель: $TO"
echo ""

# Проверка доступности порта
if ! timeout 2 bash -c "echo > /dev/tcp/$HOST/$PORT" 2>/dev/null; then
    echo "Ошибка: не удалось подключиться к $HOST:$PORT"
    echo "Убедитесь, что SMTP сервер запущен"
    exit 1
fi

echo "✓ Порт доступен"
echo ""

# Отправка тестового письма через telnet
echo "Отправка тестового письма..."
{
    sleep 1
    echo "EHLO test"
    sleep 1
    echo "MAIL FROM:<test@example.com>"
    sleep 1
    echo "RCPT TO:<$TO>"
    sleep 1
    echo "DATA"
    sleep 1
    echo "Subject: Test Message"
    echo ""
    echo "Это тестовое письмо от $(date)"
    echo "."
    sleep 1
    echo "QUIT"
} | telnet $HOST $PORT 2>/dev/null

echo ""
echo "Проверьте логи сервера и базу данных"

