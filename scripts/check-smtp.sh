#!/bin/bash

echo "=== Проверка SMTP сервера ==="
echo ""

# Проверка порта 2525
echo "1. Проверка порта 2525:"
if netstat -tuln | grep -q ":2525"; then
    echo "   ✓ Порт 2525 слушается"
    netstat -tuln | grep ":2525"
else
    echo "   ✗ Порт 2525 НЕ слушается"
fi

echo ""

# Проверка порта 25
echo "2. Проверка порта 25:"
if netstat -tuln | grep -q ":25 "; then
    echo "   ✓ Порт 25 слушается"
    netstat -tuln | grep ":25 "
else
    echo "   ✗ Порт 25 НЕ слушается (нормально, если используете 2525)"
fi

echo ""

# Проверка процессов
echo "3. Проверка процессов:"
if pgrep -f "tempmail\|smtp" > /dev/null; then
    echo "   ✓ Процесс найден:"
    ps aux | grep -E "tempmail|smtp" | grep -v grep
else
    echo "   ✗ Процесс не найден"
fi

echo ""

# Проверка переменных окружения
echo "4. Проверка переменных окружения:"
if [ -f .env ]; then
    echo "   .env файл найден"
    grep -E "SMTP_PORT|MAIL_DOMAIN" .env 2>/dev/null || echo "   Переменные не найдены в .env"
else
    echo "   .env файл не найден"
fi

echo ""

# Тест подключения к SMTP
echo "5. Тест подключения к SMTP (localhost:2525):"
if timeout 2 bash -c "echo > /dev/tcp/localhost/2525" 2>/dev/null; then
    echo "   ✓ Подключение к порту 2525 успешно"
else
    echo "   ✗ Не удалось подключиться к порту 2525"
fi

echo ""

# Проверка логов (если есть)
echo "6. Последние логи SMTP (если доступны):"
if [ -f logs/app.log ]; then
    echo "   Последние 10 строк с 'SMTP':"
    grep -i smtp logs/app.log | tail -10
elif journalctl -u tempmail > /dev/null 2>&1; then
    echo "   Последние логи systemd:"
    journalctl -u tempmail -n 20 --no-pager | grep -i smtp
else
    echo "   Логи не найдены"
fi

echo ""
echo "=== Проверка завершена ==="

