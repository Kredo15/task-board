package http

import (
	"github.com/gofiber/fiber/v2"

	authHandler "github.com/Kredo15/task-board/services/api-gateway/internal/transport/http/handler/auth"
	boardHandler "github.com/Kredo15/task-board/services/api-gateway/internal/transport/http/handler/board"
	"github.com/Kredo15/task-board/services/api-gateway/internal/transport/http/middleware"
)

// Router содержит все зависимости для настройки маршрутов
type Router struct {
	authHandler  *authHandler.AuthHandler
	boardHandler *boardHandler.BoardHandler
	authMW       *middleware.JWTMiddleware
}

// NewRouter — конструктор для инициализации роутера
func NewRouter(authH *authHandler.AuthHandler, boardH *boardHandler.BoardHandler, authMW *middleware.JWTMiddleware) *Router {
	return &Router{
		authHandler:  authH,
		boardHandler: boardH,
		authMW:       authMW,
	}
}

func (r *Router) SetupRoutes(app *fiber.App) {
	// Публичные маршруты
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// --- ПУБЛИЧНЫЕ МАРШРУТЫ ---
	auth := v1.Group("/auth")
	auth.Post("/register", r.authHandler.Register)
	auth.Post("/login", r.authHandler.Login)

	// --- ЗАЩИЩЕННЫЕ МАРШРУТЫ (JWT Required) ---
	// Используем middleware для всех маршрутов ниже
	authProtected := auth.Group("/", r.authMW.VerifyToken)

	// Обновление токена
	authProtected.Post("/refresh-token", r.authHandler.Refresh)

	// Выход (удаление сессии на стороне Auth Service)
	authProtected.Post("/logout", r.authHandler.Logout)

	// Смена пароля внутри личного кабинета
	authProtected.Post("/reset-password", r.authHandler.ResetPassword)

	// Доски
	boards := v1.Group("/boards", r.authMW.VerifyToken)
	boards.Get("/", r.boardHandler.ListBoards)
	boards.Post("/", r.boardHandler.CreateBoard)
}
