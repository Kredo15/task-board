package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kredo15/task-board/services/board-service/internal/config"
	"github.com/Kredo15/task-board/services/board-service/internal/infrastructure/repository/postgres"
	db "github.com/Kredo15/task-board/services/board-service/pkg/db/postgres"
	"github.com/Kredo15/task-board/services/board-service/pkg/db/redis"
	"github.com/Kredo15/task-board/services/board-service/pkg/lexorank"
	loggerPkg "github.com/Kredo15/task-board/services/board-service/pkg/logger"
	"github.com/Kredo15/task-board/services/board-service/pkg/uuid"

	cacheRepo "github.com/Kredo15/task-board/services/board-service/internal/infrastructure/cache/redis"
	"github.com/Kredo15/task-board/services/board-service/internal/transport/grpc"
	usecaseBoard "github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
	usecaseCache "github.com/Kredo15/task-board/services/board-service/internal/usecase/cache"
	usecaseColumn "github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
	usecaseTask "github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
)

func Run() error {
	// конфиг сюда, возвращаем ошибку
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config parser error: %w", err)
	}

	log := loggerPkg.NewLogger(cfg.Logging.Level)

	genID := uuid.NewGenerator()
	lexorankGen := lexorank.NewLexorankGen()

	db, err := db.NewClient(cfg)
	if err != nil {
		log.Fatal("fatal to connect posgres %w", err)
		return fmt.Errorf("fatal to connect posgres %w", err)
	}
	defer db.Close()

	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		log.Fatal("fatal to connect redis %w", err)
		return fmt.Errorf("fatal to connect redis %w", err)
	}
	defer redisClient.Close()

	boardRepo := postgres.NewBoardRepository(db.GetPool())

	// usecase Board
	boardCreateUC := usecaseBoard.NewCreateBoardUseCase(boardRepo, genID)
	boardGetUC := usecaseBoard.NewGetBoardUseCase(boardRepo)
	boardUpdateUC := usecaseBoard.NewUpdateBoardUseCase(boardRepo, genID)
	boardDeleteUC := usecaseBoard.NewDeleteBoardUseCase(boardRepo, genID)

	// usecase Column
	columnAddUc := usecaseColumn.NewAddColumnUseCase(boardRepo, genID, lexorankGen)
	columnUpdateUC := usecaseColumn.NewUpdateColumnUseCase(boardRepo, genID)
	columnMoveUC := usecaseColumn.NewMoveColumnUseCase(boardRepo, genID, lexorankGen)
	columnDeleteUC := usecaseColumn.NewDeleteTaskUseCase(boardRepo, genID)

	// usecase Task
	taskAddUC := usecaseTask.NewAddTaskUseCase(boardRepo, genID, lexorankGen)
	taskUpdataUC := usecaseTask.NewUpdateTaskUseCase(boardRepo, genID)
	taskMoveUC := usecaseTask.NewMoveTaskUseCase(boardRepo, genID, lexorankGen)
	taskDeleteUC := usecaseTask.NewDeleteTaskUseCase(boardRepo, genID)

	boardCache := cacheRepo.NewBoardCache(redisClient, cfg.Redis.TTL)

	// Для кэширования применяем паттерн Decorator
	//usecase cache board
	cacheBoardUpdateUC := usecaseCache.NewCachedUpdateBoardUseCase(boardUpdateUC, boardCache)
	cacheBoardGetUC := usecaseCache.NewCachedGetBoardUseCase(boardGetUC, boardCache)
	cacheBoardDeleteUC := usecaseCache.NewCachedDeleteBoardUseCase(boardDeleteUC, boardCache)
	// usecase cache column
	cacheColumnAddUC := usecaseCache.NewCachedCreateColumnUseCase(columnAddUc, boardCache)
	cacheColumnUpdateUC := usecaseCache.NewCachedUpdateColumnUseCase(columnUpdateUC, boardCache)
	cacheColumnMoveUC := usecaseCache.NewCachedMoveColumnUseCase(columnMoveUC, boardCache)
	cacheColumnDeleteUC := usecaseCache.NewCachedDeleteColumnUseCase(columnDeleteUC, boardCache)
	// usecaaw cache task
	cacheTaskAddUC := usecaseCache.NewCachedAddTaskUseCase(taskAddUC, boardCache)
	cacheTaskUpdateUC := usecaseCache.NewCachedUpdateTaskUseCase(taskUpdataUC, boardCache)
	cacheTaskMoveUC := usecaseCache.NewCachedMoveTaskUseCase(taskMoveUC, boardCache)
	cacheTaskDeleUC := usecaseCache.NewCachedDeleteTaskUseCase(taskDeleteUC, boardCache)

	boardUC := grpc.BoardUseCases{
		Create: boardCreateUC,
		Update: cacheBoardUpdateUC,
		Get:    cacheBoardGetUC,
		Delete: cacheBoardDeleteUC,
	}

	columnUC := grpc.ColumnUseCases{
		Create: cacheColumnAddUC,
		Update: cacheColumnUpdateUC,
		Move:   cacheColumnMoveUC,
		Delete: cacheColumnDeleteUC,
	}

	taskUC := grpc.TaskUseCases{
		Create: cacheTaskAddUC,
		Update: cacheTaskUpdateUC,
		Move:   cacheTaskMoveUC,
		Delete: cacheTaskDeleUC,
	}

	handler := grpc.NewHandler(
		grpc.Deps{
			Board:  boardUC,
			Column: columnUC,
			Task:   taskUC,
		},
		log,
	)

	srv := grpc.NewServer(cfg.GRPC.Address(), handler, log)
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
