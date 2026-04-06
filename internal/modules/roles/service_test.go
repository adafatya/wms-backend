package roles

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Manual mock repository
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateRole(ctx context.Context, name string) (Role, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(Role), args.Error(1)
}

func (m *mockRepository) GetRole(ctx context.Context, id int64) (Role, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Role), args.Error(1)
}

func (m *mockRepository) ListRoles(ctx context.Context) ([]Role, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Role), args.Error(1)
}

func (m *mockRepository) UpdateRole(ctx context.Context, id int64, name string) (Role, error) {
	args := m.Called(ctx, id, name)
	return args.Get(0).(Role), args.Error(1)
}

func (m *mockRepository) DeleteRole(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_CreateRole(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		repo := new(mockRepository)
		service := NewService(repo)

		req := CreateRoleRequest{Name: "Admin"}
		expectedRole := Role{ID: 1, Name: "Admin"}

		repo.On("CreateRole", mock.Anything, "Admin").Return(expectedRole, nil).Once()

		res, err := service.CreateRole(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedRole, res)
		repo.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		repo := new(mockRepository)
		service := NewService(repo)

		req := CreateRoleRequest{Name: ""}

		res, err := service.CreateRole(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, Role{}, res)
		repo.AssertNotCalled(t, "CreateRole", mock.Anything, mock.Anything)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo := new(mockRepository)
		service := NewService(repo)

		req := CreateRoleRequest{Name: "Editor"}

		repo.On("CreateRole", mock.Anything, "Editor").Return(Role{}, errors.New("db error")).Once()

		res, err := service.CreateRole(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, Role{}, res)
		assert.Equal(t, "db error", err.Error())
		repo.AssertExpectations(t)
	})
}
