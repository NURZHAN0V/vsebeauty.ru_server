package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Использование: go run test-smtp.go <smtp-host> <smtp-port> <email-to>")
		fmt.Println("Пример: go run test-smtp.go localhost 2525 test@tempmail.dev")
		os.Exit(1)
	}

	host := os.Args[1]
	port := os.Args[2]
	to := os.Args[3]

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Подключение к SMTP серверу %s...\n", addr)

	// Создаём клиент
	client, err := smtp.Dial(addr)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer client.Close()

	fmt.Println("✓ Подключение успешно!")

	// Устанавливаем отправителя
	if err := client.Mail("test@example.com"); err != nil {
		log.Fatalf("Ошибка MAIL FROM: %v", err)
	}
	fmt.Println("✓ MAIL FROM установлен")

	// Устанавливаем получателя
	if err := client.Rcpt(to); err != nil {
		log.Fatalf("Ошибка RCPT TO: %v", err)
	}
	fmt.Printf("✓ RCPT TO установлен для %s\n", to)

	// Отправляем данные
	wc, err := client.Data()
	if err != nil {
		log.Fatalf("Ошибка DATA: %v", err)
	}

	message := fmt.Sprintf("From: test@example.com\r\nTo: %s\r\nSubject: Test Message\r\n\r\nЭто тестовое письмо от %s\r\n", to, time.Now().Format(time.RFC3339))
	_, err = wc.Write([]byte(message))
	if err != nil {
		log.Fatalf("Ошибка записи: %v", err)
	}

	err = wc.Close()
	if err != nil {
		log.Fatalf("Ошибка закрытия DATA: %v", err)
	}

	fmt.Println("✓ Письмо отправлено успешно!")
	fmt.Println("Проверьте базу данных или логи сервера")
}

