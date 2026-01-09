package grpcutil

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
)

var ErrNoUserID = errors.New("user id not found in metadata")

// GetUserID достает x-user-id из входящего контекста gRPC
func GetUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrNoUserID
	}

	values := md.Get("x-user-id")
	if len(values) == 0 || values[0] == "" {
		return "", ErrNoUserID
	}

	return values[0], nil
}