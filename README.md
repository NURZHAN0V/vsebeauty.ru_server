# TempMail Server

Сервис временных email-адресов на Go с REST API и SMTP сервером.

## Быстрый старт

### Требования

- Docker и Docker Compose
- Go 1.23+ (для локальной разработки)

### Запуск через Docker Compose

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd server
```

2. Настройте переменные окружения:
```bash
cp .env.example .env
# Отредактируйте .env файл
```

3. Запустите сервисы:
```bash
docker compose up -d
```

4. Проверьте статус:
```bash
docker compose ps
docker logs tempmail-app
```

## Конфигурация

### Переменные окружения

Основные переменные (можно задать в `.env` или `docker-compose.yml`):

```bash
# Сервер
HTTP_PORT=8080          # Порт HTTP API
SMTP_PORT=25            # Порт SMTP сервера

# База данных
DB_HOST=postgres         # Хост PostgreSQL
DB_PORT=5432            # Порт PostgreSQL
DB_NAME=tempmail        # Имя базы данных
DB_USER=postgres        # Пользователь БД
DB_PASSWORD=secret      # Пароль БД (обязательно!)

# Почта
MAIL_DOMAIN=vsebeauty.ru  # Домен для email адресов
DEFAULT_TTL=1h          # Время жизни ящика по умолчанию
MAX_TTL=24h            # Максимальное время жизни

# Лимиты
MAX_MESSAGE_SIZE=10485760        # Макс. размер письма (10 MB)
MAX_ATTACHMENT_SIZE=5242880      # Макс. размер вложения (5 MB)
MAX_MESSAGES_PER_MAILBOX=100     # Макс. писем в ящике
```

## API Endpoints

### Почтовые ящики

- `POST /api/v1/mailbox` - Создать новый ящик
- `GET /api/v1/mailbox/:id` - Получить информацию о ящике
- `DELETE /api/v1/mailbox/:id` - Удалить ящик

### Письма

- `GET /api/v1/mailbox/:id/messages` - Получить список писем
- `GET /api/v1/mailbox/:id/messages/:mid` - Получить письмо
- `DELETE /api/v1/mailbox/:id/messages/:mid` - Удалить письмо

### Системные

- `GET /health` - Проверка здоровья сервера
- `GET /stats` - Статистика сервиса
- `GET /swagger/*` - Swagger UI документация

## Разработка

### Локальный запуск

```bash
# Установите зависимости
go mod download

# Запустите PostgreSQL через Docker
docker compose up -d postgres

# Запустите приложение
go run cmd/api/main.go
```

### Сборка

```bash
# Сборка бинарника
go build -o tempmail ./cmd/api

# Сборка Docker образа
docker build -t tempmail:latest .
```

## Деплой на продакшн

Подробная инструкция по деплою находится в [DEPLOY.md](./DEPLOY.md)

Кратко:
1. Настройте DNS записи (MX, A, SPF)
2. Откройте порт 25 в файрволе
3. Запустите через `docker compose up -d`
4. Проверьте логи и работу сервиса

## Структура проекта

```
server/
├── cmd/              # Точки входа приложения
│   ├── api/         # HTTP API сервер
│   └── smtp/        # SMTP сервер (отдельный)
├── internal/        # Внутренние пакеты
│   ├── config/      # Конфигурация
│   ├── domain/      # Модели данных
│   ├── handler/     # HTTP обработчики
│   ├── repository/  # Работа с БД
│   ├── service/     # Бизнес-логика
│   └── smtp/        # SMTP сервер
├── migrations/      # SQL миграции
├── scripts/         # Вспомогательные скрипты
├── docs/           # Документация
├── Dockerfile      # Docker образ
└── docker-compose.yml  # Docker Compose конфигурация
```

## Документация

Подробная документация находится в папке `docs/`:

- [Введение в Go](./docs/01-introduction.md)
- [Структура проекта](./docs/02-project-structure.md)
- [Работа с БД](./docs/03-database.md)
- [REST API](./docs/04-rest-api.md)
- [SMTP сервер](./docs/05-smtp-server.md)
- [Спам-фильтр](./docs/06-spam-filter.md)
- [Docker деплой](./docs/07-docker-deploy.md)
- [OpenAPI/Swagger](./docs/08-openapi.md)
- [Продакшн деплой](./docs/09-production-deploy.md)

## Troubleshooting

### SMTP сервер не запускается на порту 25

Ошибка: `bind: permission denied`

Решение: В `docker-compose.yml` добавлено `cap_add: - NET_BIND_SERVICE` для разрешения привязки к привилегированным портам.

### Письма не приходят из интернета

Проверьте:
1. Настроены ли DNS записи (MX, A, SPF)
2. Открыт ли порт 25 в файрволе
3. Не блокирует ли провайдер порт 25
4. Логи сервера на наличие входящих соединений

### Проблемы с подключением к БД

Проверьте:
1. Запущен ли контейнер PostgreSQL
2. Правильность пароля в переменных окружения
3. Логи контейнера: `docker logs tempmail-postgres`

## Лицензия

MIT
