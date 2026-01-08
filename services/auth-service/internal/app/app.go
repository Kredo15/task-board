package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"

	"github.com/Kredo15/task-board/services/auth-service/internal/config"
	repoPg "github.com/Kredo15/task-board/services/auth-service/internal/infrastructure/repository/postgres"
	repoRedis "github.com/Kredo15/task-board/services/auth-service/internal/infrastructure/repository/redis"
	"github.com/Kredo15/task-board/services/auth-service/internal/infrastructure/security"
	"github.com/Kredo15/task-board/services/auth-service/internal/transport/grpc"
	"github.com/Kredo15/task-board/services/auth-service/internal/usecase"
	db "github.com/Kredo15/task-board/services/auth-service/pkg/db/postgres"
	"github.com/Kredo15/task-board/services/auth-service/pkg/logger"
	"github.com/Kredo15/task-board/services/auth-service/pkg/uuid"
	"github.com/Kredo15/task-board/services/auth-service/pkg/validator"
)

func Run() error {
	// конфиг сюда, возвращаем ошибку
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config parser error: %w", err)
	}

	log := logger.NewLogger(cfg.Logging.Level)

	genID := uuid.NewGenerator()

	valid := validator.NewValidator()

	db, err := db.NewClient(cfg)
	if err != nil {
		log.Fatal("fatal to connect posgres %w", err)
		return fmt.Errorf("fatal to connect posgres %w", err)
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host,
		DB:   cfg.Redis.DB,
	})

	// Infrastructure (Реализация интерфейсов)
	userRepo := repoPg.NewUserRepository(db.GetPool())
	tokenRepo := repoRedis.NewTokenRepository(rdb)
	hasher := security.NewBcryptHasher(cfg.Auth.BcryptCost)
	jwtManager := security.NewJWTManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.AccessTTL,
		cfg.Auth.RefreshTTL,
	)

	// Usecase (Бизнес-логика)
	registerUC := usecase.NewRegisterUseCase(userRepo, tokenRepo, hasher, jwtManager, genID, log)
	loginUC := usecase.NewLoginUseCase(userRepo, tokenRepo, hasher, jwtManager, log)
	refUC := usecase.NewRefreshUseCase(tokenRepo, jwtManager, log)
	outUC := usecase.NewLogoutUseCase(tokenRepo, log)
	resetUC := usecase.NewResetUseCase(userRepo, tokenRepo, hasher, jwtManager, log)

	// Сервер
	srv := grpc.NewServer(cfg.GRPC.Address(), log, &valid, registerUC, loginUC, refUC, resetUC, outUC)
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
