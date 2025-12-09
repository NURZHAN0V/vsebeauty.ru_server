package main

// @title TempMail API
// @version 1.0
// @description Сервис временных email-адресов с REST API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@tempmail.dev

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @schemes http https

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"tempmail/internal/config"
	"tempmail/internal/handler"
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

	fmt.Println("=== TempMail Server ===")

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

	// Создаём обработчики
	mailboxHandler := handler.NewMailboxHandler(mailboxService)
	messageHandler := handler.NewMessageHandler(messageService)

	// Создаём Fiber-приложение
	app := fiber.New(fiber.Config{
		AppName: "TempMail API",
	})

	// Настраиваем маршруты
	handler.SetupRoutes(app, mailboxHandler, messageHandler)

	// Создаём SMTP-сервер
	smtpServer := smtpserver.NewServer(cfg.Server, cfg.Mail, mailboxService, messageService)

	// Запускаем SMTP-сервер в отдельной горутине
	go func() {
		if err := smtpServer.Start(); err != nil {
			log.Printf("SMTP-сервер остановлен: %v", err)
		}
	}()

	// Запускаем HTTP-сервер в отдельной горутине
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
		if err := app.Listen(addr); err != nil {
			log.Printf("HTTP-сервер остановлен: %v", err)
		}
	}()

	fmt.Printf("\nHTTP API: http://localhost:%d\n", cfg.Server.HTTPPort)
	fmt.Printf("SMTP: localhost:%d\n", cfg.Server.SMTPPort)
	fmt.Println("\nНажмите Ctrl+C для остановки")

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nОстановка серверов...")
	smtpServer.Close()
	app.Shutdown()
}
