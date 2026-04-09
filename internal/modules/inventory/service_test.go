package inventory

import (
	"context"
	"testing"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) WithTx(querier sqlc.Querier) Repository {
	args := m.Called(querier)
	return args.Get(0).(Repository)
}

func (m *mockRepository) CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(Product), args.Error(1)
}

func (m *mockRepository) GetProduct(ctx context.Context, id int64) (Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Product), args.Error(1)
}

func (m *mockRepository) ListProducts(ctx context.Context, search string, page, limit int) ([]Product, *models.Pagination, error) {
	args := m.Called(ctx, search, page, limit)
	return args.Get(0).([]Product), args.Get(1).(*models.Pagination), args.Error(2)
}

func (m *mockRepository) UpdateProduct(ctx context.Context, id int64, req UpdateProductRequest) (Product, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(Product), args.Error(1)
}

func (m *mockRepository) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepository) CreateLocation(ctx context.Context, req CreateLocationRequest) (Location, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(Location), args.Error(1)
}

func (m *mockRepository) GetLocation(ctx context.Context, id int64) (Location, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Location), args.Error(1)
}

func (m *mockRepository) ListLocations(ctx context.Context, search string, page, limit int) ([]Location, *models.Pagination, error) {
	args := m.Called(ctx, search, page, limit)
	return args.Get(0).([]Location), args.Get(1).(*models.Pagination), args.Error(2)
}

func (m *mockRepository) UpdateLocation(ctx context.Context, id int64, req UpdateLocationRequest) (Location, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(Location), args.Error(1)
}

func (m *mockRepository) DeleteLocation(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepository) BulkUpsertInventories(ctx context.Context, productIDs, locationIDs []int64, quantities []decimal.Decimal) error {
	args := m.Called(ctx, productIDs, locationIDs, quantities)
	return args.Error(0)
}

func (m *mockRepository) GetInventoriesByLocation(ctx context.Context, locationID int64) ([]LocationInventoryItem, error) {
	args := m.Called(ctx, locationID)
	return args.Get(0).([]LocationInventoryItem), args.Error(1)
}

func (m *mockRepository) GetInventoriesByProduct(ctx context.Context, productID int64) ([]ProductInventoryItem, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]ProductInventoryItem), args.Error(1)
}

func (m *mockRepository) ListInventories(ctx context.Context, page, limit int) ([]InventoryResponse, *models.Pagination, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]InventoryResponse), args.Get(1).(*models.Pagination), args.Error(2)
}

func (m *mockRepository) GetInventoryStock(ctx context.Context, productID, locationID int64) (decimal.Decimal, error) {
	args := m.Called(ctx, productID, locationID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

type mockStore struct {
	mock.Mock
	sqlc.Querier
}

func (m *mockStore) ExecTx(ctx context.Context, fn func(sqlc.Querier) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func TestService_UpsertInventories(t *testing.T) {
	mockRepo := new(mockRepository)
	mockStore := new(mockStore)
	svc := NewService(mockRepo, mockStore)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		req := []InventoryInput{
			{ProductID: 1, LocationID: 1, Quantity: decimal.NewFromFloat(10)},
		}

		mockRepo.On("BulkUpsertInventories", ctx, []int64{1}, []int64{1}, []decimal.Decimal{decimal.NewFromFloat(10)}).Return(nil)

		err := svc.UpsertInventories(ctx, req)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
