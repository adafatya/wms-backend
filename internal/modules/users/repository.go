package users

import (
	"context"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
)

type Repository interface {
	CreateUser(ctx context.Context, username string) (User, error)
}

type repository struct {
	querier sqlc.Querier
}

func NewRepository(querier sqlc.Querier) Repository {
	return &repository{
		querier: querier,
	}
}

func (r *repository) CreateUser(ctx context.Context, username string) (User, error) {
	u, err := r.querier.CreateUser(ctx, username)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	}, nil
}
