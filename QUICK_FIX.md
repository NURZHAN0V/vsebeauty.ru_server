# Быстрое исправление проблемы с SMTP

## Проблема
Ящики создаются с доменом `@vsebeauty.ru`, но SMTP сервер настроен на `tempmail.dev`.

## Решение

### 1. Исправьте docker-compose.yml
Уже исправлено: `MAIL_DOMAIN=vsebeauty.ru`

### 2. Перезапустите контейнер
```bash
docker-compose down
docker-compose up -d
```

### 3. Проверьте логи
```bash
docker logs -f tempmail-app
```

Должны увидеть:
```
SMTP-сервер запущен на порту 2525
Домен: vsebeauty.ru
```

### 4. Проверьте порт (Linux команды!)
```bash
# Правильная команда для Linux
netstat -tuln | grep :2525

# Или
ss -tuln | grep :2525
```

### 5. Отправьте тестовое письмо
```bash
# Создайте ящик
curl -X POST http://localhost:8080/api/v1/mailbox

# Отправьте тест (в контейнере или на хосте)
telnet localhost 2525
# Затем:
EHLO test
MAIL FROM:<test@example.com>
RCPT TO:<ваш-ящик@vsebeauty.ru>
DATA
Subject: Test

Тест
.
QUIT
```

### 6. Проверьте письмо
```bash
curl http://localhost:8080/api/v1/mailbox/{id}/messages
```

## Важно для получения писем из интернета

Если отправляете с Gmail/Outlook и письма не приходят:

1. **Нужны DNS записи:**
   - MX запись: `vsebeauty.ru` → ваш сервер
   - A запись: IP вашего сервера
   - SPF запись: разрешает отправку

2. **Порт 25 должен быть открыт** (для продакшена)
   - Измените `SMTP_PORT=25` в docker-compose.yml
   - Откройте порт в файрволе

3. **Проверьте логи при отправке:**
   ```bash
   docker logs -f tempmail-app
   ```
   Должны появиться записи о подключении.

