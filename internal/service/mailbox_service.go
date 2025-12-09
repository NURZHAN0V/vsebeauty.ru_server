package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"tempmail/internal/config"
	"tempmail/internal/domain"
	"tempmail/internal/repository"
)

// Ошибки сервиса
var (
	ErrMailboxNotFound = errors.New("почтовый ящик не найден")
	ErrMailboxExpired  = errors.New("срок действия ящика истёк")
	ErrInvalidTTL      = errors.New("недопустимое время жизни")
)

// MailboxService — сервис для работы с почтовыми ящиками
type MailboxService struct {
	repo   *repository.MailboxRepository // Репозиторий для работы с БД
	config config.MailConfig             // Настройки почты
}

// NewMailboxService создаёт новый сервис
func NewMailboxService(repo *repository.MailboxRepository, cfg config.MailConfig) *MailboxService {
	return &MailboxService{
		repo:   repo,
		config: cfg,
	}
}

// Create создаёт новый почтовый ящик
// Если address пустой — генерируется случайный
func (s *MailboxService) Create(address string, ttl time.Duration) (*domain.Mailbox, error) {
	// Если адрес не указан — генерируем случайный
	if address == "" {
		address = s.generateRandomAddress()
	} else {
		// Добавляем домен к адресу
		address = fmt.Sprintf("%s@%s", address, s.config.Domain)
	}

	// Проверяем TTL
	if ttl <= 0 {
		ttl = s.config.DefaultTTL
	}
	if ttl > s.config.MaxTTL {
		return nil, ErrInvalidTTL
	}

	// Проверяем, не занят ли адрес
	existing, err := s.repo.GetByAddress(address)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		// Адрес занят — генерируем новый
		address = s.generateRandomAddress()
	}

	// Создаём ящик
	return s.repo.Create(address, ttl)
}

// GetByID возвращает ящик по ID
func (s *MailboxService) GetByID(id string) (*domain.Mailbox, error) {
	mailbox, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if mailbox == nil {
		return nil, ErrMailboxNotFound
	}

	// Проверяем, не истёк ли срок
	if mailbox.IsExpired() {
		return nil, ErrMailboxExpired
	}

	return mailbox, nil
}

// Delete удаляет почтовый ящик
func (s *MailboxService) Delete(id string) error {
	// Проверяем существование
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// generateRandomAddress генерирует случайный email-адрес
func (s *MailboxService) generateRandomAddress() string {
	// Символы для генерации
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

	// Генерируем 10 случайных символов
	result := make([]byte, 10)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return fmt.Sprintf("%s@%s", string(result), s.config.Domain)
}

// init вызывается при загрузке пакета
// Инициализируем генератор случайных чисел
func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetByAddress возвращает ящик по email-адресу
func (s *MailboxService) GetByAddress(address string) (*domain.Mailbox, error) {
	mailbox, err := s.repo.GetByAddress(address)
	if err != nil {
		return nil, err
	}
	if mailbox == nil {
		return nil, nil // Ящик не найден — это не ошибка
	}

	// Проверяем, не истёк ли срок
	if mailbox.IsExpired() {
		return nil, nil // Истёкший ящик считаем несуществующим
	}

	return mailbox, nil
}
