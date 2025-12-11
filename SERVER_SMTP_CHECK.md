# Проверка SMTP на Linux сервере

## Важно: Домен должен совпадать!

В ящике используется домен `@vsebeauty.ru`, но в конфигурации указан `tempmail.dev`.
**Измените `MAIL_DOMAIN=vsebeauty.ru` в конфигурации!**

## Команды для Linux (не Windows!)

### 1. Проверка портов

```bash
# Проверка порта 2525
netstat -tuln | grep :2525

# Или через ss
ss -tuln | grep :2525

# Проверка порта 25
netstat -tuln | grep :25
```

### 2. Проверка процессов

```bash
# Поиск процесса
ps aux | grep tempmail

# Или через pgrep
pgrep -a tempmail
```

### 3. Проверка Docker контейнера

```bash
# Список контейнеров
docker ps

# Логи контейнера
docker logs tempmail-app

# Логи с фильтром SMTP
docker logs tempmail-app 2>&1 | grep -i smtp
```

### 4. Проверка переменных окружения

```bash
# В Docker контейнере
docker exec tempmail-app env | grep -E "SMTP|MAIL"

# Или если запущено напрямую
env | grep -E "SMTP|MAIL"
```

### 5. Тест SMTP подключения

```bash
# Проверка доступности порта
timeout 2 bash -c "echo > /dev/tcp/localhost/2525" && echo "Порт доступен" || echo "Порт недоступен"

# Тест через telnet (если установлен)
telnet localhost 2525
```

### 6. Отправка тестового письма

```bash
# Используйте скрипт
chmod +x scripts/test-smtp-linux.sh
./scripts/test-smtp-linux.sh localhost 2525 ваш-ящик@vsebeauty.ru

# Или через telnet вручную
telnet localhost 2525
# Затем введите:
EHLO test
MAIL FROM:<test@example.com>
RCPT TO:<ваш-ящик@vsebeauty.ru>
DATA
Subject: Test

Тест
.
QUIT
```

## Исправление конфигурации

### В Docker Compose

Измените `docker-compose.yml`:
```yaml
environment:
  - MAIL_DOMAIN=vsebeauty.ru  # Было: tempmail.dev
```

Затем перезапустите:
```bash
docker-compose down
docker-compose up -d
```

### В .env файле

```bash
MAIL_DOMAIN=vsebeauty.ru
SMTP_PORT=2525
```

## Проверка логов

```bash
# Логи Docker контейнера
docker logs -f tempmail-app

# Ищите строки:
# "SMTP-сервер запущен на порту 2525"
# "Домен: vsebeauty.ru"
# "Новое SMTP-соединение от ..."
# "MAIL FROM: ..."
# "RCPT TO: ..."
```

## Проверка базы данных

```bash
# Подключение к БД в контейнере
docker exec -it tempmail-postgres psql -U postgres -d tempmail

# Проверка ящиков
SELECT id, address, created_at, expires_at FROM mailboxes ORDER BY created_at DESC LIMIT 5;

# Проверка писем
SELECT id, mailbox_id, from_address, subject, received_at FROM messages ORDER BY received_at DESC LIMIT 10;
```

## Быстрая проверка

```bash
# Запустите скрипт проверки
chmod +x scripts/check-smtp.sh
./scripts/check-smtp.sh
```

