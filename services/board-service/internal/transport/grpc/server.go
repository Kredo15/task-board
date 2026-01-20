package grpc

import (
	"net"

	"google.golang.org/grpc"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	"github.com/Kredo15/task-board/services/board-service/pkg/logger"
	"github.com/Kredo15/task-board/services/board-service/pkg/validator"
)

type Server struct {
	grpcServ *grpc.Server
	address  string
	validate *validator.Validator
	log      logger.Logger
}

func NewServer(addr string, handler *Handler) *Server {
	serv := grpc.NewServer()

	boardv1.RegisterBoardServiceServer(serv, handler)

	return &Server{
		grpcServ: serv,
		address:  addr,
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
