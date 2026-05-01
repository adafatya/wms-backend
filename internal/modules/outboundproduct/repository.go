package outboundproduct

import (
	"context"
	"database/sql"
	"time"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/models"
	"github.com/shopspring/decimal"
)

type Repository interface {
	// Customer
	CreateCustomer(ctx context.Context, arg sqlc.CreateCustomerParams) (Customer, error)
	GetCustomer(ctx context.Context, id int64) (Customer, error)
	ListCustomers(ctx context.Context, page, limit int) ([]Customer, *models.Pagination, error)
	UpdateCustomer(ctx context.Context, arg sqlc.UpdateCustomerParams) (Customer, error)
	DeleteCustomer(ctx context.Context, id int64) error

	// Delivery Order
	CreateDeliveryOrder(ctx context.Context, arg sqlc.CreateDeliveryOrderParams) (DeliveryOrder, error)
	GetDeliveryOrder(ctx context.Context, id int64) (DeliveryOrder, error)
	ListDeliveryOrders(ctx context.Context, page, limit int) ([]DeliveryOrder, *models.Pagination, error)
	UpdateDeliveryOrder(ctx context.Context, arg sqlc.UpdateDeliveryOrderParams) (DeliveryOrder, error)
	UpdateDeliveryOrderStatus(ctx context.Context, arg sqlc.UpdateDeliveryOrderStatusParams) (DeliveryOrder, error)
	DeleteDeliveryOrder(ctx context.Context, id int64) error

	// Delivery Order Items
	BulkCreateDeliveryOrderItems(ctx context.Context, doID int64, productIDs []int64, quantities []decimal.Decimal) error
	GetDeliveryOrderItems(ctx context.Context, doID int64) ([]DeliveryOrderItem, error)
	UpsertDeliveryOrderItem(ctx context.Context, arg sqlc.UpsertDeliveryOrderItemParams) (DeliveryOrderItem, error)
	DeleteDeliveryOrderItems(ctx context.Context, doID int64) error
	IncrementDeliveryOrderItemDeliveredQty(ctx context.Context, arg sqlc.IncrementDeliveryOrderItemDeliveredQtyParams) error

	// Delivery
	CreateDelivery(ctx context.Context, arg sqlc.CreateDeliveryParams) (Delivery, error)
	GetDelivery(ctx context.Context, id int64) (Delivery, error)
	ListDeliveries(ctx context.Context, page, limit int) ([]Delivery, *models.Pagination, error)
	UpdateDelivery(ctx context.Context, arg sqlc.UpdateDeliveryParams) (Delivery, error)
	DeleteDelivery(ctx context.Context, id int64) error

	// Delivery Items
	BulkCreateDeliveryItems(ctx context.Context, deliveryID int64, productIDs []int64, quantities []decimal.Decimal) error
	GetDeliveryItems(ctx context.Context, deliveryID int64) ([]DeliveryItem, error)

	// Inventory
	GetInventoryStock(ctx context.Context, arg sqlc.GetInventoryStockParams) (decimal.Decimal, error)
	BulkDeductInventories(ctx context.Context, arg sqlc.BulkDeductInventoriesParams) error

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

// Customer
func (r *repository) CreateCustomer(ctx context.Context, arg sqlc.CreateCustomerParams) (Customer, error) {
	c, err := r.querier.CreateCustomer(ctx, arg)
	if err != nil {
		return Customer{}, err
	}
	return mapCustomer(c), nil
}

func (r *repository) GetCustomer(ctx context.Context, id int64) (Customer, error) {
	c, err := r.querier.GetCustomer(ctx, id)
	if err != nil {
		return Customer{}, err
	}
	return mapCustomer(c), nil
}

func (r *repository) ListCustomers(ctx context.Context, page, limit int) ([]Customer, *models.Pagination, error) {
	offset := (page - 1) * limit
	rows, err := r.querier.ListCustomers(ctx, sqlc.ListCustomersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, nil, err
	}

	total, err := r.querier.CountCustomers(ctx)
	if err != nil {
		return nil, nil, err
	}

	customers := make([]Customer, len(rows))
	for i, row := range rows {
		customers[i] = mapCustomer(row)
	}

	return customers, models.NewPagination(page, limit, total), nil
}

func (r *repository) UpdateCustomer(ctx context.Context, arg sqlc.UpdateCustomerParams) (Customer, error) {
	c, err := r.querier.UpdateCustomer(ctx, arg)
	if err != nil {
		return Customer{}, err
	}
	return mapCustomer(c), nil
}

func (r *repository) DeleteCustomer(ctx context.Context, id int64) error {
	res, err := r.querier.DeleteCustomer(ctx, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delivery Order
func (r *repository) CreateDeliveryOrder(ctx context.Context, arg sqlc.CreateDeliveryOrderParams) (DeliveryOrder, error) {
	do, err := r.querier.CreateDeliveryOrder(ctx, arg)
	if err != nil {
		return DeliveryOrder{}, err
	}
	return mapDeliveryOrder(do), nil
}

func (r *repository) GetDeliveryOrder(ctx context.Context, id int64) (DeliveryOrder, error) {
	row, err := r.querier.GetDeliveryOrder(ctx, id)
	if err != nil {
		return DeliveryOrder{}, err
	}
	return mapGetDeliveryOrderRow(row), nil
}

func (r *repository) ListDeliveryOrders(ctx context.Context, page, limit int) ([]DeliveryOrder, *models.Pagination, error) {
	offset := (page - 1) * limit
	rows, err := r.querier.ListDeliveryOrders(ctx, sqlc.ListDeliveryOrdersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, nil, err
	}

	total, err := r.querier.CountDeliveryOrders(ctx)
	if err != nil {
		return nil, nil, err
	}

	orders := make([]DeliveryOrder, len(rows))
	for i, row := range rows {
		orders[i] = mapListDeliveryOrdersRow(row)
	}

	return orders, models.NewPagination(page, limit, total), nil
}

func (r *repository) UpdateDeliveryOrder(ctx context.Context, arg sqlc.UpdateDeliveryOrderParams) (DeliveryOrder, error) {
	do, err := r.querier.UpdateDeliveryOrder(ctx, arg)
	if err != nil {
		return DeliveryOrder{}, err
	}
	return mapDeliveryOrder(do), nil
}

func (r *repository) UpdateDeliveryOrderStatus(ctx context.Context, arg sqlc.UpdateDeliveryOrderStatusParams) (DeliveryOrder, error) {
	do, err := r.querier.UpdateDeliveryOrderStatus(ctx, arg)
	if err != nil {
		return DeliveryOrder{}, err
	}
	return mapDeliveryOrder(do), nil
}

func (r *repository) DeleteDeliveryOrder(ctx context.Context, id int64) error {
	res, err := r.querier.DeleteDeliveryOrder(ctx, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delivery Order Items
func (r *repository) BulkCreateDeliveryOrderItems(ctx context.Context, doID int64, productIDs []int64, quantities []decimal.Decimal) error {
	qtyStrs := make([]string, len(quantities))
	for i, q := range quantities {
		qtyStrs[i] = q.String()
	}
	return r.querier.BulkCreateDeliveryOrderItems(ctx, sqlc.BulkCreateDeliveryOrderItemsParams{
		DeliveryOrderID: doID,
		Column2:         productIDs,
		Column3:         qtyStrs,
	})
}

func (r *repository) GetDeliveryOrderItems(ctx context.Context, doID int64) ([]DeliveryOrderItem, error) {
	rows, err := r.querier.GetDeliveryOrderItems(ctx, doID)
	if err != nil {
		return nil, err
	}
	items := make([]DeliveryOrderItem, len(rows))
	for i, row := range rows {
		items[i] = mapDeliveryOrderItemRow(row)
	}
	return items, nil
}

func (r *repository) UpsertDeliveryOrderItem(ctx context.Context, arg sqlc.UpsertDeliveryOrderItemParams) (DeliveryOrderItem, error) {
	item, err := r.querier.UpsertDeliveryOrderItem(ctx, arg)
	if err != nil {
		return DeliveryOrderItem{}, err
	}
	return mapDeliveryOrderItem(item), nil
}

func (r *repository) DeleteDeliveryOrderItems(ctx context.Context, doID int64) error {
	return r.querier.DeleteDeliveryOrderItems(ctx, doID)
}

func (r *repository) IncrementDeliveryOrderItemDeliveredQty(ctx context.Context, arg sqlc.IncrementDeliveryOrderItemDeliveredQtyParams) error {
	return r.querier.IncrementDeliveryOrderItemDeliveredQty(ctx, arg)
}

// Delivery
func (r *repository) CreateDelivery(ctx context.Context, arg sqlc.CreateDeliveryParams) (Delivery, error) {
	d, err := r.querier.CreateDelivery(ctx, arg)
	if err != nil {
		return Delivery{}, err
	}
	return mapDelivery(d), nil
}

func (r *repository) GetDelivery(ctx context.Context, id int64) (Delivery, error) {
	row, err := r.querier.GetDelivery(ctx, id)
	if err != nil {
		return Delivery{}, err
	}
	return mapGetDeliveryRow(row), nil
}

func (r *repository) ListDeliveries(ctx context.Context, page, limit int) ([]Delivery, *models.Pagination, error) {
	offset := (page - 1) * limit
	rows, err := r.querier.ListDeliveries(ctx, sqlc.ListDeliveriesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, nil, err
	}

	total, err := r.querier.CountDeliveries(ctx)
	if err != nil {
		return nil, nil, err
	}

	deliveries := make([]Delivery, len(rows))
	for i, row := range rows {
		deliveries[i] = mapListDeliveriesRow(row)
	}

	return deliveries, models.NewPagination(page, limit, total), nil
}

func (r *repository) UpdateDelivery(ctx context.Context, arg sqlc.UpdateDeliveryParams) (Delivery, error) {
	d, err := r.querier.UpdateDelivery(ctx, arg)
	if err != nil {
		return Delivery{}, err
	}
	return mapDelivery(d), nil
}

func (r *repository) DeleteDelivery(ctx context.Context, id int64) error {
	res, err := r.querier.DeleteDelivery(ctx, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delivery Items
func (r *repository) BulkCreateDeliveryItems(ctx context.Context, deliveryID int64, productIDs []int64, quantities []decimal.Decimal) error {
	qtyStrs := make([]string, len(quantities))
	for i, q := range quantities {
		qtyStrs[i] = q.String()
	}
	return r.querier.BulkCreateDeliveryItems(ctx, sqlc.BulkCreateDeliveryItemsParams{
		DeliveryID: deliveryID,
		Column2:    productIDs,
		Column3:    qtyStrs,
	})
}

func (r *repository) GetDeliveryItems(ctx context.Context, deliveryID int64) ([]DeliveryItem, error) {
	rows, err := r.querier.GetDeliveryItems(ctx, deliveryID)
	if err != nil {
		return nil, err
	}
	items := make([]DeliveryItem, len(rows))
	for i, row := range rows {
		items[i] = mapDeliveryItemRow(row)
	}
	return items, nil
}

// Inventory
func (r *repository) GetInventoryStock(ctx context.Context, arg sqlc.GetInventoryStockParams) (decimal.Decimal, error) {
	qtyStr, err := r.querier.GetInventoryStock(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return decimal.Zero, nil
		}
		return decimal.Zero, err
	}
	qty, _ := decimal.NewFromString(qtyStr)
	return qty, nil
}

func (r *repository) BulkDeductInventories(ctx context.Context, arg sqlc.BulkDeductInventoriesParams) error {
	return r.querier.BulkDeductInventories(ctx, arg)
}

// Mappers
func mapCustomer(c sqlc.Customer) Customer {
	var deletedAt *time.Time
	if c.DeletedAt.Valid {
		deletedAt = &c.DeletedAt.Time
	}
	return Customer{
		ID:          c.ID,
		Name:        c.Name,
		Address:     c.Address,
		ContactName: nullStringPtr(c.ContactName),
		ContactInfo: nullStringPtr(c.ContactInfo),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

func mapDeliveryOrder(do sqlc.DeliveryOrder) DeliveryOrder {
	var deletedAt *time.Time
	if do.DeletedAt.Valid {
		deletedAt = &do.DeletedAt.Time
	}
	return DeliveryOrder{
		ID:           do.ID,
		CustomerID:   do.CustomerID,
		OrderNumber:  do.OrderNumber,
		DeliveryDate: do.DeliveryDate,
		Status:       do.Status,
		Note:         nullStringPtr(do.Note),
		CreatedAt:    do.CreatedAt,
		UpdatedAt:    do.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func mapGetDeliveryOrderRow(row sqlc.GetDeliveryOrderRow) DeliveryOrder {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	return DeliveryOrder{
		ID:           row.ID,
		CustomerID:   row.CustomerID,
		Customer: &CustomerSummary{
			Name:    row.CustomerName,
			Address: row.CustomerAddress,
		},
		OrderNumber:  row.OrderNumber,
		DeliveryDate: row.DeliveryDate,
		Status:       row.Status,
		Note:         nullStringPtr(row.Note),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func mapListDeliveryOrdersRow(row sqlc.ListDeliveryOrdersRow) DeliveryOrder {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	return DeliveryOrder{
		ID:           row.ID,
		CustomerID:   row.CustomerID,
		Customer: &CustomerSummary{
			Name: row.CustomerName,
		},
		OrderNumber:  row.OrderNumber,
		DeliveryDate: row.DeliveryDate,
		Status:       row.Status,
		Note:         nullStringPtr(row.Note),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func mapDeliveryOrderItem(item sqlc.DeliveryOrderItem) DeliveryOrderItem {
	var deletedAt *time.Time
	if item.DeletedAt.Valid {
		deletedAt = &item.DeletedAt.Time
	}
	qty, _ := decimal.NewFromString(item.Quantity)
	recvQty, _ := decimal.NewFromString(item.DeliveredQuantity)
	return DeliveryOrderItem{
		ID:                item.ID,
		DeliveryOrderID:   item.DeliveryOrderID,
		ProductID:         item.ProductID,
		Quantity:          qty,
		DeliveredQuantity: recvQty,
		CreatedAt:         item.CreatedAt,
		UpdatedAt:         item.UpdatedAt,
		DeletedAt:         deletedAt,
	}
}

func mapDeliveryOrderItemRow(row sqlc.GetDeliveryOrderItemsRow) DeliveryOrderItem {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	qty, _ := decimal.NewFromString(row.Quantity)
	recvQty, _ := decimal.NewFromString(row.DeliveredQuantity)
	return DeliveryOrderItem{
		ID:              row.ID,
		DeliveryOrderID: row.DeliveryOrderID,
		ProductID:       row.ProductID,
		Product: &ProductSummary{
			Name:    row.ProductName,
			SKUCode: row.ProductSkuCode,
			UOM:     row.ProductUom,
		},
		Quantity:          qty,
		DeliveredQuantity: recvQty,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
		DeletedAt:         deletedAt,
	}
}

func mapDelivery(d sqlc.Delivery) Delivery {
	var deletedAt *time.Time
	if d.DeletedAt.Valid {
		deletedAt = &d.DeletedAt.Time
	}
	return Delivery{
		ID:              d.ID,
		DeliveryOrderID: d.DeliveryOrderID,
		UserID:          d.UserID,
		LocationID:      d.LocationID,
		DeliveredAt:     d.DeliveredAt,
		VehicleNumber:   nullStringPtr(d.VehicleNumber),
		Note:            nullStringPtr(d.Note),
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
		DeletedAt:       deletedAt,
	}
}

func mapGetDeliveryRow(row sqlc.GetDeliveryRow) Delivery {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	return Delivery{
		ID:              row.ID,
		DeliveryOrderID: row.DeliveryOrderID,
		OrderNumber:     row.OrderNumber,
		UserID:          row.UserID,
		UserName:        row.UserName,
		LocationID:      row.LocationID,
		LocationName:    row.LocationName,
		DeliveredAt:     row.DeliveredAt,
		VehicleNumber:   nullStringPtr(row.VehicleNumber),
		Note:            nullStringPtr(row.Note),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		DeletedAt:       deletedAt,
	}
}

func mapListDeliveriesRow(row sqlc.ListDeliveriesRow) Delivery {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	return Delivery{
		ID:              row.ID,
		DeliveryOrderID: row.DeliveryOrderID,
		OrderNumber:     row.OrderNumber,
		UserID:          row.UserID,
		LocationID:      row.LocationID,
		DeliveredAt:     row.DeliveredAt,
		VehicleNumber:   nullStringPtr(row.VehicleNumber),
		Note:            nullStringPtr(row.Note),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		DeletedAt:       deletedAt,
	}
}

func mapDeliveryItemRow(row sqlc.GetDeliveryItemsRow) DeliveryItem {
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		deletedAt = &row.DeletedAt.Time
	}
	qty, _ := decimal.NewFromString(row.Quantity)
	return DeliveryItem{
		ID:         row.ID,
		DeliveryID: row.DeliveryID,
		ProductID:  row.ProductID,
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
