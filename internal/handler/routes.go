package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"tempmail/internal/service"
)

// SetupRoutes настраивает все маршруты приложения
func SetupRoutes(
	app *fiber.App,
	mailboxHandler *MailboxHandler,
	messageHandler *MessageHandler,
) {
	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API v1
	api := app.Group("/api/v1")

	// Mailbox routes
	mailbox := api.Group("/mailbox")
	mailbox.Post("/", mailboxHandler.Create)
	mailbox.Get("/:id", mailboxHandler.Get)
	mailbox.Delete("/:id", mailboxHandler.Delete)

	// Message routes
	mailbox.Get("/:id/messages", messageHandler.GetMessages)
	mailbox.Get("/:id/messages/:mid", messageHandler.GetMessage)
	mailbox.Delete("/:id/messages/:mid", messageHandler.DeleteMessage)

	// Health check
	// @Summary Проверка здоровья
	// @Description Возвращает статус сервера
	// @Tags system
	// @Produce json
	// @Success 200 {object} map[string]string "Статус сервера"
	// @Router /health [get]
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Stats
	// @Summary Статистика сервиса
	// @Description Возвращает статистику работы сервиса
	// @Tags system
	// @Produce json
	// @Success 200 {object} map[string]interface{} "Статистика"
	// @Router /stats [get]
	app.Get("/stats", func(c *fiber.Ctx) error {
		stats := service.GlobalStats.GetStats()
		return c.JSON(fiber.Map{
			"total_mailboxes":   stats.TotalMailboxes,
			"total_messages":    stats.TotalMessages,
			"total_spam":        stats.TotalSpamMessages,
			"deleted_mailboxes": stats.DeletedMailboxes,
			"last_cleanup":      stats.LastCleanup.Format("2006-01-02 15:04:05"),
		})
	})
}
