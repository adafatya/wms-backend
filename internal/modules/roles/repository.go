package roles

import (
	"context"
	"time"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/models"
)

type Repository interface {
	CreateRole(ctx context.Context, name string) (Role, error)
	GetRole(ctx context.Context, id int64) (Role, error)
	ListRoles(ctx context.Context, page, limit int) ([]Role, *models.Pagination, error)
	CountRoles(ctx context.Context) (int64, error)
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

func (r *repository) ListRoles(ctx context.Context, page, limit int) ([]Role, *models.Pagination, error) {
	offset := (page - 1) * limit
	arg := sqlc.ListRolesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	rows, err := r.querier.ListRoles(ctx, arg)
	if err != nil {
		return nil, nil, err
	}
	roles := make([]Role, len(rows))
	for i, row := range rows {
		roles[i] = mapRole(row)
	}

	totalData, err := r.querier.CountRoles(ctx)
	if err != nil {
		return nil, nil, err
	}

	pagination := models.NewPagination(page, limit, totalData)

	return roles, pagination, nil
}

func (r *repository) CountRoles(ctx context.Context) (int64, error) {
	return r.querier.CountRoles(ctx)
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
