package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Kredo15/task-board/services/api-gateway/internal/config"
	"github.com/Kredo15/task-board/services/api-gateway/internal/infrastructure/clients"
	"github.com/Kredo15/task-board/services/api-gateway/internal/transport/http"
	"github.com/Kredo15/task-board/services/api-gateway/internal/transport/http/handler/auth"
	"github.com/Kredo15/task-board/services/api-gateway/internal/transport/http/handler/board"
	"github.com/Kredo15/task-board/services/api-gateway/internal/transport/http/middleware"
	loggerPkg "github.com/Kredo15/task-board/services/api-gateway/pkg/logger"
)

type App struct {
	fiber     *fiber.App
	container *clients.Clients // Сохраняем контейнер, чтобы закрыть соединения
	cfg       *config.Config
	log       loggerPkg.Logger
}

func NewApp() (*App, error) {
	// конфиг сюда, возвращаем ошибку
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("config parser error: %w", err)
	}

	log := loggerPkg.NewLogger(cfg.Level)

	ctx := context.Background()
	// Инициализация gRPC клиентов
	clientsContainer := clients.NewClients(log)
	err_conn := clientsContainer.InitConnections(ctx, cfg.AuthServiceAddr, cfg.BoardServiceAddr)
	if err_conn != nil {
		return nil, fmt.Errorf("fatal to connect grpc service %w", err)
	}
	// Инициализация хендлеров
	authHandler := auth.NewAuthHandler(clientsContainer.Auth)
	boardHandler := board.NewBoardHandler(clientsContainer.Board)

	// Инициализация middleware
	authMW := middleware.NewJWTMiddleware(cfg.JWTSecret)

	// Настройка Fiber и Роутера
	app := fiber.New()
	router := http.NewRouter(authHandler, boardHandler, authMW)
	router.SetupRoutes(app)
	return &App{
		fiber:     app,
		container: clientsContainer,
		log:       log,
		cfg:       cfg,
	}, nil
}

func (a *App) Run() error {

	go func() {
		if err := a.fiber.Listen(fmt.Sprintf(":%d", a.cfg.Port)); err != nil {
			a.log.Fatal("Failed to start HTTP server: %v", err)
		}
	}()

	// Ожидаем сигналы от ОС (Ctrl+C, kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // Блокируемся до получения сигнала
	a.log.Info("Shutting down server...")
	return a.Stop()
}

func (a *App) Stop() error {
	// 1. Устанавливаем таймаут на закрытие (например, 10 секунд)
	// Это дает время Fiber завершить текущие HTTP-запросы
	const shutdownTimeout = 10 * time.Second

	// 2. Останавливаем HTTP сервер
	if err := a.fiber.ShutdownWithTimeout(shutdownTimeout); err != nil {
		a.log.Error("Fiber shutdown error: %v", err)
	}

	// 3. Закрываем gRPC соединения в контейнере
	// В container.go должен быть метод Close()
	if err := a.container.Close(); err != nil {
		a.log.Error("gRPC clients close error: %v", err)
	}

	a.log.Info("Application stopped gracefully")
	return nil
}
