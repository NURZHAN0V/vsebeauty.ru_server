#!/bin/bash

# Тест SMTP через bash (без telnet)
# Использование: ./test-smtp-bash.sh <host> <port> <email-to>

if [ -z "$3" ]; then
    echo "Использование: $0 <host> <port> <email-to>"
    echo "Пример: $0 localhost 2525 test@vsebeauty.ru"
    exit 1
fi

HOST=$1
PORT=$2
TO=$3

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

# Используем nc (netcat) если доступен
if command -v nc >/dev/null 2>&1; then
    echo "Используем netcat для теста..."
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
        echo "Тестовое письмо от $(date)"
        echo "."
        sleep 1
        echo "QUIT"
    } | nc $HOST $PORT
    
    echo ""
    echo "Проверьте логи сервера и базу данных"
else
    echo "netcat (nc) не установлен. Установите: apt-get install netcat"
    echo "Или используйте Python скрипт: python3 scripts/test-smtp-python.py $HOST $PORT $TO"
fi

