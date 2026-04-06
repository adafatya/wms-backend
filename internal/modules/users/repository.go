package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
)

type Repository interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (User, error)
}

type repository struct {
	querier sqlc.Querier
}

func NewRepository(querier sqlc.Querier) Repository {
	return &repository{
		querier: querier,
	}
}

func (r *repository) CreateUser(ctx context.Context, req CreateUserRequest) (User, error) {
	arg := sqlc.CreateUserParams{
		Username: req.Username,
		Nik:      sql.NullString{String: req.NIK, Valid: true},
		Password: sql.NullString{String: req.Password, Valid: true},
		FullName: sql.NullString{String: req.FullName, Valid: true},
		RoleID:   sql.NullInt64{Int64: req.RoleID, Valid: true},
	}
	u, err := r.querier.CreateUser(ctx, arg)
	if err != nil {
		return User{}, err
	}

	return mapUser(u), nil
}

func mapUser(u sqlc.User) User {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}
	return User{
		ID:        u.ID,
		Username:  u.Username,
		NIK:       u.Nik.String,
		FullName:  u.FullName.String,
		RoleID:    u.RoleID.Int64,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
