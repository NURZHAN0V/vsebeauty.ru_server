# Настройка проброса порта 25 → 2525 через iptables

## Проблема:
Docker не может привязаться к порту 25 даже с `cap_add: NET_BIND_SERVICE`.

## Решение:
Используем проброс порта 25 → 2525 через iptables на хосте.

## Настройка:

### 1. Измените docker-compose.yml:
- SMTP_PORT=2525 (внутри контейнера)
- Порт 2525:2525 проброшен

### 2. Настройте iptables проброс:

```bash
# Пробросить порт 25 на 2525
iptables -t nat -A PREROUTING -p tcp --dport 25 -j REDIRECT --to-port 2525

# Сохранить правила
iptables-save > /etc/iptables/rules.v4

# Или для Ubuntu/Debian
netfilter-persistent save
```

### 3. Перезапустите контейнер:

```bash
docker compose down
docker compose up -d
```

### 4. Проверьте:

```bash
# Порт 25 должен перенаправляться на 2525
netstat -tuln | grep ":25 "
netstat -tuln | grep ":2525"

# Проверьте логи
docker logs tempmail-app | grep -i smtp
```

## Альтернатива: network_mode: host

Если iptables не работает, можно использовать:

```yaml
network_mode: host
```

Но это менее безопасно и требует изменения всей конфигурации.

## После настройки:

Письма из интернета будут приходить на порт 25, который автоматически перенаправится на 2525 в контейнере.

