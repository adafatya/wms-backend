package roles

import (
	"context"
	"testing"

	"github.com/adafatya/wms-backend/internal/db/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db, querier, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := NewRepository(querier)
	ctx := context.Background()

	t.Run("CreateRole", func(t *testing.T) {
		role, err := repo.CreateRole(ctx, "Admin")
		assert.NoError(t, err)
		assert.NotEmpty(t, role.ID)
		assert.Equal(t, "Admin", role.Name)
	})

	t.Run("GetRole", func(t *testing.T) {
		role, _ := repo.CreateRole(ctx, "Manager")
		
		fetched, err := repo.GetRole(ctx, role.ID)
		assert.NoError(t, err)
		assert.Equal(t, role.ID, fetched.ID)
		assert.Equal(t, "Manager", fetched.Name)
	})

	t.Run("ListRoles", func(t *testing.T) {
		testutil.TruncateTables(db) // start clean
		repo.CreateRole(ctx, "Role 1")
		repo.CreateRole(ctx, "Role 2")

		roles, err := repo.ListRoles(ctx)
		assert.NoError(t, err)
		assert.Len(t, roles, 2)
	})

	t.Run("UpdateRole", func(t *testing.T) {
		role, _ := repo.CreateRole(ctx, "Old Name")
		
		updated, err := repo.UpdateRole(ctx, role.ID, "New Name")
		assert.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
	})

	t.Run("DeleteRole", func(t *testing.T) {
		role, _ := repo.CreateRole(ctx, "To Delete")
		
		err := repo.DeleteRole(ctx, role.ID)
		assert.NoError(t, err)

		fetched, err := repo.GetRole(ctx, role.ID)
		assert.NoError(t, err)
		assert.NotNil(t, fetched.DeletedAt)
	})
}
