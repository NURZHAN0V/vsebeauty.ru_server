package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"tempmail/internal/domain"
)

// MessageRepository — репозиторий для работы с письмами
type MessageRepository struct {
	db *sql.DB
}

// NewMessageRepository создаёт новый репозиторий
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create создаёт новое письмо
func (r *MessageRepository) Create(msg *domain.Message) error {
	// Генерируем ID, если не задан
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}

	// Устанавливаем время получения
	if msg.ReceivedAt.IsZero() {
		msg.ReceivedAt = time.Now()
	}

	query := `
        INSERT INTO messages (id, mailbox_id, from_address, subject, body_text, body_html, received_at, is_read, is_spam)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	_, err := r.db.Exec(query,
		msg.ID,
		msg.MailboxID,
		msg.FromAddress,
		msg.Subject,
		msg.BodyText,
		msg.BodyHTML,
		msg.ReceivedAt,
		msg.IsRead,
		msg.IsSpam,
	)

	return err
}

// GetByMailboxID возвращает все письма для указанного ящика
func (r *MessageRepository) GetByMailboxID(mailboxID string) ([]*domain.Message, error) {
	query := `
        SELECT id, mailbox_id, from_address, subject, body_text, body_html, received_at, is_read, is_spam
        FROM messages
        WHERE mailbox_id = $1
        ORDER BY received_at DESC
    `

	// Query возвращает несколько строк
	rows, err := r.db.Query(query, mailboxID)
	if err != nil {
		return nil, err
	}
	// defer гарантирует, что rows.Close() выполнится при выходе из функции
	defer rows.Close()

	// Создаём срез для результатов
	var messages []*domain.Message

	// Перебираем все строки результата
	for rows.Next() {
		msg := &domain.Message{}
		err := rows.Scan(
			&msg.ID,
			&msg.MailboxID,
			&msg.FromAddress,
			&msg.Subject,
			&msg.BodyText,
			&msg.BodyHTML,
			&msg.ReceivedAt,
			&msg.IsRead,
			&msg.IsSpam,
		)
		if err != nil {
			return nil, err
		}
		// Добавляем письмо в срез
		messages = append(messages, msg)
	}

	// Проверяем ошибки, возникшие при переборе
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// GetByID находит письмо по ID
func (r *MessageRepository) GetByID(id string) (*domain.Message, error) {
	query := `
        SELECT id, mailbox_id, from_address, subject, body_text, body_html, received_at, is_read, is_spam
        FROM messages
        WHERE id = $1
    `

	msg := &domain.Message{}
	err := r.db.QueryRow(query, id).Scan(
		&msg.ID,
		&msg.MailboxID,
		&msg.FromAddress,
		&msg.Subject,
		&msg.BodyText,
		&msg.BodyHTML,
		&msg.ReceivedAt,
		&msg.IsRead,
		&msg.IsSpam,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// MarkAsRead помечает письмо как прочитанное
func (r *MessageRepository) MarkAsRead(id string) error {
	query := `UPDATE messages SET is_read = true WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Delete удаляет письмо
func (r *MessageRepository) Delete(id string) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// CountByMailboxID возвращает количество писем в ящике
func (r *MessageRepository) CountByMailboxID(mailboxID string) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE mailbox_id = $1`

	var count int
	err := r.db.QueryRow(query, mailboxID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
