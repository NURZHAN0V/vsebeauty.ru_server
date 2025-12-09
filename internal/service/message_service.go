package service

import (
	"errors"

	"tempmail/internal/config"
	"tempmail/internal/domain"
	"tempmail/internal/repository"
)

// Ошибки сервиса
var (
	ErrMessageNotFound = errors.New("письмо не найдено")
	ErrMailboxFull     = errors.New("ящик переполнен")
	ErrMessageTooLarge = errors.New("письмо слишком большое")
)

// MessageService — сервис для работы с письмами
type MessageService struct {
	msgRepo     *repository.MessageRepository
	mailboxRepo *repository.MailboxRepository
	limits      config.LimitsConfig
}

// NewMessageService создаёт новый сервис
func NewMessageService(
	msgRepo *repository.MessageRepository,
	mailboxRepo *repository.MailboxRepository,
	limits config.LimitsConfig,
) *MessageService {
	return &MessageService{
		msgRepo:     msgRepo,
		mailboxRepo: mailboxRepo,
		limits:      limits,
	}
}

// Create создаёт новое письмо
func (s *MessageService) Create(msg *domain.Message) error {
	// Проверяем существование ящика
	mailbox, err := s.mailboxRepo.GetByID(msg.MailboxID)
	if err != nil {
		return err
	}
	if mailbox == nil {
		return ErrMailboxNotFound
	}

	// Проверяем, не истёк ли срок ящика
	if mailbox.IsExpired() {
		return ErrMailboxExpired
	}

	// Проверяем количество писем в ящике
	count, err := s.msgRepo.CountByMailboxID(msg.MailboxID)
	if err != nil {
		return err
	}
	if count >= s.limits.MaxMessagesPerMailbox {
		return ErrMailboxFull
	}

	// Проверяем размер письма
	messageSize := len(msg.BodyText) + len(msg.BodyHTML)
	if messageSize > s.limits.MaxMessageSize {
		return ErrMessageTooLarge
	}

	return s.msgRepo.Create(msg)
}

// GetByMailboxID возвращает все письма ящика
func (s *MessageService) GetByMailboxID(mailboxID string) ([]*domain.Message, error) {
	// Проверяем существование ящика
	mailbox, err := s.mailboxRepo.GetByID(mailboxID)
	if err != nil {
		return nil, err
	}
	if mailbox == nil {
		return nil, ErrMailboxNotFound
	}

	return s.msgRepo.GetByMailboxID(mailboxID)
}

// GetByID возвращает письмо по ID
func (s *MessageService) GetByID(id string) (*domain.Message, error) {
	msg, err := s.msgRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, ErrMessageNotFound
	}

	// Помечаем как прочитанное
	_ = s.msgRepo.MarkAsRead(id)

	return msg, nil
}

// Delete удаляет письмо
func (s *MessageService) Delete(id string) error {
	msg, err := s.msgRepo.GetByID(id)
	if err != nil {
		return err
	}
	if msg == nil {
		return ErrMessageNotFound
	}

	return s.msgRepo.Delete(id)
}
