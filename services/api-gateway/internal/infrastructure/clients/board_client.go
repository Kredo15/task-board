package clients

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
)

type BoardClient struct {
	Client boardv1.BoardServiceClient
	Conn   *grpc.ClientConn
}

// InitBoardClient создает соединение с сервисом Board
func InitBoardClient(BoardServiceAddr string) (*BoardClient, error) {
	// Настройки KeepAlive, чтобы соединение не протухало
	ka := keepalive.ClientParameters{
		Time:                10 * time.Second, // пинговать сервер каждые 10 сек
		Timeout:             time.Second,      // ждать ответ на пинг 1 сек
		PermitWithoutStream: true,             // пинговать даже если нет активных вызовов
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(ka),
	}

	// NewClient cоздает объект соединения
	conn, err := grpc.NewClient(BoardServiceAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %w", err)
	}

	return &BoardClient{
		Client: boardv1.NewBoardServiceClient(conn),
		Conn:   conn,
	}, nil
}

// Не забудь метод для закрытия, чтобы вызвать его в main.go при graceful shutdown
func (a *BoardClient) Close() {
	_ = a.Conn.Close()
}
