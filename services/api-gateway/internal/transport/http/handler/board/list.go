package board

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/metadata"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
)

// CreateBoard обрабатывает создание новой доски
func (h *BoardHandler) ListBoards(c *fiber.Ctx) error {
	// Извлекаем userID, который сохранил Middleware AuthRequired
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Парсим тело запроса в Proto-структуру
	var req boardv1.ListBoardsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request body"})
	}

	// Создаем контекст с Metadata для передачи userID в gRPC
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", userID)

	// Вызываем микросервис Board Service
	resp, err := h.boardClient.ListBoards(ctx, &req)
	if err != nil {
		// маппинг ошибок gRPC -> HTTP
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Возвращаем результат
	return c.Status(http.StatusCreated).JSON(resp)
}
