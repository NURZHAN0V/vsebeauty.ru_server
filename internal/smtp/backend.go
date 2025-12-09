package smtp

import (
	"log"

	"github.com/emersion/go-smtp"

	"tempmail/internal/service"
)

// Backend реализует интерфейс smtp.Backend
// Он создаёт сессии для каждого входящего соединения
type Backend struct {
	mailboxService *service.MailboxService // Сервис для проверки ящиков
	messageService *service.MessageService // Сервис для сохранения писем
	domain         string                  // Наш домен (tempmail.dev)
}

// NewBackend создаёт новый SMTP-бэкенд
func NewBackend(
	mailboxService *service.MailboxService,
	messageService *service.MessageService,
	domain string,
) *Backend {
	return &Backend{
		mailboxService: mailboxService,
		messageService: messageService,
		domain:         domain,
	}
}

// NewSession создаёт новую сессию для входящего соединения
// Вызывается при каждом новом подключении к SMTP-серверу
func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	log.Printf("Новое SMTP-соединение от %s", c.Hostname())

	return &Session{
		backend: b,
	}, nil
}
