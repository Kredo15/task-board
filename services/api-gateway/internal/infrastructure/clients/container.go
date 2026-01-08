package clients

import (
	"context"

	"google.golang.org/grpc"

	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	"github.com/Kredo15/task-board/services/api-gateway/pkg/logger"
)

type Clients struct {
	Auth  authv1.AuthClient
	Board boardv1.BoardServiceClient
	// Список соединений для корректного закрытия (Graceful Shutdown)
	connections []*grpc.ClientConn
	log         logger.Logger
}

func NewClients(l logger.Logger) *Clients {
	return &Clients{
		connections: make([]*grpc.ClientConn, 0),
		log:         l,
	}
}

func (c *Clients) InitConnections(ctx context.Context, authAddr, boardAddr string) error {
	// Инициализируем Auth
	c.log.Info("initializing grpc auth clients")
	authClient, err := InitAuthClient(authAddr)
	if err != nil {
		c.log.Error("failed to connect to auth service", err)
		return err
	}
	c.Auth = authClient.Client
	c.connections = append(c.connections, authClient.Conn)

	// Инициализируем Board
	c.log.Info("initializing grpc board clients")
	boardClient, err := InitBoardClient(boardAddr)
	if err != nil {
		c.log.Error("failed to connect to board service", err)
		return err
	}
	c.Board = boardClient.Client
	c.connections = append(c.connections, boardClient.Conn)

	return nil
}

// Close закрывает все открытые gRPC соединения
func (c *Clients) Close() error {
	for _, conn := range c.connections {
		if err := conn.Close(); err != nil {
			c.log.Error("error closing grpc connection", err)
			return err
		}
	}
	return nil
}
