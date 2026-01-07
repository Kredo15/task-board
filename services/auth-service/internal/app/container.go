package app

import (
	"github.com/Kredo15/task-board/services/auth-service/internal/usecase"
	"github.com/Kredo15/task-board/services/auth-service/pkg/logger"
)

// Container содержит все UseCase сервиса
type Container struct {
	Logger          logger.Logger
	RegisterUC      *usecase.RegisterUseCase
	LoginUC         *usecase.LoginUseCase
	RefreshUC       *usecase.RefreshUseCase
	LogoutUC        *usecase.LogoutUseCase
	ResetPasswordUC *usecase.ResetUseCase
}

func NewContainer(
	logger logger.Logger,
	reg *usecase.RegisterUseCase,
	log *usecase.LoginUseCase,
	ref *usecase.RefreshUseCase,
	out *usecase.LogoutUseCase,
	reset *usecase.ResetUseCase,
) *Container {
	return &Container{
		Logger:          logger,
		RegisterUC:      reg,
		LoginUC:         log,
		RefreshUC:       ref,
		LogoutUC:        out,
		ResetPasswordUC: reset,
	}
}
