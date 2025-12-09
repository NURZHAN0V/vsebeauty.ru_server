package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config — главная структура конфигурации приложения
// Все поля заполняются из переменных окружения
type Config struct {
	Server   ServerConfig   // Настройки серверов
	Database DatabaseConfig // Настройки базы данных
	Redis    RedisConfig    // Настройки Redis
	Mail     MailConfig     // Настройки почты
	Limits   LimitsConfig   // Лимиты
}

// ServerConfig — настройки HTTP и SMTP серверов
type ServerConfig struct {
	HTTPPort int `envconfig:"HTTP_PORT" default:"8080"` // Порт HTTP сервера
	SMTPPort int `envconfig:"SMTP_PORT" default:"2525"` // Порт SMTP сервера
}

// DatabaseConfig — настройки подключения к PostgreSQL
type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"` // Адрес сервера БД
	Port     int    `envconfig:"DB_PORT" default:"5432"`      // Порт БД
	Name     string `envconfig:"DB_NAME" default:"tempmail"`  // Имя базы данных
	User     string `envconfig:"DB_USER" default:"postgres"`  // Пользователь БД
	Password string `envconfig:"DB_PASSWORD" required:"true"` // Пароль БД (обязательный)
}

// RedisConfig — настройки подключения к Redis
type RedisConfig struct {
	Host string `envconfig:"REDIS_HOST" default:"localhost"` // Адрес Redis
	Port int    `envconfig:"REDIS_PORT" default:"6379"`      // Порт Redis
}

// MailConfig — настройки почтовых ящиков
type MailConfig struct {
	Domain     string        `envconfig:"MAIL_DOMAIN" default:"tempmail.dev"` // Домен для email
	DefaultTTL time.Duration `envconfig:"DEFAULT_TTL" default:"1h"`           // Время жизни по умолчанию
	MaxTTL     time.Duration `envconfig:"MAX_TTL" default:"24h"`              // Максимальное время жизни
}

// LimitsConfig — лимиты и ограничения
type LimitsConfig struct {
	MaxMessageSize        int `envconfig:"MAX_MESSAGE_SIZE" default:"10485760"`    // Макс. размер письма (10 MB)
	MaxAttachmentSize     int `envconfig:"MAX_ATTACHMENT_SIZE" default:"5242880"`  // Макс. размер вложения (5 MB)
	MaxMessagesPerMailbox int `envconfig:"MAX_MESSAGES_PER_MAILBOX" default:"100"` // Макс. писем в ящике
}

// Load загружает конфигурацию из переменных окружения
// Сначала пытается прочитать файл .env, затем читает переменные окружения
func Load() (*Config, error) {
	// Пытаемся загрузить .env файл
	// Если файла нет — не страшно, будем читать из системных переменных
	_ = godotenv.Load()

	// Создаём пустую структуру конфигурации
	var cfg Config

	// Заполняем структуру из переменных окружения
	// Если обязательное поле отсутствует — вернётся ошибка
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	// Возвращаем указатель на конфигурацию
	return &cfg, nil
}
