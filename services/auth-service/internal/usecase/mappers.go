package usecase

import (
	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
)

func MapUserToResponse(u *domain.User) *UserResponse {
	return &UserResponse{
		UserID:   u.ID(),
		Email:    u.Email(),
		Username: u.Username(),
	}
}
