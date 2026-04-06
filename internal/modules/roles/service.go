package roles

import (
	"context"
)

type Service interface {
	CreateRole(ctx context.Context, req CreateRoleRequest) (Role, error)
	GetRole(ctx context.Context, id int64) (Role, error)
	ListRoles(ctx context.Context) ([]Role, error)
	UpdateRole(ctx context.Context, id int64, req UpdateRoleRequest) (Role, error)
	DeleteRole(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateRole(ctx context.Context, req CreateRoleRequest) (Role, error) {
	if err := ValidateCreateRoleRequest(&req); err != nil {
		return Role{}, err
	}
	return s.repo.CreateRole(ctx, req.Name)
}

func (s *service) GetRole(ctx context.Context, id int64) (Role, error) {
	return s.repo.GetRole(ctx, id)
}

func (s *service) ListRoles(ctx context.Context) ([]Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *service) UpdateRole(ctx context.Context, id int64, req UpdateRoleRequest) (Role, error) {
	if err := ValidateUpdateRoleRequest(&req); err != nil {
		return Role{}, err
	}
	return s.repo.UpdateRole(ctx, id, req.Name)
}

func (s *service) DeleteRole(ctx context.Context, id int64) error {
	return s.repo.DeleteRole(ctx, id)
}
