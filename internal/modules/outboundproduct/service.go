package outboundproduct

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/models"
	"github.com/shopspring/decimal"
)

type Service interface {
	// Customer
	CreateCustomer(ctx context.Context, req CreateCustomerRequest) (Customer, error)
	GetCustomer(ctx context.Context, id int64) (Customer, error)
	ListCustomers(ctx context.Context, page, limit int) ([]Customer, *models.Pagination, error)
	UpdateCustomer(ctx context.Context, id int64, req UpdateCustomerRequest) (Customer, error)
	DeleteCustomer(ctx context.Context, id int64) error

	// Delivery Order
	CreateDeliveryOrder(ctx context.Context, req CreateDeliveryOrderRequest) (DeliveryOrder, error)
	GetDeliveryOrder(ctx context.Context, id int64) (DeliveryOrder, error)
	ListDeliveryOrders(ctx context.Context, page, limit int) ([]DeliveryOrder, *models.Pagination, error)
	UpdateDeliveryOrder(ctx context.Context, id int64, req UpdateDeliveryOrderRequest) (DeliveryOrder, error)
	DeleteDeliveryOrder(ctx context.Context, id int64) error

	// Delivery
	CreateDelivery(ctx context.Context, req CreateDeliveryRequest) (Delivery, error)
	GetDelivery(ctx context.Context, id int64) (Delivery, error)
	ListDeliveries(ctx context.Context, page, limit int) ([]Delivery, *models.Pagination, error)
	UpdateDelivery(ctx context.Context, id int64, req UpdateDeliveryRequest) (Delivery, error)
	DeleteDelivery(ctx context.Context, id int64) error
}

type service struct {
	repo  Repository
	store sqlc.Store
}

func NewService(repo Repository, store sqlc.Store) Service {
	return &service{repo: repo, store: store}
}

// Customer
func (s *service) CreateCustomer(ctx context.Context, req CreateCustomerRequest) (Customer, error) {
	if err := ValidateCreateCustomer(req); err != nil {
		return Customer{}, err
	}

	return s.repo.CreateCustomer(ctx, sqlc.CreateCustomerParams{
		Name:        req.Name,
		Address:     req.Address,
		ContactName: sql.NullString{String: req.ContactName, Valid: req.ContactName != ""},
		ContactInfo: sql.NullString{String: req.ContactInfo, Valid: req.ContactInfo != ""},
	})
}

func (s *service) GetCustomer(ctx context.Context, id int64) (Customer, error) {
	return s.repo.GetCustomer(ctx, id)
}

func (s *service) ListCustomers(ctx context.Context, page, limit int) ([]Customer, *models.Pagination, error) {
	return s.repo.ListCustomers(ctx, page, limit)
}

func (s *service) UpdateCustomer(ctx context.Context, id int64, req UpdateCustomerRequest) (Customer, error) {
	if err := ValidateUpdateCustomer(req); err != nil {
		return Customer{}, err
	}

	return s.repo.UpdateCustomer(ctx, sqlc.UpdateCustomerParams{
		ID:          id,
		Name:        req.Name,
		Address:     req.Address,
		ContactName: sql.NullString{String: req.ContactName, Valid: req.ContactName != ""},
		ContactInfo: sql.NullString{String: req.ContactInfo, Valid: req.ContactInfo != ""},
	})
}

func (s *service) DeleteCustomer(ctx context.Context, id int64) error {
	return s.repo.DeleteCustomer(ctx, id)
}

// Delivery Order
func (s *service) CreateDeliveryOrder(ctx context.Context, req CreateDeliveryOrderRequest) (DeliveryOrder, error) {
	if err := ValidateCreateDeliveryOrder(req); err != nil {
		return DeliveryOrder{}, err
	}

	deliveryDate, _ := time.Parse("2006-01-02", req.DeliveryDate)

	var result DeliveryOrder
	err := s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)

		do, err := txRepo.CreateDeliveryOrder(ctx, sqlc.CreateDeliveryOrderParams{
			CustomerID:   req.CustomerID,
			OrderNumber:  req.OrderNumber,
			DeliveryDate: deliveryDate,
			Note:         sql.NullString{String: req.Note, Valid: req.Note != ""},
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

		if err := txRepo.BulkCreateDeliveryOrderItems(ctx, do.ID, productIDs, quantities); err != nil {
			return err
		}

		result = do
		return nil
	})

	if err != nil {
		return DeliveryOrder{}, err
	}

	return s.GetDeliveryOrder(ctx, result.ID)
}

func (s *service) GetDeliveryOrder(ctx context.Context, id int64) (DeliveryOrder, error) {
	do, err := s.repo.GetDeliveryOrder(ctx, id)
	if err != nil {
		return DeliveryOrder{}, err
	}

	items, err := s.repo.GetDeliveryOrderItems(ctx, id)
	if err != nil {
		return DeliveryOrder{}, err
	}

	do.Items = items
	return do, nil
}

func (s *service) ListDeliveryOrders(ctx context.Context, page, limit int) ([]DeliveryOrder, *models.Pagination, error) {
	return s.repo.ListDeliveryOrders(ctx, page, limit)
}

func (s *service) UpdateDeliveryOrder(ctx context.Context, id int64, req UpdateDeliveryOrderRequest) (DeliveryOrder, error) {
	if err := ValidateUpdateDeliveryOrder(req); err != nil {
		return DeliveryOrder{}, err
	}

	deliveryDate, _ := time.Parse("2006-01-02", req.DeliveryDate)

	err := s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)

		_, err := txRepo.UpdateDeliveryOrder(ctx, sqlc.UpdateDeliveryOrderParams{
			ID:           id,
			CustomerID:   req.CustomerID,
			OrderNumber:  req.OrderNumber,
			DeliveryDate: deliveryDate,
			Status:       req.Status,
			Note:         sql.NullString{String: req.Note, Valid: req.Note != ""},
		})
		if err != nil {
			return err
		}

		// Reconciliation
		if err := txRepo.DeleteDeliveryOrderItems(ctx, id); err != nil {
			return err
		}

		for _, item := range req.Items {
			_, err := txRepo.UpsertDeliveryOrderItem(ctx, sqlc.UpsertDeliveryOrderItemParams{
				DeliveryOrderID:   id,
				ProductID:         item.ProductID,
				Quantity:          item.Quantity.String(),
				DeliveredQuantity: item.DeliveredQuantity.String(),
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return DeliveryOrder{}, err
	}

	return s.GetDeliveryOrder(ctx, id)
}

func (s *service) DeleteDeliveryOrder(ctx context.Context, id int64) error {
	return s.repo.DeleteDeliveryOrder(ctx, id)
}

// Delivery
func (s *service) CreateDelivery(ctx context.Context, req CreateDeliveryRequest) (Delivery, error) {
	if err := ValidateCreateDelivery(req); err != nil {
		return Delivery{}, err
	}

	deliveredAt, _ := time.Parse("2006-01-02 15:04:05", req.DeliveredAt)

	var result Delivery
	err := s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)

		// 1. Stock Check
		for _, item := range req.Items {
			stock, err := txRepo.GetInventoryStock(ctx, sqlc.GetInventoryStockParams{
				ProductID:  item.ProductID,
				LocationID: req.LocationID,
			})
			if err != nil {
				return err
			}
			if stock.LessThan(item.Quantity) {
				return fmt.Errorf("insufficient stock for product ID %d: have %s, need %s", item.ProductID, stock.String(), item.Quantity.String())
			}
		}

		// 2. Create Delivery Header
		d, err := txRepo.CreateDelivery(ctx, sqlc.CreateDeliveryParams{
			DeliveryOrderID: req.DeliveryOrderID,
			UserID:          req.UserID,
			LocationID:      req.LocationID,
			DeliveredAt:     deliveredAt,
			VehicleNumber:   sql.NullString{String: req.VehicleNumber, Valid: req.VehicleNumber != ""},
			Note:            sql.NullString{String: req.Note, Valid: req.Note != ""},
		})
		if err != nil {
			return err
		}

		productIDs := make([]int64, len(req.Items))
		quantities := make([]decimal.Decimal, len(req.Items))
		locationIDs := make([]int64, len(req.Items))
		for i, item := range req.Items {
			productIDs[i] = item.ProductID
			quantities[i] = item.Quantity
			locationIDs[i] = req.LocationID

			// 3. Update Delivery Order Items delivered quantity
			err := txRepo.IncrementDeliveryOrderItemDeliveredQty(ctx, sqlc.IncrementDeliveryOrderItemDeliveredQtyParams{
				DeliveryOrderID:   req.DeliveryOrderID,
				ProductID:         item.ProductID,
				DeliveredQuantity: item.Quantity.String(),
			})
			if err != nil {
				return err
			}
		}

		// 4. Create Delivery Items
		if err := txRepo.BulkCreateDeliveryItems(ctx, d.ID, productIDs, quantities); err != nil {
			return err
		}

		// 5. Deduct Inventory
		qtyStrs := make([]string, len(quantities))
		for i, q := range quantities {
			qtyStrs[i] = q.String()
		}
		if err := txRepo.BulkDeductInventories(ctx, sqlc.BulkDeductInventoriesParams{
			ProductIds:  productIDs,
			LocationIds: locationIDs,
			Quantities:  qtyStrs,
		}); err != nil {
			return err
		}

		// 6. Update Delivery Order Status
		// Check if all items delivered
		doItems, err := txRepo.GetDeliveryOrderItems(ctx, req.DeliveryOrderID)
		if err != nil {
			return err
		}

		allCompleted := true
		anyStarted := false
		for _, doi := range doItems {
			if doi.DeliveredQuantity.GreaterThan(decimal.Zero) {
				anyStarted = true
			}
			if doi.DeliveredQuantity.LessThan(doi.Quantity) {
				allCompleted = false
			}
		}

		newStatus := StatusPending
		if allCompleted {
			newStatus = StatusCompleted
		} else if anyStarted {
			newStatus = StatusProcessing
		}

		if _, err := txRepo.UpdateDeliveryOrderStatus(ctx, sqlc.UpdateDeliveryOrderStatusParams{
			ID:     req.DeliveryOrderID,
			Status: newStatus,
		}); err != nil {
			return err
		}

		result = d
		return nil
	})

	if err != nil {
		return Delivery{}, err
	}

	return s.GetDelivery(ctx, result.ID)
}

func (s *service) GetDelivery(ctx context.Context, id int64) (Delivery, error) {
	d, err := s.repo.GetDelivery(ctx, id)
	if err != nil {
		return Delivery{}, err
	}

	items, err := s.repo.GetDeliveryItems(ctx, id)
	if err != nil {
		return Delivery{}, err
	}

	d.Items = items
	return d, nil
}

func (s *service) ListDeliveries(ctx context.Context, page, limit int) ([]Delivery, *models.Pagination, error) {
	return s.repo.ListDeliveries(ctx, page, limit)
}

func (s *service) UpdateDelivery(ctx context.Context, id int64, req UpdateDeliveryRequest) (Delivery, error) {
	if err := ValidateUpdateDelivery(req); err != nil {
		return Delivery{}, err
	}

	deliveredAt, _ := time.Parse("2006-01-02 15:04:05", req.DeliveredAt)

	// Note: Update delivery is complex because it might affect stock again.
	// For simplicity in this plan, we just update the header info.
	// In real WMS, updating quantities would require reverting previous stock and applying new ones.
	
	return s.repo.UpdateDelivery(ctx, sqlc.UpdateDeliveryParams{
		ID:              id,
		DeliveryOrderID: req.DeliveryOrderID,
		UserID:          req.UserID,
		LocationID:      req.LocationID,
		DeliveredAt:     deliveredAt,
		VehicleNumber:   sql.NullString{String: req.VehicleNumber, Valid: req.VehicleNumber != ""},
		Note:            sql.NullString{String: req.Note, Valid: req.Note != ""},
	})
}

func (s *service) DeleteDelivery(ctx context.Context, id int64) error {
	return s.repo.DeleteDelivery(ctx, id)
}
