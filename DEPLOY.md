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

**Рекомендуемый способ - через .env файл:**

```bash
# Создайте .env файл
nano .env
```

Добавьте следующие переменные в `.env`:

```bash
# Обязательные
DB_PASSWORD=ваш-надежный-пароль  # Измените пароль БД!
MAIL_DOMAIN=vsebeauty.ru          # Ваш домен
SMTP_PORT=25                      # Стандартный SMTP порт

# Сервер
HTTP_PORT=8080

# База данных
DB_HOST=postgres
DB_PORT=5432
DB_NAME=tempmail
DB_USER=postgres

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# Время жизни ящиков
DEFAULT_TTL=1h
MAX_TTL=24h

# Лимиты
MAX_MESSAGE_SIZE=10485760
MAX_ATTACHMENT_SIZE=5242880
MAX_MESSAGES_PER_MAILBOX=100
CLEANUP_INTERVAL=5m
```

**Важно:** 
- Файл `.env` не должен попадать в git (уже в .gitignore)
- Docker Compose автоматически читает переменные из `.env` через `env_file: - .env`
- Переменные из `.env` используются в `docker-compose.yml` через синтаксис `${VARIABLE:-default}`
- Если `.env` файл отсутствует, используются значения по умолчанию из `docker-compose.yml`
- **Обязательно измените `DB_PASSWORD`** на надежный пароль перед запуском в продакшене!

### 4. Настройка DNS

В панели управления доменом настройте:

- **MX запись**: `vsebeauty.ru` → `mail.vsebeauty.ru` (приоритет 10)
- **A запись**: `mail.vsebeauty.ru` → IP адрес вашего сервера
- **SPF запись**: `v=spf1 ip4:ВАШ-IP ~all`

### 5. Настройка Nginx (опционально, но рекомендуется)

Nginx используется как reverse proxy для API, чтобы не открывать порт 8080 напрямую.

**Установка Nginx:**
```bash
apt install nginx -y
```

**Настройка конфигурации:**
```bash
# Скопируйте конфигурацию
cp nginx/nginx-simple.conf /etc/nginx/sites-available/tempmail
ln -s /etc/nginx/sites-available/tempmail /etc/nginx/sites-enabled/

# Или для HTTPS используйте nginx.conf (после настройки SSL)
# cp nginx/nginx.conf /etc/nginx/sites-available/tempmail

# Проверьте конфигурацию
nginx -t

# Перезапустите Nginx
systemctl restart nginx
```

**Настройка SSL (Let's Encrypt):**
```bash
# Установите Certbot
apt install certbot python3-certbot-nginx -y

# Получите сертификат
certbot --nginx -d vsebeauty.ru -d www.vsebeauty.ru

# Автоматическое обновление
certbot renew --dry-run
```

### 6. Настройка файрвола

```bash
# Откройте необходимые порты
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw allow 25/tcp    # SMTP
# Порт 8080 можно закрыть, если используете Nginx

# Включите файрвол
ufw enable
```

### 7. Запуск

```bash
# Запустите все сервисы
docker compose up -d

# Проверьте статус
docker compose ps

# Проверьте логи
docker logs tempmail-app
```

### 8. Проверка работы

```bash
# Проверьте SMTP сервер
docker logs tempmail-app | grep -i smtp
# Должно быть: "SMTP-сервер запущен на порту 25"

# Проверьте порт
netstat -tuln | grep ":25 "

# Проверьте API (через Nginx, если настроен)
curl http://localhost/api/v1/health
# или напрямую
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

