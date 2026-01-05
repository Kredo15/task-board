package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kredo15/task-board/services/board-service/internal/config"
	"github.com/Kredo15/task-board/services/board-service/internal/infrastructure/repository/postgres"
	db "github.com/Kredo15/task-board/services/board-service/pkg/db/postgres"
	loggerPkg "github.com/Kredo15/task-board/services/board-service/pkg/logger"
	"github.com/Kredo15/task-board/services/board-service/pkg/uuid"
	"github.com/Kredo15/task-board/services/board-service/pkg/validator"

	"github.com/Kredo15/task-board/services/board-service/internal/transport/grpc"
	usecase "github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

func Run() error {
	// конфиг сюда, возвращаем ошибку
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config parser error: %w", err)
	}

	log := loggerPkg.NewLogger(cfg)

	genID := uuid.NewGenerator()

	valid := validator.NewValidator()

	db, err := db.NewClient(cfg)
	if err != nil {
		log.Fatal("fatal to connect posgres %w", err)
		return fmt.Errorf("fatal to connect posgres %w", err)
	}
	defer db.Close()

	boardRepo := postgres.NewBoardRepository(db.GetPool())

	boardUC := usecase.NewCreateBoardUseCase(boardRepo, genID)

	srv := grpc.NewServer(cfg.GRPC.Address(), boardUC, &valid, log)
	// Запуск сервера в отдельной горутине
	go func() {
		log.Info("Starting gRPC server on %s", cfg.GRPC.Address())
		if err := srv.Run(); err != nil {
			log.Fatal("failed to serve: %v", err)
		}
	}()
	// Ожидаем сигналы от ОС (Ctrl+C, kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // Блокируемся до получения сигнала
	log.Info("Shutting down server...")

	srv.GracefulStop()
	log.Info("Server stopped")
	return nil
}
