package smtp

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-smtp"

	"tempmail/internal/config"
	"tempmail/internal/service"
)

// Server — SMTP-сервер для приёма писем
type Server struct {
	server  *smtp.Server
	backend *Backend
	config  config.ServerConfig
}

// NewServer создаёт новый SMTP-сервер
func NewServer(
	cfg config.ServerConfig,
	mailCfg config.MailConfig,
	mailboxService *service.MailboxService,
	messageService *service.MessageService,
) *Server {
	// Создаём бэкенд
	backend := NewBackend(mailboxService, messageService, mailCfg.Domain)

	// Создаём SMTP-сервер
	server := smtp.NewServer(backend)

	// Настраиваем параметры сервера
	server.Addr = fmt.Sprintf(":%d", cfg.SMTPPort) // Адрес для прослушивания
	server.Domain = mailCfg.Domain                 // Наш домен
	server.ReadTimeout = 30 * time.Second          // Таймаут чтения
	server.WriteTimeout = 30 * time.Second         // Таймаут записи
	server.MaxMessageBytes = 10 * 1024 * 1024      // Макс. размер письма (10 MB)
	server.MaxRecipients = 10                      // Макс. получателей
	server.AllowInsecureAuth = true                // Разрешаем без TLS (для разработки)

	return &Server{
		server:  server,
		backend: backend,
		config:  cfg,
	}
}

// Start запускает SMTP-сервер
func (s *Server) Start() error {
	log.Printf("SMTP-сервер запущен на порту %d", s.config.SMTPPort)
	log.Printf("Домен: %s", s.server.Domain)

	// ListenAndServe блокирует выполнение
	return s.server.ListenAndServe()
}

// Close останавливает SMTP-сервер
func (s *Server) Close() error {
	return s.server.Close()
}
