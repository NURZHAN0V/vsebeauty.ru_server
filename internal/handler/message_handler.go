package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"

	"tempmail/internal/service"
)

// MessageHandler — обработчик запросов для писем
type MessageHandler struct {
	service *service.MessageService
}

// NewMessageHandler создаёт новый обработчик
func NewMessageHandler(svc *service.MessageService) *MessageHandler {
	return &MessageHandler{service: svc}
}

// MessageResponse — структура ответа с данными письма
type MessageResponse struct {
	ID          string `json:"id"`
	MailboxID   string `json:"mailbox_id"`
	FromAddress string `json:"from_address"`
	Subject     string `json:"subject"`
	BodyText    string `json:"body_text,omitempty"`
	BodyHTML    string `json:"body_html,omitempty"`
	ReceivedAt  string `json:"received_at"`
	IsRead      bool   `json:"is_read"`
	IsSpam      bool   `json:"is_spam"`
}

// MessageListResponse — краткая информация о письме для списка
type MessageListResponse struct {
	ID          string `json:"id"`
	FromAddress string `json:"from_address"`
	Subject     string `json:"subject"`
	ReceivedAt  string `json:"received_at"`
	IsRead      bool   `json:"is_read"`
	IsSpam      bool   `json:"is_spam"`
}

// GetMessages возвращает список писем
// @Summary Получить список писем
// @Description Возвращает список всех писем в почтовом ящике (без содержимого)
// @Tags messages
// @Produce json
// @Param id path string true "ID почтового ящика" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200 {array} MessageListResponse "Список писем"
// @Failure 404 {object} ErrorResponse "Ящик не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /mailbox/{id}/messages [get]
func (h *MessageHandler) GetMessages(c *fiber.Ctx) error {
	mailboxID := c.Params("id")

	messages, err := h.service.GetByMailboxID(mailboxID)
	if err != nil {
		if errors.Is(err, service.ErrMailboxNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Почтовый ящик не найден",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Внутренняя ошибка сервера",
		})
	}

	// Преобразуем в формат ответа
	response := make([]MessageListResponse, len(messages))
	for i, msg := range messages {
		response[i] = MessageListResponse{
			ID:          msg.ID,
			FromAddress: msg.FromAddress,
			Subject:     msg.Subject,
			ReceivedAt:  msg.ReceivedAt.Format(time.RFC3339),
			IsRead:      msg.IsRead,
			IsSpam:      msg.IsSpam,
		}
	}

	return c.JSON(response)
}

// GetMessage возвращает письмо по ID
// @Summary Получить письмо
// @Description Возвращает полную информацию о письме включая содержимое. Автоматически помечает письмо как прочитанное.
// @Tags messages
// @Produce json
// @Param id path string true "ID почтового ящика" example("550e8400-e29b-41d4-a716-446655440000")
// @Param mid path string true "ID письма" example("550e8400-e29b-41d4-a716-446655440001")
// @Success 200 {object} MessageResponse "Информация о письме"
// @Failure 404 {object} ErrorResponse "Письмо не найдено"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /mailbox/{id}/messages/{mid} [get]
func (h *MessageHandler) GetMessage(c *fiber.Ctx) error {
	// mid — ID письма
	messageID := c.Params("mid")

	msg, err := h.service.GetByID(messageID)
	if err != nil {
		if errors.Is(err, service.ErrMessageNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Письмо не найдено",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Внутренняя ошибка сервера",
		})
	}

	return c.JSON(MessageResponse{
		ID:          msg.ID,
		MailboxID:   msg.MailboxID,
		FromAddress: msg.FromAddress,
		Subject:     msg.Subject,
		BodyText:    msg.BodyText,
		BodyHTML:    msg.BodyHTML,
		ReceivedAt:  msg.ReceivedAt.Format(time.RFC3339),
		IsRead:      msg.IsRead,
		IsSpam:      msg.IsSpam,
	})
}

// DeleteMessage удаляет письмо
// @Summary Удалить письмо
// @Description Удаляет письмо из почтового ящика
// @Tags messages
// @Param id path string true "ID почтового ящика" example("550e8400-e29b-41d4-a716-446655440000")
// @Param mid path string true "ID письма" example("550e8400-e29b-41d4-a716-446655440001")
// @Success 204 "Письмо успешно удалено"
// @Failure 404 {object} ErrorResponse "Письмо не найдено"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /mailbox/{id}/messages/{mid} [delete]
func (h *MessageHandler) DeleteMessage(c *fiber.Ctx) error {
	messageID := c.Params("mid")

	err := h.service.Delete(messageID)
	if err != nil {
		if errors.Is(err, service.ErrMessageNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Письмо не найдено",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Внутренняя ошибка сервера",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
