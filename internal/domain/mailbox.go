package domain

import (
	"time"
)

// Mailbox — почтовый ящик
// Каждый ящик имеет уникальный адрес и время жизни
type Mailbox struct {
	ID        string    `json:"id"`         // Уникальный идентификатор (UUID)
	Address   string    `json:"address"`    // Email адрес (например, abc123@tempmail.dev)
	CreatedAt time.Time `json:"created_at"` // Дата создания
	ExpiresAt time.Time `json:"expires_at"` // Дата истечения срока
	IsActive  bool      `json:"is_active"`  // Активен ли ящик
}

// IsExpired проверяет, истёк ли срок действия ящика
func (m *Mailbox) IsExpired() bool {
	// time.Now() возвращает текущее время
	// After проверяет, наступило ли время ExpiresAt
	return time.Now().After(m.ExpiresAt)
}
