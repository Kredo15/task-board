package auth

import (
	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
)

// AuthHandler инкапсулирует gRPC клиент
type AuthHandler struct {
	authClient authv1.AuthClient
}

func NewAuthHandler(client authv1.AuthClient) *AuthHandler {
	return &AuthHandler{authClient: client}
}
