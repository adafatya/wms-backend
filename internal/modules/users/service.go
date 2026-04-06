package users

import (
	"context"

	"github.com/adafatya/wms-backend/pkg/utils"
)

type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (User, error) {
	if err := ValidateCreateUser(&req); err != nil {
		return User{}, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}
	req.Password = hashedPassword

	return s.repo.CreateUser(ctx, req)
}
