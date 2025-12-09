# Этап 1: Сборка приложения
# Используем официальный образ Go для сборки
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
# Это отдельный шаг для кэширования (если зависимости не изменились, этот шаг пропустится)
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
# CGO_ENABLED=0 — отключаем CGO для статической сборки
# -ldflags="-s -w" — уменьшаем размер бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /tempmail ./cmd/api

# Этап 2: Финальный образ
# Используем минимальный образ Alpine
FROM alpine:3.19

# Устанавливаем ca-certificates для HTTPS и tzdata для часовых поясов
RUN apk --no-cache add ca-certificates tzdata

# Создаём пользователя для безопасности (не запускаем от root)
RUN adduser -D -g '' appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение из этапа сборки
COPY --from=builder /tempmail .

# Копируем миграции
COPY migrations ./migrations

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порты
EXPOSE 8080 2525

# Команда запуска
CMD ["./tempmail"]