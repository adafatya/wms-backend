package roles

import (
	"context"
	"time"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
)

type Repository interface {
	CreateRole(ctx context.Context, name string) (Role, error)
	GetRole(ctx context.Context, id int64) (Role, error)
	ListRoles(ctx context.Context) ([]Role, error)
	UpdateRole(ctx context.Context, id int64, name string) (Role, error)
	DeleteRole(ctx context.Context, id int64) error
}

type repository struct {
	querier sqlc.Querier
}

func NewRepository(querier sqlc.Querier) Repository {
	return &repository{
		querier: querier,
	}
}

func (r *repository) CreateRole(ctx context.Context, name string) (Role, error) {
	row, err := r.querier.CreateRole(ctx, name)
	if err != nil {
		return Role{}, err
	}
	return mapRole(row), nil
}

func (r *repository) GetRole(ctx context.Context, id int64) (Role, error) {
	row, err := r.querier.GetRole(ctx, id)
	if err != nil {
		return Role{}, err
	}
	return mapRole(row), nil
}

func (r *repository) ListRoles(ctx context.Context) ([]Role, error) {
	rows, err := r.querier.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	roles := make([]Role, len(rows))
	for i, row := range rows {
		roles[i] = mapRole(row)
	}
	return roles, nil
}

func (r *repository) UpdateRole(ctx context.Context, id int64, name string) (Role, error) {
	arg := sqlc.UpdateRoleParams{
		ID:   id,
		Name: name,
	}
	row, err := r.querier.UpdateRole(ctx, arg)
	if err != nil {
		return Role{}, err
	}
	return mapRole(row), nil
}

func (r *repository) DeleteRole(ctx context.Context, id int64) error {
	return r.querier.DeleteRole(ctx, id)
}

func mapRole(row sqlc.Role) Role {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	return Role{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
