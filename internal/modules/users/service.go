package users

import (
	"context"
	"errors"

	"github.com/adafatya/wms-backend/internal/models"
	"github.com/adafatya/wms-backend/internal/modules/roles"
	"github.com/adafatya/wms-backend/pkg/utils"
)

type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (User, error)
	ListUsers(ctx context.Context, page, limit int) ([]User, *models.Pagination, error)
}

type service struct {
	repo     Repository
	roleRepo roles.Repository
}

func NewService(repo Repository, roleRepo roles.Repository) Service {
	return &service{
		repo:     repo,
		roleRepo: roleRepo,
	}
}

func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (User, error) {
	if err := ValidateCreateUser(&req); err != nil {
		return User{}, err
	}

	// Check if role exists
	_, err := s.roleRepo.GetRole(ctx, req.RoleID)
	if err != nil {
		return User{}, errors.New("role not found or has been deleted")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}
	req.Password = hashedPassword

	return s.repo.CreateUser(ctx, req)
}

func (s *service) ListUsers(ctx context.Context, page, limit int) ([]User, *models.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return s.repo.ListUsers(ctx, page, limit)
}
