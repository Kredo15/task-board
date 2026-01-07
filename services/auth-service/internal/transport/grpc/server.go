package grpc

import (
	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
	"github.com/Kredo15/task-board/services/auth-service/internal/app"
)

// AuthHandler реализует интерфейс, сгенерированный protoc
type Handler struct {
	authv1.UnimplementedAuthServer
	deps *app.Container
}

func NewHandler(deps *app.Container) *Handler {
	return &Handler{
		deps: deps,
	}
}
