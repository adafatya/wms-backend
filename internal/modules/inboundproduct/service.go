package inboundproduct

import (
	"context"
	"database/sql"
	"time"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/models"
	"github.com/shopspring/decimal"
)

type Service interface {
	CreateSchedule(ctx context.Context, req CreateScheduleRequest) (IncomingSchedule, error)
	GetSchedule(ctx context.Context, id int64) (IncomingSchedule, error)
	ListSchedules(ctx context.Context, page, limit int) ([]IncomingSchedule, *models.Pagination, error)
	UpdateSchedule(ctx context.Context, id int64, req UpdateScheduleRequest) (IncomingSchedule, error)
	DeleteSchedule(ctx context.Context, id int64) error

	CreateReceipt(ctx context.Context, req CreateReceiptRequest) (ProductReceipt, error)
	GetReceipt(ctx context.Context, id int64) (ProductReceipt, error)
	ListReceipts(ctx context.Context, page, limit int) ([]ProductReceipt, *models.Pagination, error)
}

type service struct {
	repo  Repository
	store sqlc.Store
}

func NewService(repo Repository, store sqlc.Store) Service {
	return &service{repo: repo, store: store}
}

func (s *service) CreateSchedule(ctx context.Context, req CreateScheduleRequest) (IncomingSchedule, error) {
	if err := ValidateCreateSchedule(&req); err != nil {
		return IncomingSchedule{}, err
	}

	expectedDate, _ := time.Parse("2006-01-02", req.ExpectedDate)

	var schedule IncomingSchedule
	err := s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)
		
		created, err := txRepo.CreateSchedule(ctx, sqlc.CreateIncomingScheduleParams{
			LocationID:   req.LocationID,
			PoNumber:     req.PONumber,
			ExpectedDate: expectedDate,
			Note:         sql.NullString{String: req.Note, Valid: req.Note != ""},
		})
		if err != nil {
			return err
		}

		for _, item := range req.Items {
			_, err := txRepo.UpsertScheduleItem(ctx, sqlc.UpsertIncomingScheduleItemParams{
				IncomingScheduleID: created.ID,
				ProductID:          item.ProductID,
				Quantity:           item.Quantity.String(),
				ReceivedQuantity:   decimal.Zero.String(),
				Status:             StatusPending,
			})
			if err != nil {
				return err
			}
		}

		schedule = created
		return nil
	})

	if err != nil {
		return IncomingSchedule{}, err
	}

	return s.GetSchedule(ctx, schedule.ID)
}

func (s *service) GetSchedule(ctx context.Context, id int64) (IncomingSchedule, error) {
	schedule, err := s.repo.GetSchedule(ctx, id)
	if err != nil {
		return IncomingSchedule{}, err
	}

	items, err := s.repo.GetScheduleItems(ctx, id)
	if err != nil {
		return IncomingSchedule{}, err
	}
	schedule.Items = items

	return schedule, nil
}

func (s *service) ListSchedules(ctx context.Context, page, limit int) ([]IncomingSchedule, *models.Pagination, error) {
	return s.repo.ListSchedules(ctx, page, limit)
}

func (s *service) UpdateSchedule(ctx context.Context, id int64, req UpdateScheduleRequest) (IncomingSchedule, error) {
	if err := ValidateUpdateSchedule(&req); err != nil {
		return IncomingSchedule{}, err
	}

	expectedDate, _ := time.Parse("2006-01-02", req.ExpectedDate)

	err := s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)

		_, err := txRepo.UpdateSchedule(ctx, sqlc.UpdateIncomingScheduleParams{
			ID:           id,
			LocationID:   req.LocationID,
			PoNumber:     req.PONumber,
			ExpectedDate: expectedDate,
			Status:       req.Status,
			Note:         sql.NullString{String: req.Note, Valid: req.Note != ""},
		})
		if err != nil {
			return err
		}

		// Logic for items reconciliation/upsert
		for _, item := range req.Items {
			_, err := txRepo.UpsertScheduleItem(ctx, sqlc.UpsertIncomingScheduleItemParams{
				IncomingScheduleID: id,
				ProductID:          item.ProductID,
				Quantity:           item.Quantity.String(),
				ReceivedQuantity:   item.ReceivedQuantity.String(),
				Status:             item.Status,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return IncomingSchedule{}, err
	}

	return s.GetSchedule(ctx, id)
}

func (s *service) DeleteSchedule(ctx context.Context, id int64) error {
	return s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)
		if err := txRepo.DeleteScheduleItems(ctx, id); err != nil {
			return err
		}
		return txRepo.DeleteSchedule(ctx, id)
	})
}

func (s *service) CreateReceipt(ctx context.Context, req CreateReceiptRequest) (ProductReceipt, error) {
	if err := ValidateCreateReceipt(&req); err != nil {
		return ProductReceipt{}, err
	}

	receivedDate, _ := time.Parse("2006-01-02", req.ReceivedDate)

	var receipt ProductReceipt
	err := s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)

		var scheduleID sql.NullInt64
		if req.IncomingScheduleID != nil {
			scheduleID = sql.NullInt64{Int64: *req.IncomingScheduleID, Valid: true}
		}

		created, err := txRepo.CreateReceipt(ctx, sqlc.CreateProductReceiptParams{
			IncomingScheduleID: scheduleID,
			LocationID:         req.LocationID,
			ReceivedDate:       receivedDate,
			ReceivedBy:         req.ReceivedBy,
			Note:               sql.NullString{String: req.Note, Valid: req.Note != ""},
		})
		if err != nil {
			return err
		}

		productIDs := make([]int64, len(req.Items))
		quantities := make([]decimal.Decimal, len(req.Items))
		for i, item := range req.Items {
			productIDs[i] = item.ProductID
			quantities[i] = item.Quantity
		}

		if err := txRepo.BulkCreateReceiptItems(ctx, created.ID, productIDs, quantities); err != nil {
			return err
		}

		receipt = created
		return nil
	})

	if err != nil {
		return ProductReceipt{}, err
	}

	return s.GetReceipt(ctx, receipt.ID)
}

func (s *service) GetReceipt(ctx context.Context, id int64) (ProductReceipt, error) {
	receipt, err := s.repo.GetReceipt(ctx, id)
	if err != nil {
		return ProductReceipt{}, err
	}

	items, err := s.repo.GetReceiptItems(ctx, id)
	if err != nil {
		return ProductReceipt{}, err
	}
	receipt.Items = items

	return receipt, nil
}

func (s *service) ListReceipts(ctx context.Context, page, limit int) ([]ProductReceipt, *models.Pagination, error) {
	return s.repo.ListReceipts(ctx, page, limit)
}
