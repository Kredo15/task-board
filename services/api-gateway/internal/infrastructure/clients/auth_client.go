package clients

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
)

type AuthClient struct {
	Client authv1.AuthClient
	Conn   *grpc.ClientConn
}

// InitAuthClient создает соединение с сервисом Auth
func InitAuthClient(AuthServiceAddr string) (*AuthClient, error) {
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
	conn, err := grpc.NewClient(AuthServiceAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %w", err)
	}

	return &AuthClient{
		Client: authv1.NewAuthClient(conn),
		Conn:   conn,
	}, nil
}

func (a *AuthClient) Close() {
	_ = a.Conn.Close()
}
