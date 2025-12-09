package smtp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"

	"github.com/emersion/go-smtp"

	"tempmail/internal/domain"
)

// Session обрабатывает одну SMTP-сессию (одно письмо)
type Session struct {
	backend *Backend // Ссылка на бэкенд
	from    string   // Адрес отправителя
	to      []string // Адреса получателей
}

// AuthPlain обрабатывает PLAIN-аутентификацию
// Мы не требуем аутентификацию, поэтому всегда возвращаем nil
func (s *Session) AuthPlain(username, password string) error {
	// Аутентификация не требуется для приёма писем
	return nil
}

// Mail вызывается, когда клиент сообщает адрес отправителя (MAIL FROM)
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Printf("MAIL FROM: %s", from)
	s.from = from
	return nil
}

// Rcpt вызывается для каждого получателя (RCPT TO)
// Здесь мы проверяем, существует ли почтовый ящик
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	log.Printf("RCPT TO: %s", to)

	// Извлекаем email из формата "Name <email@domain.com>"
	address := extractEmail(to)

	// Проверяем, что письмо для нашего домена
	if !strings.HasSuffix(address, "@"+s.backend.domain) {
		return fmt.Errorf("мы не принимаем письма для домена %s", address)
	}

	// Проверяем, существует ли ящик
	mailbox, err := s.backend.mailboxService.GetByAddress(address)
	if err != nil {
		log.Printf("Ошибка проверки ящика: %v", err)
		return &smtp.SMTPError{
			Code:    550,
			Message: "Почтовый ящик не найден",
		}
	}
	if mailbox == nil {
		return &smtp.SMTPError{
			Code:    550,
			Message: "Почтовый ящик не существует",
		}
	}

	// Добавляем получателя
	s.to = append(s.to, address)
	return nil
}

// Data вызывается, когда клиент отправляет содержимое письма
func (s *Session) Data(r io.Reader) error {
	log.Println("Получение данных письма...")

	// Читаем всё письмо в буфер
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return err
	}

	// Парсим письмо
	msg, err := mail.ReadMessage(&buf)
	if err != nil {
		log.Printf("Ошибка парсинга письма: %v", err)
		return err
	}

	// Извлекаем заголовки
	subject := decodeHeader(msg.Header.Get("Subject"))
	from := msg.Header.Get("From")
	contentType := msg.Header.Get("Content-Type")

	if from == "" {
		from = s.from
	}

	// Парсим тело письма
	bodyText, bodyHTML := parseBody(msg.Body, contentType)

	log.Printf("Письмо от %s, тема: %s", from, subject)

	// Сохраняем письмо для каждого получателя
	for _, to := range s.to {
		err := s.saveMessage(to, from, subject, bodyText, bodyHTML)
		if err != nil {
			log.Printf("Ошибка сохранения письма для %s: %v", to, err)
		}
	}

	return nil
}

// saveMessage сохраняет письмо в базу данных
func (s *Session) saveMessage(to, from, subject, bodyText, bodyHTML string) error {
	mailbox, err := s.backend.mailboxService.GetByAddress(to)
	if err != nil {
		return err
	}
	if mailbox == nil {
		return fmt.Errorf("ящик %s не найден", to)
	}

	message := &domain.Message{
		MailboxID:   mailbox.ID,
		FromAddress: extractEmail(from),
		Subject:     subject,
		BodyText:    bodyText,
		BodyHTML:    bodyHTML,
		IsRead:      false,
		IsSpam:      false,
	}

	return s.backend.messageService.Create(message)
}

// parseBody парсит тело письма и извлекает текст и HTML
func parseBody(body io.Reader, contentType string) (text, html string) {
	// Если Content-Type не указан, считаем plain text
	if contentType == "" {
		data, _ := io.ReadAll(body)
		return string(data), ""
	}

	// Парсим Content-Type
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		data, _ := io.ReadAll(body)
		return string(data), ""
	}

	// Если это multipart (письмо с несколькими частями)
	if strings.HasPrefix(mediaType, "multipart/") {
		boundary := params["boundary"]
		if boundary == "" {
			data, _ := io.ReadAll(body)
			return string(data), ""
		}

		// Читаем все части
		mr := multipart.NewReader(body, boundary)
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}

			partType := part.Header.Get("Content-Type")
			partData, _ := io.ReadAll(part)

			if strings.HasPrefix(partType, "text/plain") {
				text = string(partData)
			} else if strings.HasPrefix(partType, "text/html") {
				html = string(partData)
			}
		}
		return text, html
	}

	// Простое письмо (не multipart)
	data, _ := io.ReadAll(body)
	if strings.HasPrefix(mediaType, "text/html") {
		return "", string(data)
	}
	return string(data), ""
}

// decodeHeader декодирует заголовок письма (поддержка UTF-8)
func decodeHeader(s string) string {
	// Декодируем MIME-encoded слова (=?UTF-8?B?...?=)
	dec := new(mime.WordDecoder)
	decoded, err := dec.DecodeHeader(s)
	if err != nil {
		return s
	}
	return decoded
}

// Reset вызывается для сброса сессии
func (s *Session) Reset() {
	s.from = ""
	s.to = nil
}

// Logout вызывается при завершении сессии
func (s *Session) Logout() error {
	log.Println("SMTP-сессия завершена")
	return nil
}

// extractEmail извлекает email из строки вида "Name <email@domain.com>"
func extractEmail(s string) string {
	// Если есть угловые скобки, извлекаем email из них
	if start := strings.Index(s, "<"); start != -1 {
		if end := strings.Index(s, ">"); end != -1 {
			return strings.TrimSpace(s[start+1 : end])
		}
	}
	// Иначе возвращаем как есть
	return strings.TrimSpace(s)
}
