package domain

import (
	"time"
)

// Message — входящее письмо
type Message struct {
	ID          string    `json:"id"`           // Уникальный идентификатор
	MailboxID   string    `json:"mailbox_id"`   // ID почтового ящика
	FromAddress string    `json:"from_address"` // Адрес отправителя
	Subject     string    `json:"subject"`      // Тема письма
	BodyText    string    `json:"body_text"`    // Текстовое содержимое
	BodyHTML    string    `json:"body_html"`    // HTML содержимое
	ReceivedAt  time.Time `json:"received_at"`  // Дата получения
	IsRead      bool      `json:"is_read"`      // Прочитано ли
	IsSpam      bool      `json:"is_spam"`      // Помечено как спам
}

// Attachment — вложение к письму
type Attachment struct {
	ID          string `json:"id"`           // Уникальный идентификатор
	MessageID   string `json:"message_id"`   // ID письма
	Filename    string `json:"filename"`     // Имя файла
	ContentType string `json:"content_type"` // MIME-тип (например, image/png)
	SizeBytes   int64  `json:"size_bytes"`   // Размер в байтах
	StoragePath string `json:"storage_path"` // Путь к файлу на диске
}
