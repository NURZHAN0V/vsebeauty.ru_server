package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"tempmail/internal/domain"
)

// MailboxRepository — репозиторий для работы с почтовыми ящиками
type MailboxRepository struct {
	db *sql.DB // Подключение к базе данных
}

// NewMailboxRepository создаёт новый репозиторий
func NewMailboxRepository(db *sql.DB) *MailboxRepository {
	return &MailboxRepository{db: db}
}

// Create создаёт новый почтовый ящик
func (r *MailboxRepository) Create(address string, ttl time.Duration) (*domain.Mailbox, error) {
	// Генерируем уникальный ID
	id := uuid.New().String()

	// Вычисляем время истечения
	now := time.Now()
	expiresAt := now.Add(ttl)

	// SQL-запрос для вставки записи
	// $1, $2, $3, $4 — это плейсхолдеры для параметров
	// Они защищают от SQL-инъекций
	query := `
        INSERT INTO mailboxes (id, address, created_at, expires_at, is_active)
        VALUES ($1, $2, $3, $4, $5)
    `

	// Выполняем запрос
	// Exec используется для запросов, которые не возвращают данные (INSERT, UPDATE, DELETE)
	_, err := r.db.Exec(query, id, address, now, expiresAt, true)
	if err != nil {
		return nil, err
	}

	// Возвращаем созданный ящик
	return &domain.Mailbox{
		ID:        id,
		Address:   address,
		CreatedAt: now,
		ExpiresAt: expiresAt,
		IsActive:  true,
	}, nil
}

// GetByID находит ящик по ID
func (r *MailboxRepository) GetByID(id string) (*domain.Mailbox, error) {
	// SQL-запрос для выборки одной записи
	query := `
        SELECT id, address, created_at, expires_at, is_active
        FROM mailboxes
        WHERE id = $1
    `

	// Создаём пустую структуру для результата
	mailbox := &domain.Mailbox{}

	// QueryRow выполняет запрос и возвращает одну строку
	// Scan читает значения из строки в поля структуры
	err := r.db.QueryRow(query, id).Scan(
		&mailbox.ID,
		&mailbox.Address,
		&mailbox.CreatedAt,
		&mailbox.ExpiresAt,
		&mailbox.IsActive,
	)

	// Проверяем ошибки
	if err == sql.ErrNoRows {
		// Специальная ошибка — запись не найдена
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mailbox, nil
}

// GetByAddress находит ящик по email-адресу
func (r *MailboxRepository) GetByAddress(address string) (*domain.Mailbox, error) {
	query := `
        SELECT id, address, created_at, expires_at, is_active
        FROM mailboxes
        WHERE address = $1 AND is_active = true
    `

	mailbox := &domain.Mailbox{}
	err := r.db.QueryRow(query, address).Scan(
		&mailbox.ID,
		&mailbox.Address,
		&mailbox.CreatedAt,
		&mailbox.ExpiresAt,
		&mailbox.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mailbox, nil
}

// Delete удаляет почтовый ящик
func (r *MailboxRepository) Delete(id string) error {
	query := `DELETE FROM mailboxes WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// DeleteExpired удаляет все истёкшие ящики
func (r *MailboxRepository) DeleteExpired() (int64, error) {
	query := `DELETE FROM mailboxes WHERE expires_at < NOW()`

	// Exec возвращает Result, из которого можно узнать количество затронутых строк
	result, err := r.db.Exec(query)
	if err != nil {
		return 0, err
	}

	// RowsAffected возвращает количество удалённых записей
	return result.RowsAffected()
}
