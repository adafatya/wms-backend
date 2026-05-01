package inboundproduct

import (
	"context"
	"testing"
	"time"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateSchedule(ctx context.Context, arg sqlc.CreateIncomingScheduleParams) (IncomingSchedule, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(IncomingSchedule), args.Error(1)
}

func (m *mockRepository) GetSchedule(ctx context.Context, id int64) (IncomingSchedule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(IncomingSchedule), args.Error(1)
}

func (m *mockRepository) ListSchedules(ctx context.Context, page, limit int) ([]IncomingSchedule, *models.Pagination, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]IncomingSchedule), args.Get(1).(*models.Pagination), args.Error(2)
}

func (m *mockRepository) UpdateSchedule(ctx context.Context, arg sqlc.UpdateIncomingScheduleParams) (IncomingSchedule, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(IncomingSchedule), args.Error(1)
}

func (m *mockRepository) DeleteSchedule(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepository) UpsertScheduleItem(ctx context.Context, arg sqlc.UpsertIncomingScheduleItemParams) (IncomingScheduleItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(IncomingScheduleItem), args.Error(1)
}

func (m *mockRepository) GetScheduleItems(ctx context.Context, scheduleID int64) ([]IncomingScheduleItem, error) {
	args := m.Called(ctx, scheduleID)
	return args.Get(0).([]IncomingScheduleItem), args.Error(1)
}

func (m *mockRepository) DeleteScheduleItems(ctx context.Context, scheduleID int64) error {
	args := m.Called(ctx, scheduleID)
	return args.Error(0)
}

func (m *mockRepository) CreateReceipt(ctx context.Context, arg sqlc.CreateProductReceiptParams) (ProductReceipt, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(ProductReceipt), args.Error(1)
}

func (m *mockRepository) GetReceipt(ctx context.Context, id int64) (ProductReceipt, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ProductReceipt), args.Error(1)
}

func (m *mockRepository) ListReceipts(ctx context.Context, page, limit int) ([]ProductReceipt, *models.Pagination, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]ProductReceipt), args.Get(1).(*models.Pagination), args.Error(2)
}

func (m *mockRepository) BulkCreateReceiptItems(ctx context.Context, receiptID int64, productIDs []int64, quantities []decimal.Decimal) error {
	args := m.Called(ctx, receiptID, productIDs, quantities)
	return args.Error(0)
}

func (m *mockRepository) GetReceiptItems(ctx context.Context, receiptID int64) ([]ProductReceiptItem, error) {
	args := m.Called(ctx, receiptID)
	return args.Get(0).([]ProductReceiptItem), args.Error(1)
}

func (m *mockRepository) WithTx(querier sqlc.Querier) Repository {
	return m
}

type mockStore struct {
	mock.Mock
	sqlc.Querier
}

func (m *mockStore) ExecTx(ctx context.Context, fn func(sqlc.Querier) error) error {
	return fn(m.Querier)
}

func TestService_CreateSchedule(t *testing.T) {
	mockRepo := new(mockRepository)
	mockStore := new(mockStore)
	svc := NewService(mockRepo, mockStore)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		req := CreateScheduleRequest{
			LocationID:   1,
			PONumber:     "PO-001",
			ExpectedDate: "2024-05-20",
			Items: []CreateScheduleItemDTO{
				{ProductID: 1, Quantity: decimal.NewFromFloat(10)},
			},
		}

		expectedDate, _ := time.Parse("2006-01-02", "2024-05-20")

		mockRepo.On("CreateSchedule", ctx, mock.MatchedBy(func(p sqlc.CreateIncomingScheduleParams) bool {
			return p.LocationID == 1 && p.PoNumber == "PO-001" && p.ExpectedDate.Equal(expectedDate)
		})).Return(IncomingSchedule{ID: 1}, nil)

		mockRepo.On("UpsertScheduleItem", ctx, mock.MatchedBy(func(p sqlc.UpsertIncomingScheduleItemParams) bool {
			return p.IncomingScheduleID == 1 && p.ProductID == 1 && p.Quantity == "10"
		})).Return(IncomingScheduleItem{}, nil)

		mockRepo.On("GetSchedule", ctx, int64(1)).Return(IncomingSchedule{ID: 1}, nil)
		mockRepo.On("GetScheduleItems", ctx, int64(1)).Return([]IncomingScheduleItem{}, nil)

		res, err := svc.CreateSchedule(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.ID)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_UpdateSchedule(t *testing.T) {
	mockRepo := new(mockRepository)
	mockStore := new(mockStore)
	svc := NewService(mockRepo, mockStore)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		id := int64(1)
		req := UpdateScheduleRequest{
			LocationID:   1,
			PONumber:     "PO-001-UPD",
			ExpectedDate: "2024-05-21",
			Status:       StatusPending,
			Items: []UpdateScheduleItemDTO{
				{ProductID: 1, Quantity: decimal.NewFromFloat(15), ReceivedQuantity: decimal.Zero, Status: StatusPending},
			},
		}

		expectedDate, _ := time.Parse("2006-01-02", "2024-05-21")

		mockRepo.On("UpdateSchedule", ctx, mock.MatchedBy(func(p sqlc.UpdateIncomingScheduleParams) bool {
			return p.ID == id && p.PoNumber == "PO-001-UPD" && p.ExpectedDate.Equal(expectedDate)
		})).Return(IncomingSchedule{ID: id}, nil)

		mockRepo.On("UpsertScheduleItem", ctx, mock.MatchedBy(func(p sqlc.UpsertIncomingScheduleItemParams) bool {
			return p.IncomingScheduleID == id && p.ProductID == 1 && p.Quantity == "15"
		})).Return(IncomingScheduleItem{}, nil)

		mockRepo.On("GetSchedule", ctx, id).Return(IncomingSchedule{ID: id}, nil)
		mockRepo.On("GetScheduleItems", ctx, id).Return([]IncomingScheduleItem{}, nil)

		res, err := svc.UpdateSchedule(ctx, id, req)
		assert.NoError(t, err)
		assert.Equal(t, id, res.ID)
		mockRepo.AssertExpectations(t)
	})
}
