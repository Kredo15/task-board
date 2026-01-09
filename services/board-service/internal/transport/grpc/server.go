package grpc

import (
	"net"

	"google.golang.org/grpc"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	handler "github.com/Kredo15/task-board/services/board-service/internal/transport/grpc/handler"
	"github.com/Kredo15/task-board/services/board-service/pkg/logger"
	"github.com/Kredo15/task-board/services/board-service/pkg/validator"
)

type Server struct {
	grpcServ *grpc.Server
	address  string
	validate *validator.Validator
	log      logger.Logger
}

func NewServer(
	addr string,
	boardceratedUC handler.CreateBoard,
	boardgetUC handler.GetBoard,
	valid *validator.Validator,
	log logger.Logger,
) *Server {
	serv := grpc.NewServer()

	boardHandler := handler.NewBoardHandler(boardceratedUC, boardgetUC, valid, log)

	boardv1.RegisterBoardServiceServer(serv, boardHandler)

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
