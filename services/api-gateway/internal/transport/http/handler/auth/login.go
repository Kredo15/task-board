package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
)

// Login обрабатывает вход пользователя
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Парсим тело запроса в Proto-структуру
	var req authv1.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		// маппинг ошибок gRPC -> HTTP
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request body"})
	}

	// Создаем контекст с Metadata для передачи userID в gRPC
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Вызываем микросервис
	resp, err := h.authClient.Login(ctx, &req)
	if err != nil {
		// маппинг ошибок gRPC -> HTTP
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Возвращаем результат
	return c.Status(http.StatusAccepted).JSON(resp)
}
