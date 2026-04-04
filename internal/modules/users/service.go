package users

import (
	"context"
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
	return s.repo.CreateUser(ctx, req.Username)
}
