package users

import (
	"context"
	"testing"

	"github.com/adafatya/wms-backend/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Manual mock repository
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateUser(ctx context.Context, req CreateUserRequest) (User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(User), args.Error(1)
}

func TestService_CreateUser(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success Hashing", func(t *testing.T) {
		repo := new(mockRepository)
		service := NewService(repo)

		rawPassword := "securepassword"
		req := CreateUserRequest{
			Username: "jdoe",
			NIK:      "1234567890",
			Password: rawPassword,
			FullName: "John Doe",
			RoleID:   1,
		}

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
	})

	t.Run("Validation Error", func(t *testing.T) {
		repo := new(mockRepository)
		service := NewService(repo)

		req := CreateUserRequest{Username: ""}

		res, err := service.CreateUser(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, User{}, res)
		repo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
	})
}
