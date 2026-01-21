package grpc

import (
	"net"

	"google.golang.org/grpc"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	"github.com/Kredo15/task-board/services/board-service/pkg/logger"
)

type Server struct {
	grpcServ *grpc.Server
	address  string
	log      logger.Logger
}

func NewServer(addr string, handler *Handler, l logger.Logger) *Server {
	serv := grpc.NewServer()

	boardv1.RegisterBoardServiceServer(serv, handler)

	return &Server{
		grpcServ: serv,
		address:  addr,
		log:      l,
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
