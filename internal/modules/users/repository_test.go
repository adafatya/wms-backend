package users

import (
	"context"
	"testing"

	"github.com/adafatya/wms-backend/internal/db/testutil"
	"github.com/adafatya/wms-backend/internal/modules/roles"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db, querier, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := NewRepository(querier)
	roleRepo := roles.NewRepository(querier)
	ctx := context.Background()

	t.Run("CreateUser", func(t *testing.T) {
		role, _ := roleRepo.CreateRole(ctx, "Operator")

		req := CreateUserRequest{
			Username: "johndoe",
			NIK:      "1234567890",
			Password: "hashedpassword",
			FullName: "John Doe",
			RoleID:   role.ID,
		}

		user, err := repo.CreateUser(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, req.Username, user.Username)
		assert.Equal(t, req.FullName, user.FullName)
		assert.Equal(t, req.RoleID, user.RoleID)
	})

	t.Run("ListUsers", func(t *testing.T) {
		// Truncate first to be sure
		testutil.TruncateTables(db)
		role, err := roleRepo.CreateRole(ctx, "Admin")
		assert.NoError(t, err)

		_, err = repo.CreateUser(ctx, CreateUserRequest{Username: "user1", NIK: "nik1", RoleID: role.ID})
		assert.NoError(t, err)
		_, err = repo.CreateUser(ctx, CreateUserRequest{Username: "user2", NIK: "nik2", RoleID: role.ID})
		assert.NoError(t, err)

		users, pagination, err := repo.ListUsers(ctx, 1, 10)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.NotNil(t, pagination)
		assert.Equal(t, 1, pagination.TotalPage)
	})
}
