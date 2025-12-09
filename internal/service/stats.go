package service

import (
	"sync"
	"time"
)

// Stats хранит статистику работы сервиса
type Stats struct {
	mu                sync.RWMutex // Мьютекс для безопасного доступа
	TotalMailboxes    int64        // Всего создано ящиков
	TotalMessages     int64        // Всего получено писем
	TotalSpamMessages int64        // Всего спам-писем
	DeletedMailboxes  int64        // Удалено ящиков
	LastCleanup       time.Time    // Время последней очистки
}

// GlobalStats — глобальная статистика
var GlobalStats = &Stats{}

// IncrementMailboxes увеличивает счётчик ящиков
func (s *Stats) IncrementMailboxes() {
	s.mu.Lock()         // Блокируем для записи
	defer s.mu.Unlock() // Разблокируем при выходе
	s.TotalMailboxes++
}

// IncrementMessages увеличивает счётчик писем
func (s *Stats) IncrementMessages(isSpam bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalMessages++
	if isSpam {
		s.TotalSpamMessages++
	}
}

// AddDeletedMailboxes добавляет к счётчику удалённых ящиков
func (s *Stats) AddDeletedMailboxes(count int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DeletedMailboxes += count
	s.LastCleanup = time.Now()
}

// GetStats возвращает копию статистики
func (s *Stats) GetStats() Stats {
	s.mu.RLock()         // Блокируем для чтения
	defer s.mu.RUnlock() // Разблокируем при выходе
	return Stats{
		TotalMailboxes:    s.TotalMailboxes,
		TotalMessages:     s.TotalMessages,
		TotalSpamMessages: s.TotalSpamMessages,
		DeletedMailboxes:  s.DeletedMailboxes,
		LastCleanup:       s.LastCleanup,
	}
}

