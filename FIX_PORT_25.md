# Исправление ошибки "permission denied" на порту 25

## Проблема:
```
listen tcp :25: bind: permission denied
```

Порт 25 требует root прав, но Docker контейнер не может привязаться к нему.

## Решение:

Добавлено `cap_add: - NET_BIND_SERVICE` в docker-compose.yml

Это разрешает контейнеру привязываться к привилегированным портам (< 1024).

## Перезапустите контейнер:

```bash
docker compose down
docker compose up -d
```

## Проверьте:

```bash
# Логи должны показать успешный запуск
docker logs tempmail-app | grep -i smtp

# Порт должен слушаться
netstat -tuln | grep ":25 "
```

## Альтернативное решение (если не работает):

Если `cap_add` не помогает, можно использовать проброс порта через iptables:

```bash
# Пробросить порт 25 на 2525 внутри контейнера
iptables -t nat -A PREROUTING -p tcp --dport 25 -j REDIRECT --to-port 2525
```

И изменить в docker-compose.yml:
- SMTP_PORT=2525 (внутри контейнера)
- Порт 25:2525 (проброс)

Но лучше использовать cap_add - это правильный способ.

