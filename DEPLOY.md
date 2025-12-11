# Инструкция по деплою

## Быстрый деплой

### 1. Подготовка сервера

```bash
# Обновление системы
apt update && apt upgrade -y

# Установка Docker и Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Установка Docker Compose
apt install docker-compose-plugin -y
```

### 2. Клонирование репозитория

```bash
cd /var/www
git clone <repository-url> tempserver
cd tempserver/server
```

### 3. Настройка конфигурации

Отредактируйте `docker-compose.yml`:

```yaml
environment:
  - DB_PASSWORD=ваш-надежный-пароль  # Измените пароль БД!
  - MAIL_DOMAIN=vsebeauty.ru          # Ваш домен
  - SMTP_PORT=25                      # Стандартный SMTP порт
```

### 4. Настройка DNS

В панели управления доменом настройте:

- **MX запись**: `vsebeauty.ru` → `mail.vsebeauty.ru` (приоритет 10)
- **A запись**: `mail.vsebeauty.ru` → IP адрес вашего сервера
- **SPF запись**: `v=spf1 ip4:ВАШ-IP ~all`

### 5. Настройка файрвола

```bash
# Откройте необходимые порты
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw allow 25/tcp    # SMTP
ufw allow 8080/tcp  # API (или через reverse proxy)

# Включите файрвол
ufw enable
```

### 6. Запуск

```bash
# Запустите все сервисы
docker compose up -d

# Проверьте статус
docker compose ps

# Проверьте логи
docker logs tempmail-app
```

### 7. Проверка работы

```bash
# Проверьте SMTP сервер
docker logs tempmail-app | grep -i smtp
# Должно быть: "SMTP-сервер запущен на порту 25"

# Проверьте порт
netstat -tuln | grep ":25 "

# Проверьте API
curl http://localhost:8080/health
```

## Обновление

```bash
# Остановите контейнеры
docker compose down

# Обновите код
git pull

# Пересоберите и запустите
docker compose up -d --build
```

## Мониторинг

```bash
# Логи приложения
docker logs -f tempmail-app

# Логи базы данных
docker logs -f tempmail-postgres

# Статистика контейнеров
docker stats

# Использование диска
docker system df
```

## Резервное копирование

```bash
# Бэкап базы данных
docker exec tempmail-postgres pg_dump -U postgres tempmail > backup.sql

# Восстановление
docker exec -i tempmail-postgres psql -U postgres tempmail < backup.sql
```

## Troubleshooting

### Порт 25 не открывается

Проверьте:
1. Файрвол: `ufw status`
2. Провайдер может блокировать порт 25 - обратитесь в поддержку
3. Используйте `cap_add: - NET_BIND_SERVICE` в docker-compose.yml

### Письма не приходят

1. Проверьте DNS: `dig MX vsebeauty.ru`
2. Проверьте логи: `docker logs tempmail-app | grep -i smtp`
3. Проверьте порт: `netstat -tuln | grep :25`

### Проблемы с производительностью

```bash
# Увеличьте лимиты в docker-compose.yml
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 2G
```

