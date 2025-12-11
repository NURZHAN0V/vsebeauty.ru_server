# Настройка порта 25 для получения писем из интернета

## ✅ MX запись настроена:
```
vsebeauty.ru.  MX  10  mail.vsebeauty.ru.
```

## Проверьте A запись для mail.vsebeauty.ru:

```bash
dig A mail.vsebeauty.ru
# или
host mail.vsebeauty.ru
```

Должна указывать на IP вашего сервера.

## Настройка порта 25:

### 1. Измените docker-compose.yml:

```yaml
environment:
  - SMTP_PORT=25  # Вместо 2525
```

### 2. Откройте порт 25 в файрволе:

```bash
# Проверьте файрвол
ufw status

# Откройте порт 25
ufw allow 25/tcp

# Или для iptables
iptables -A INPUT -p tcp --dport 25 -j ACCEPT
```

### 3. Перезапустите контейнер:

```bash
docker compose down
docker compose up -d
```

### 4. Проверьте, что порт 25 слушается:

```bash
netstat -tuln | grep ":25 "
# Должно показать: tcp  0  0  0.0.0.0:25  LISTEN
```

### 5. Проверьте логи:

```bash
docker logs tempmail-app | grep -i smtp
# Должно быть: "SMTP-сервер запущен на порту 25"
```

## Важно:

- Порт 25 требует root прав
- Многие хостинги блокируют порт 25
- Проверьте, что провайдер разрешает использовать порт 25

## После настройки:

Письма из интернета (Яндекс, Gmail и т.д.) должны начать приходить на ваш SMTP сервер.

