package users

import (
	"context"
	"errors"
	"testing"

	"github.com/adafatya/wms-backend/internal/models"
	"github.com/adafatya/wms-backend/internal/modules/roles"
	"github.com/adafatya/wms-backend/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Manual mock repositories
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateUser(ctx context.Context, req CreateUserRequest) (User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(User), args.Error(1)
}

func (m *mockRepository) ListUsers(ctx context.Context, page, limit int) ([]User, *models.Pagination, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]User), args.Get(1).(*models.Pagination), args.Error(2)
}

func (m *mockRepository) CountUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type mockRoleRepository struct {
	mock.Mock
}

func (m *mockRoleRepository) CreateRole(ctx context.Context, name string) (roles.Role, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(roles.Role), args.Error(1)
}

func (m *mockRoleRepository) UpdateRole(ctx context.Context, id int64, name string) (roles.Role, error) {
	args := m.Called(ctx, id, name)
	return args.Get(0).(roles.Role), args.Error(1)
}

func (m *mockRoleRepository) GetRole(ctx context.Context, id int64) (roles.Role, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(roles.Role), args.Error(1)
}

func (m *mockRoleRepository) DeleteRole(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRoleRepository) ListRoles(ctx context.Context, page, limit int) ([]roles.Role, *models.Pagination, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]roles.Role), args.Get(1).(*models.Pagination), args.Error(2)
}

func (m *mockRoleRepository) CountRoles(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestService_CreateUser(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success Hashing", func(t *testing.T) {
		repo := new(mockRepository)
		roleRepo := new(mockRoleRepository)
		service := NewService(repo, roleRepo)

		rawPassword := "securepassword"
		req := CreateUserRequest{
			Username: "jdoe",
			NIK:      "1234567890",
			Password: rawPassword,
			FullName: "John Doe",
			RoleID:   1,
		}

		// Expectation: Role exists
		roleRepo.On("GetRole", mock.Anything, int64(1)).Return(roles.Role{ID: 1}, nil).Once()

		// Expectation: The password passed to the repo should be hashed, NOT the raw one.
		repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(r CreateUserRequest) bool {
			// Check if it's hashed by verifying with raw password
			err := utils.CheckPassword(rawPassword, r.Password)
			return err == nil && r.Username == "jdoe"
		})).Return(User{ID: 1, Username: "jdoe"}, nil).Once()

		res, err := service.CreateUser(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.ID)
		repo.AssertExpectations(t)
		roleRepo.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		repo := new(mockRepository)
		roleRepo := new(mockRoleRepository)
		service := NewService(repo, roleRepo)

		req := CreateUserRequest{Username: ""}

		res, err := service.CreateUser(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, User{}, res)
		repo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
	})

	t.Run("Role Not Found", func(t *testing.T) {
		repo := new(mockRepository)
		roleRepo := new(mockRoleRepository)
		service := NewService(repo, roleRepo)

		req := CreateUserRequest{
			Username: "jdoe",
			NIK:      "1234567890",
			Password: "securepassword",
			FullName: "John Doe",
			RoleID:   99,
		}

		roleRepo.On("GetRole", mock.Anything, int64(99)).Return(roles.Role{}, errors.New("not found")).Once()

		res, err := service.CreateUser(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "role not found or has been deleted", err.Error())
		assert.Equal(t, User{}, res)
		repo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
	})
}
