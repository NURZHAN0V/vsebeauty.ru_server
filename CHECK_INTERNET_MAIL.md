# Проверка получения писем из интернета

## Проблема
Письма из Яндекс/Gmail не доходят, потому что SMTP сервер работает только локально.

## Проверьте сейчас:

### 1. Проверьте логи - пришло ли письмо:
```bash
docker logs tempmail-app --tail 50 | grep -i "smtp\|mail\|rcpt\|соединение"
```

Если письмо пришло, должны быть записи:
- "Новое SMTP-соединение от ..."
- "MAIL FROM: ..."
- "RCPT TO: ..."

### 2. Проверьте базу данных:
```bash
docker exec -it tempmail-postgres psql -U postgres -d tempmail -c "SELECT id, from_address, subject, received_at FROM messages WHERE mailbox_id = (SELECT id FROM mailboxes WHERE address = '3akqus4qwy@vsebeauty.ru') ORDER BY received_at DESC;"
```

### 3. Проверьте MX записи DNS:
```bash
# Проверьте MX запись для vsebeauty.ru
dig MX vsebeauty.ru
# или
nslookup -type=MX vsebeauty.ru
```

## Что нужно для получения писем из интернета:

### 1. Настроить DNS записи:
- **MX запись**: `vsebeauty.ru` → указывает на ваш сервер
- **A запись**: IP адрес вашего сервера
- **SPF запись**: `v=spf1 ip4:ВАШ-IP ~all`

### 2. Открыть порт 25:
```bash
# Проверьте, открыт ли порт 25
netstat -tuln | grep :25

# Если нет, нужно открыть в файрволе
```

### 3. Изменить порт на 25 в docker-compose.yml:
```yaml
- SMTP_PORT=25  # Вместо 2525
```

И перезапустить:
```bash
docker compose down
docker compose up -d
```

### 4. Проверить доступность из интернета:
```bash
# С другого сервера или через онлайн-сервис
telnet ваш-ip-адрес 25
```

## Временное решение для тестирования:

Пока DNS не настроен, можно использовать локальную отправку через telnet (как мы делали раньше).

## Проверка сейчас:

Выполните команды выше, чтобы понять:
1. Пришло ли письмо на сервер (логи)
2. Сохранилось ли в БД
3. Настроены ли DNS записи

