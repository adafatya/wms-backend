package users

import (
	"context"
	"testing"

	"github.com/adafatya/wms-backend/internal/db/testutil"
	"github.com/adafatya/wms-backend/internal/modules/roles"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	_, querier, cleanup := testutil.SetupTestDB(t)
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
}
