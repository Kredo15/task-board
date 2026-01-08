package grpc

import (
	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
	"github.com/Kredo15/task-board/services/auth-service/internal/usecase"
	"github.com/Kredo15/task-board/services/auth-service/pkg/logger"
	"github.com/Kredo15/task-board/services/auth-service/pkg/validator"
)

// AuthHandler реализует интерфейс, сгенерированный protoc
type Handler struct {
	authv1.UnimplementedAuthServer
	log        logger.Logger
	valid      *validator.Validator
	registerUC usecase.UserRegister
	loginUC    usecase.UserLogin
	refreshUC  usecase.TokenRefresher
	resetUC    usecase.PasswordReseter
	logoutUC   usecase.UserLogout
}

func NewHandler(
	l logger.Logger,
	valid *validator.Validator,
	regUC usecase.UserRegister,
	loginUC usecase.UserLogin,
	refreshUC usecase.TokenRefresher,
	resetUC usecase.PasswordReseter,
	logoutUC usecase.UserLogout,
) *Handler {
	return &Handler{
		log:        l,
		registerUC: regUC,
		loginUC:    loginUC,
		refreshUC:  refreshUC,
		resetUC:    resetUC,
		logoutUC:   logoutUC,
	}
}
