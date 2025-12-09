package main

import (
	"fmt"
	"log"

	"tempmail/internal/config"
	"tempmail/internal/repository"
	"tempmail/internal/service"
	smtpserver "tempmail/internal/smtp"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	fmt.Println("=== TempMail SMTP Server ===")

	// Подключаемся к базе данных
	fmt.Println("Подключение к PostgreSQL...")
	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()
	fmt.Println("Подключение успешно!")

	// Создаём репозитории
	mailboxRepo := repository.NewMailboxRepository(db.DB)
	messageRepo := repository.NewMessageRepository(db.DB)

	// Создаём сервисы
	mailboxService := service.NewMailboxService(mailboxRepo, cfg.Mail)
	messageService := service.NewMessageService(messageRepo, mailboxRepo, cfg.Limits)

	// Создаём и запускаем SMTP-сервер
	server := smtpserver.NewServer(cfg.Server, cfg.Mail, mailboxService, messageService)

	fmt.Printf("\nSMTP-сервер запущен на порту %d\n", cfg.Server.SMTPPort)
	fmt.Printf("Домен: %s\n", cfg.Mail.Domain)
	fmt.Println("Нажмите Ctrl+C для остановки")

	if err := server.Start(); err != nil {
		log.Fatal("Ошибка SMTP-сервера:", err)
	}
}
