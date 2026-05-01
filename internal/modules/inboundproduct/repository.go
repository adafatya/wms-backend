package inboundproduct

import (
	"context"
	"database/sql"
	"time"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/models"
	"github.com/shopspring/decimal"
)

type Repository interface {
	// Schedule
	CreateSchedule(ctx context.Context, arg sqlc.CreateIncomingScheduleParams) (IncomingSchedule, error)
	GetSchedule(ctx context.Context, id int64) (IncomingSchedule, error)
	ListSchedules(ctx context.Context, page, limit int) ([]IncomingSchedule, *models.Pagination, error)
	UpdateSchedule(ctx context.Context, arg sqlc.UpdateIncomingScheduleParams) (IncomingSchedule, error)
	DeleteSchedule(ctx context.Context, id int64) error

	// Schedule Items
	UpsertScheduleItem(ctx context.Context, arg sqlc.UpsertIncomingScheduleItemParams) (IncomingScheduleItem, error)
	GetScheduleItems(ctx context.Context, scheduleID int64) ([]IncomingScheduleItem, error)
	DeleteScheduleItems(ctx context.Context, scheduleID int64) error

	// Receipt
	CreateReceipt(ctx context.Context, arg sqlc.CreateProductReceiptParams) (ProductReceipt, error)
	GetReceipt(ctx context.Context, id int64) (ProductReceipt, error)
	ListReceipts(ctx context.Context, page, limit int) ([]ProductReceipt, *models.Pagination, error)
	BulkCreateReceiptItems(ctx context.Context, receiptID int64, productIDs []int64, quantities []decimal.Decimal) error
	GetReceiptItems(ctx context.Context, receiptID int64) ([]ProductReceiptItem, error)

	WithTx(querier sqlc.Querier) Repository
}

type repository struct {
	querier sqlc.Querier
}

func NewRepository(querier sqlc.Querier) Repository {
	return &repository{querier: querier}
}

func (r *repository) WithTx(querier sqlc.Querier) Repository {
	return &repository{querier: querier}
}

func (r *repository) CreateSchedule(ctx context.Context, arg sqlc.CreateIncomingScheduleParams) (IncomingSchedule, error) {
	s, err := r.querier.CreateIncomingSchedule(ctx, arg)
	if err != nil {
		return IncomingSchedule{}, err
	}
	return mapSchedule(s), nil
}

func (r *repository) GetSchedule(ctx context.Context, id int64) (IncomingSchedule, error) {
	s, err := r.querier.GetIncomingSchedule(ctx, id)
	if err != nil {
		return IncomingSchedule{}, err
	}
	return mapSchedule(s), nil
}

func (r *repository) ListSchedules(ctx context.Context, page, limit int) ([]IncomingSchedule, *models.Pagination, error) {
	offset := (page - 1) * limit
	rows, err := r.querier.ListIncomingSchedules(ctx, sqlc.ListIncomingSchedulesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, nil, err
	}

	total, err := r.querier.CountIncomingSchedules(ctx)
	if err != nil {
		return nil, nil, err
	}

	schedules := make([]IncomingSchedule, len(rows))
	for i, row := range rows {
		schedules[i] = mapSchedule(row)
	}

	return schedules, models.NewPagination(page, limit, total), nil
}

func (r *repository) UpdateSchedule(ctx context.Context, arg sqlc.UpdateIncomingScheduleParams) (IncomingSchedule, error) {
	s, err := r.querier.UpdateIncomingSchedule(ctx, arg)
	if err != nil {
		return IncomingSchedule{}, err
	}
	return mapSchedule(s), nil
}

func (r *repository) DeleteSchedule(ctx context.Context, id int64) error {
	res, err := r.querier.DeleteIncomingSchedule(ctx, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) UpsertScheduleItem(ctx context.Context, arg sqlc.UpsertIncomingScheduleItemParams) (IncomingScheduleItem, error) {
	item, err := r.querier.UpsertIncomingScheduleItem(ctx, arg)
	if err != nil {
		return IncomingScheduleItem{}, err
	}
	return mapScheduleItem(item), nil
}

func (r *repository) GetScheduleItems(ctx context.Context, scheduleID int64) ([]IncomingScheduleItem, error) {
	rows, err := r.querier.GetIncomingScheduleItems(ctx, scheduleID)
	if err != nil {
		return nil, err
	}

	items := make([]IncomingScheduleItem, len(rows))
	for i, row := range rows {
		items[i] = mapIncomingScheduleItemRow(row)
	}
	return items, nil
}

func (r *repository) DeleteScheduleItems(ctx context.Context, scheduleID int64) error {
	return r.querier.DeleteIncomingScheduleItems(ctx, scheduleID)
}

func (r *repository) CreateReceipt(ctx context.Context, arg sqlc.CreateProductReceiptParams) (ProductReceipt, error) {
	receipt, err := r.querier.CreateProductReceipt(ctx, arg)
	if err != nil {
		return ProductReceipt{}, err
	}
	return mapReceipt(receipt), nil
}

func (r *repository) GetReceipt(ctx context.Context, id int64) (ProductReceipt, error) {
	row, err := r.querier.GetProductReceipt(ctx, id)
	if err != nil {
		return ProductReceipt{}, err
	}
	return mapGetProductReceiptRow(row), nil
}

func (r *repository) ListReceipts(ctx context.Context, page, limit int) ([]ProductReceipt, *models.Pagination, error) {
	offset := (page - 1) * limit
	rows, err := r.querier.ListProductReceipts(ctx, sqlc.ListProductReceiptsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, nil, err
	}

	total, err := r.querier.CountProductReceipts(ctx)
	if err != nil {
		return nil, nil, err
	}

	receipts := make([]ProductReceipt, len(rows))
	for i, row := range rows {
		receipts[i] = mapListProductReceiptsRow(row)
	}

	return receipts, models.NewPagination(page, limit, total), nil
}

func (r *repository) BulkCreateReceiptItems(ctx context.Context, receiptID int64, productIDs []int64, quantities []decimal.Decimal) error {
	qStrs := make([]string, len(quantities))
	for i, q := range quantities {
		qStrs[i] = q.String()
	}

	return r.querier.BulkCreateProductReceiptItems(ctx, sqlc.BulkCreateProductReceiptItemsParams{
		ProductReceiptID: receiptID,
		ProductIds:       productIDs,
		Quantities:       qStrs,
	})
}

func (r *repository) GetReceiptItems(ctx context.Context, receiptID int64) ([]ProductReceiptItem, error) {
	rows, err := r.querier.GetProductReceiptItems(ctx, receiptID)
	if err != nil {
		return nil, err
	}

	items := make([]ProductReceiptItem, len(rows))
	for i, row := range rows {
		items[i] = mapProductReceiptItemRow(row)
	}
	return items, nil
}

// Mappers

func mapSchedule(u sqlc.IncomingSchedule) IncomingSchedule {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}
	qty, _ := decimal.NewFromString(u.ReceivedQuantity)
	return IncomingSchedule{
		ID:               u.ID,
		LocationID:       u.LocationID,
		PONumber:         u.PoNumber,
		ExpectedDate:     u.ExpectedDate,
		Status:           u.Status,
		Note:             nullStringPtr(u.Note),
		ReceivedQuantity: qty,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
		DeletedAt:        deletedAt,
	}
}

func mapScheduleItem(u sqlc.IncomingScheduleItem) IncomingScheduleItem {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}
	qty, _ := decimal.NewFromString(u.Quantity)
	recvQty, _ := decimal.NewFromString(u.ReceivedQuantity)
	return IncomingScheduleItem{
		ID:                 u.ID,
		IncomingScheduleID: u.IncomingScheduleID,
		ProductID:          u.ProductID,
		Quantity:           qty,
		ReceivedQuantity:   recvQty,
		Status:             u.Status,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
		DeletedAt:          deletedAt,
	}
}

func mapIncomingScheduleItemRow(row sqlc.GetIncomingScheduleItemsRow) IncomingScheduleItem {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	qty, _ := decimal.NewFromString(row.Quantity)
	recvQty, _ := decimal.NewFromString(row.ReceivedQuantity)
	return IncomingScheduleItem{
		ID:                 row.ID,
		IncomingScheduleID: row.IncomingScheduleID,
		ProductID:          row.ProductID,
		Product: &ProductSummary{
			Name:    row.ProductName,
			SKUCode: row.ProductSkuCode,
			UOM:     row.ProductUom,
		},
		Quantity:         qty,
		ReceivedQuantity: recvQty,
		Status:           row.Status,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
		DeletedAt:        deletedAt,
	}
}

func mapReceipt(u sqlc.ProductReceipt) ProductReceipt {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}
	var scheduleID *int64
	if u.IncomingScheduleID.Valid {
		scheduleID = &u.IncomingScheduleID.Int64
	}
	return ProductReceipt{
		ID:                 u.ID,
		IncomingScheduleID: scheduleID,
		LocationID:         u.LocationID,
		ReceivedDate:       u.ReceivedDate,
		ReceivedBy:         u.ReceivedBy,
		Note:               nullStringPtr(u.Note),
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
		DeletedAt:          deletedAt,
	}
}

func mapGetProductReceiptRow(row sqlc.GetProductReceiptRow) ProductReceipt {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	var scheduleID *int64
	if row.IncomingScheduleID.Valid {
		scheduleID = &row.IncomingScheduleID.Int64
	}
	var schedule *ScheduleSummary
	if row.SchedulePoNumber.Valid {
		schedule = &ScheduleSummary{
			PONumber:     row.SchedulePoNumber.String,
			ExpectedDate: row.ScheduleExpectedDate.Time,
		}
	}
	return ProductReceipt{
		ID:                 row.ID,
		IncomingScheduleID: scheduleID,
		IncomingSchedule:   schedule,
		LocationID:         row.LocationID,
		Location: &LocationSummary{
			Name: row.LocationName.String,
			Code: row.LocationCode.String,
		},
		ReceivedDate: row.ReceivedDate,
		ReceivedBy:   row.ReceivedBy,
		Note:         nullStringPtr(row.Note),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func mapListProductReceiptsRow(row sqlc.ListProductReceiptsRow) ProductReceipt {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	var scheduleID *int64
	if row.IncomingScheduleID.Valid {
		scheduleID = &row.IncomingScheduleID.Int64
	}
	var schedule *ScheduleSummary
	if row.SchedulePoNumber.Valid {
		schedule = &ScheduleSummary{
			PONumber:     row.SchedulePoNumber.String,
			ExpectedDate: row.ScheduleExpectedDate.Time,
		}
	}
	return ProductReceipt{
		ID:                 row.ID,
		IncomingScheduleID: scheduleID,
		IncomingSchedule:   schedule,
		LocationID:         row.LocationID,
		Location: &LocationSummary{
			Name: row.LocationName.String,
			Code: row.LocationCode.String,
		},
		ReceivedDate: row.ReceivedDate,
		ReceivedBy:   row.ReceivedBy,
		Note:         nullStringPtr(row.Note),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func mapProductReceiptItemRow(row sqlc.GetProductReceiptItemsRow) ProductReceiptItem {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	qty, _ := decimal.NewFromString(row.Quantity)
	return ProductReceiptItem{
		ID:               row.ID,
		ProductReceiptID: row.ProductReceiptID,
		ProductID:        row.ProductID,
		Product: &ProductSummary{
			Name:    row.ProductName,
			SKUCode: row.ProductSkuCode,
			UOM:     row.ProductUom,
		},
		Quantity:  qty,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func nullStringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
