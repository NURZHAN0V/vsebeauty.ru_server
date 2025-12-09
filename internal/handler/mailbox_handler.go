package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"

	"tempmail/internal/service"
)

// MailboxHandler — обработчик запросов для почтовых ящиков
type MailboxHandler struct {
	service *service.MailboxService
}

// NewMailboxHandler создаёт новый обработчик
func NewMailboxHandler(svc *service.MailboxService) *MailboxHandler {
	return &MailboxHandler{service: svc}
}

// CreateRequest — структура запроса на создание ящика
type CreateRequest struct {
	Address string `json:"address"` // Желаемый адрес (необязательно)
	TTL     string `json:"ttl"`     // Время жизни (например, "1h", "30m")
}

// MailboxResponse — структура ответа с данными ящика
type MailboxResponse struct {
	ID        string `json:"id"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
	IsActive  bool   `json:"is_active"`
}

// Create создаёт новый почтовый ящик
// @Summary Создать почтовый ящик
// @Description Создаёт новый временный почтовый ящик. Если адрес не указан, генерируется случайный.
// @Tags mailbox
// @Accept json
// @Produce json
// @Param request body CreateRequest false "Параметры создания (необязательно)"
// @Success 201 {object} MailboxResponse "Ящик успешно создан"
// @Failure 400 {object} ErrorResponse "Неверные параметры запроса"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /mailbox [post]
func (h *MailboxHandler) Create(c *fiber.Ctx) error {
	// Парсим тело запроса
	var req CreateRequest

	// BodyParser читает JSON из тела запроса и заполняет структуру
	if err := c.BodyParser(&req); err != nil {
		// Если тело пустое — это нормально, создадим ящик с настройками по умолчанию
		req = CreateRequest{}
	}

	// Парсим TTL
	var ttl time.Duration
	if req.TTL != "" {
		var err error
		ttl, err = time.ParseDuration(req.TTL)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Неверный формат TTL. Используйте формат: 1h, 30m, 24h",
			})
		}
	}

	// Создаём ящик
	mailbox, err := h.service.Create(req.Address, ttl)
	if err != nil {
		// Проверяем тип ошибки
		if errors.Is(err, service.ErrInvalidTTL) {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "TTL превышает максимально допустимое значение",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Внутренняя ошибка сервера",
		})
	}

	// Возвращаем успешный ответ
	// Status(201) — код "Created" (создано)
	return c.Status(fiber.StatusCreated).JSON(MailboxResponse{
		ID:        mailbox.ID,
		Address:   mailbox.Address,
		CreatedAt: mailbox.CreatedAt.Format(time.RFC3339),
		ExpiresAt: mailbox.ExpiresAt.Format(time.RFC3339),
		IsActive:  mailbox.IsActive,
	})
}

// Get возвращает информацию о ящике
// @Summary Получить информацию о ящике
// @Description Возвращает информацию о почтовом ящике по его ID
// @Tags mailbox
// @Produce json
// @Param id path string true "ID почтового ящика" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200 {object} MailboxResponse "Информация о ящике"
// @Failure 404 {object} ErrorResponse "Ящик не найден"
// @Failure 410 {object} ErrorResponse "Срок действия ящика истёк"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /mailbox/{id} [get]
func (h *MailboxHandler) Get(c *fiber.Ctx) error {
	// Params получает параметр из URL
	id := c.Params("id")

	mailbox, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrMailboxNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Почтовый ящик не найден",
			})
		}
		if errors.Is(err, service.ErrMailboxExpired) {
			return c.Status(fiber.StatusGone).JSON(ErrorResponse{
				Error: "Срок действия ящика истёк",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Внутренняя ошибка сервера",
		})
	}

	return c.JSON(MailboxResponse{
		ID:        mailbox.ID,
		Address:   mailbox.Address,
		CreatedAt: mailbox.CreatedAt.Format(time.RFC3339),
		ExpiresAt: mailbox.ExpiresAt.Format(time.RFC3339),
		IsActive:  mailbox.IsActive,
	})
}

// Delete удаляет почтовый ящик
// @Summary Удалить почтовый ящик
// @Description Удаляет почтовый ящик и все связанные письма
// @Tags mailbox
// @Param id path string true "ID почтового ящика" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 204 "Ящик успешно удалён"
// @Failure 404 {object} ErrorResponse "Ящик не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /mailbox/{id} [delete]
func (h *MailboxHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.Delete(id)
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

	// 204 No Content — успешное удаление без тела ответа
	return c.SendStatus(fiber.StatusNoContent)
}
