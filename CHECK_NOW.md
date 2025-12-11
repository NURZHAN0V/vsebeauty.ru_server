# Проверка SMTP прямо сейчас

## ✅ Что уже работает:
- Порт 2525 слушается (видно из netstat)
- Контейнер перезапущен
- Домен исправлен на vsebeauty.ru

## 🔍 Проверьте логи контейнера:

```bash
# Последние логи
docker logs tempmail-app --tail 50

# Логи с фильтром SMTP
docker logs tempmail-app 2>&1 | grep -i smtp

# Должны увидеть:
# "SMTP-сервер запущен на порту 2525"
# "Домен: vsebeauty.ru"
```

## 📧 Отправьте тестовое письмо:

### Вариант 1: Python (если установлен)
```bash
python3 scripts/test-smtp-python.py localhost 2525 ooq961x6sj@vsebeauty.ru
```

### Вариант 2: netcat (если установлен)
```bash
chmod +x scripts/test-smtp-bash.sh
./scripts/test-smtp-bash.sh localhost 2525 ooq961x6sj@vsebeauty.ru
```

### Вариант 3: Установите telnet
```bash
apt-get update
apt-get install -y telnet

# Затем
telnet localhost 2525
# Введите:
EHLO test
MAIL FROM:<test@example.com>
RCPT TO:<ooq961x6sj@vsebeauty.ru>
DATA
Subject: Test

Тест
.
QUIT
```

### Вариант 4: Через Docker exec
```bash
# Если в контейнере есть Python
docker exec tempmail-app python3 /path/to/test-smtp-python.py localhost 2525 ooq961x6sj@vsebeauty.ru
```

## 🔍 Проверьте базу данных:

```bash
# Подключитесь к БД
docker exec -it tempmail-postgres psql -U postgres -d tempmail

# Проверьте ящик
SELECT id, address, created_at, expires_at FROM mailboxes WHERE address = 'ooq961x6sj@vsebeauty.ru';

# Проверьте письма
SELECT id, mailbox_id, from_address, subject, received_at FROM messages ORDER BY received_at DESC LIMIT 10;
```

## 📊 Проверьте через API:

```bash
# Получите список писем
curl http://localhost:8080/api/v1/mailbox/d714e111-8041-4c4b-a5c4-cdce44f77e4a/messages

# Должен вернуть массив писем (может быть пустым [])
```

## ⚠️ Важно для получения писем из интернета:

Если отправляете с Gmail/Outlook и письма не приходят:

1. **SMTP сервер работает только локально** - для получения из интернета нужно:
   - Настроить DNS (MX записи для vsebeauty.ru)
   - Открыть порт 25 в файрволе
   - Изменить SMTP_PORT=25 в docker-compose.yml

2. **Проверьте логи при отправке:**
   ```bash
   docker logs -f tempmail-app
   ```
   При отправке письма должны появиться:
   - "Новое SMTP-соединение от ..."
   - "MAIL FROM: ..."
   - "RCPT TO: ..."
   - "Получение данных письма..."

## 🚀 Быстрая проверка всего:

```bash
# 1. Проверка порта
netstat -tuln | grep :2525

# 2. Проверка логов
docker logs tempmail-app --tail 20 | grep -i smtp

# 3. Проверка переменных
docker exec tempmail-app env | grep -E "SMTP|MAIL"

# 4. Создание ящика
curl -X POST http://localhost:8080/api/v1/mailbox

# 5. Отправка теста (выберите один вариант выше)

# 6. Проверка писем
curl http://localhost:8080/api/v1/mailbox/{id}/messages
```

