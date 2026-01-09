package grpc

import (
	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
	"github.com/Kredo15/task-board/services/auth-service/pkg/logger"
	"github.com/Kredo15/task-board/services/auth-service/pkg/validator"
)

// AuthHandler реализует интерфейс, сгенерированный protoc
type Handler struct {
	authv1.UnimplementedAuthServer
	log        logger.Logger
	valid      *validator.Validator
	registerUC UserRegister
	loginUC    UserLogin
	refreshUC  TokenRefresher
	resetUC    PasswordReseter
	logoutUC   UserLogout
}

func NewHandler(
	l logger.Logger,
	valid *validator.Validator,
	regUC UserRegister,
	loginUC UserLogin,
	refreshUC TokenRefresher,
	resetUC PasswordReseter,
	logoutUC UserLogout,
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
