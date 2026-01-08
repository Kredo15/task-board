package grpc

import (
	"net"

	"google.golang.org/grpc"

	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
	handler "github.com/Kredo15/task-board/services/auth-service/internal/transport/grpc/handler"
	usecase "github.com/Kredo15/task-board/services/auth-service/internal/usecase"
	"github.com/Kredo15/task-board/services/auth-service/pkg/logger"
	"github.com/Kredo15/task-board/services/auth-service/pkg/validator"
)

type Server struct {
	grpcServ *grpc.Server
	address  string
	validate *validator.Validator
	log      logger.Logger
}

func NewServer(
	addr string,
	l logger.Logger,
	valid *validator.Validator,
	regUC usecase.UserRegister,
	loginUC usecase.UserLogin,
	refreshUC usecase.TokenRefresher,
	resetUC usecase.PasswordReseter,
	logoutUC usecase.UserLogout,
) *Server {
	serv := grpc.NewServer()

	boardHandler := handler.NewHandler(
		l,
		valid,
		regUC,
		loginUC,
		refreshUC,
		resetUC,
		logoutUC,
	)

	authv1.RegisterAuthServer(serv, boardHandler)

	return &Server{
		grpcServ: serv,
		address:  addr,
		log:      l,
		validate: valid,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	return s.grpcServ.Serve(lis)
}

func (s *Server) GracefulStop() {
	s.grpcServ.GracefulStop()
}
