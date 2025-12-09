package repository

import (
	"database/sql"
	"fmt"

	// Импортируем драйвер PostgreSQL
	_ "github.com/lib/pq"

	"tempmail/internal/config"
)

// PostgresDB — обёртка над подключением к PostgreSQL
type PostgresDB struct {
	DB *sql.DB // Стандартный интерфейс Go для работы с БД
}

// NewPostgresDB создаёт новое подключение к PostgreSQL
func NewPostgresDB(cfg config.DatabaseConfig) (*PostgresDB, error) {
	// Формируем строку подключения
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
	fmt.Printf("DEBUG connStr: postgres://%s:***@%s:%d/%s?sslmode=disable\n", cfg.User, cfg.Host, cfg.Port, cfg.Name)

	// Открываем соединение с базой данных
	// sql.Open не устанавливает соединение сразу, только проверяет параметры
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Проверяем, что соединение работает
	// Ping отправляет запрос к БД и ждёт ответа
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Возвращаем обёртку с подключением
	return &PostgresDB{DB: db}, nil
}

// Close закрывает соединение с базой данных
func (p *PostgresDB) Close() error {
	return p.DB.Close()
}
