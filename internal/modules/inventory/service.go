package inventory

import (
	"context"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/models"
	"github.com/shopspring/decimal"
)

type Service interface {
	// Product
	CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error)
	GetProduct(ctx context.Context, id int64) (ProductDetail, error)
	ListProducts(ctx context.Context, search string, page, limit int) ([]Product, *models.Pagination, error)
	UpdateProduct(ctx context.Context, id int64, req UpdateProductRequest) (Product, error)
	DeleteProduct(ctx context.Context, id int64) error

	// Location
	CreateLocation(ctx context.Context, req CreateLocationRequest) (Location, error)
	GetLocation(ctx context.Context, id int64) (LocationDetail, error)
	ListLocations(ctx context.Context, search string, page, limit int) ([]Location, *models.Pagination, error)
	UpdateLocation(ctx context.Context, id int64, req UpdateLocationRequest) (Location, error)
	DeleteLocation(ctx context.Context, id int64) error

	// Inventory
	UpsertInventories(ctx context.Context, req []InventoryInput) error
	ListInventories(ctx context.Context, page, limit int) ([]InventoryResponse, *models.Pagination, error)
	GetInventoryStock(ctx context.Context, productID, locationID int64) (GetStokResponse, error)
}

type service struct {
	repo  Repository
	store sqlc.Store
}

func NewService(repo Repository, store sqlc.Store) Service {
	return &service{
		repo:  repo,
		store: store,
	}
}

// Product implementations
func (s *service) CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error) {
	if err := ValidateCreateProduct(&req); err != nil {
		return Product{}, err
	}

	var p Product
	err := s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)
		var err error
		p, err = txRepo.CreateProduct(ctx, req)
		if err != nil {
			return err
		}

		if len(req.Inventory) > 0 {
			productIDs := make([]int64, len(req.Inventory))
			locationIDs := make([]int64, len(req.Inventory))
			quantities := make([]decimal.Decimal, len(req.Inventory))

			for i, inv := range req.Inventory {
				productIDs[i] = p.ID
				locationIDs[i] = inv.LocationID
				quantities[i] = inv.Quantity
			}

			if err := txRepo.BulkUpsertInventories(ctx, productIDs, locationIDs, quantities); err != nil {
				return err
			}
		}
		return nil
	})

	return p, err
}

func (s *service) GetProduct(ctx context.Context, id int64) (ProductDetail, error) {
	p, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return ProductDetail{}, err
	}

	inv, err := s.repo.GetInventoriesByProduct(ctx, id)
	if err != nil {
		return ProductDetail{}, err
	}

	return ProductDetail{
		Product:   p,
		Inventory: inv,
	}, nil
}

func (s *service) ListProducts(ctx context.Context, search string, page, limit int) ([]Product, *models.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.ListProducts(ctx, search, page, limit)
}

func (s *service) UpdateProduct(ctx context.Context, id int64, req UpdateProductRequest) (Product, error) {
	if err := ValidateUpdateProduct(&req); err != nil {
		return Product{}, err
	}

	var p Product
	err := s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)
		var err error
		p, err = txRepo.UpdateProduct(ctx, id, req)
		if err != nil {
			return err
		}

		if len(req.Inventory) > 0 {
			productIDs := make([]int64, len(req.Inventory))
			locationIDs := make([]int64, len(req.Inventory))
			quantities := make([]decimal.Decimal, len(req.Inventory))

			for i, inv := range req.Inventory {
				productIDs[i] = id
				locationIDs[i] = inv.LocationID
				quantities[i] = inv.Quantity
			}

			if err := txRepo.BulkUpsertInventories(ctx, productIDs, locationIDs, quantities); err != nil {
				return err
			}
		}
		return nil
	})

	return p, err
}

func (s *service) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.DeleteProduct(ctx, id)
}

// Location implementations
func (s *service) CreateLocation(ctx context.Context, req CreateLocationRequest) (Location, error) {
	if err := ValidateCreateLocation(&req); err != nil {
		return Location{}, err
	}

	var l Location
	err := s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)
		var err error
		l, err = txRepo.CreateLocation(ctx, req)
		if err != nil {
			return err
		}

		if len(req.Inventory) > 0 {
			productIDs := make([]int64, len(req.Inventory))
			locationIDs := make([]int64, len(req.Inventory))
			quantities := make([]decimal.Decimal, len(req.Inventory))

			for i, inv := range req.Inventory {
				productIDs[i] = inv.ProductID
				locationIDs[i] = l.ID
				quantities[i] = inv.Quantity
			}

			if err := txRepo.BulkUpsertInventories(ctx, productIDs, locationIDs, quantities); err != nil {
				return err
			}
		}
		return nil
	})

	return l, err
}

func (s *service) GetLocation(ctx context.Context, id int64) (LocationDetail, error) {
	l, err := s.repo.GetLocation(ctx, id)
	if err != nil {
		return LocationDetail{}, err
	}

	inv, err := s.repo.GetInventoriesByLocation(ctx, id)
	if err != nil {
		return LocationDetail{}, err
	}

	return LocationDetail{
		Location:  l,
		Inventory: inv,
	}, nil
}

func (s *service) ListLocations(ctx context.Context, search string, page, limit int) ([]Location, *models.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.ListLocations(ctx, search, page, limit)
}

func (s *service) UpdateLocation(ctx context.Context, id int64, req UpdateLocationRequest) (Location, error) {
	if err := ValidateUpdateLocation(&req); err != nil {
		return Location{}, err
	}

	var l Location
	err := s.store.ExecTx(ctx, func(q sqlc.Querier) error {
		txRepo := s.repo.WithTx(q)
		var err error
		l, err = txRepo.UpdateLocation(ctx, id, req)
		if err != nil {
			return err
		}

		if len(req.Inventory) > 0 {
			productIDs := make([]int64, len(req.Inventory))
			locationIDs := make([]int64, len(req.Inventory))
			quantities := make([]decimal.Decimal, len(req.Inventory))

			for i, inv := range req.Inventory {
				productIDs[i] = inv.ProductID
				locationIDs[i] = id
				quantities[i] = inv.Quantity
			}

			if err := txRepo.BulkUpsertInventories(ctx, productIDs, locationIDs, quantities); err != nil {
				return err
			}
		}
		return nil
	})

	return l, err
}

func (s *service) DeleteLocation(ctx context.Context, id int64) error {
	return s.repo.DeleteLocation(ctx, id)
}

// Inventory implementations
func (s *service) UpsertInventories(ctx context.Context, req []InventoryInput) error {
	if err := ValidateUpsertInventories(req); err != nil {
		return err
	}

	if len(req) == 0 {
		return nil
	}

	productIDs := make([]int64, len(req))
	locationIDs := make([]int64, len(req))
	quantities := make([]decimal.Decimal, len(req))

	for i, inv := range req {
		productIDs[i] = inv.ProductID
		locationIDs[i] = inv.LocationID
		quantities[i] = inv.Quantity
	}

	return s.repo.BulkUpsertInventories(ctx, productIDs, locationIDs, quantities)
}

func (s *service) ListInventories(ctx context.Context, page, limit int) ([]InventoryResponse, *models.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.ListInventories(ctx, page, limit)
}

func (s *service) GetInventoryStock(ctx context.Context, productID, locationID int64) (GetStokResponse, error) {
	qty, err := s.repo.GetInventoryStock(ctx, productID, locationID)
	if err != nil {
		return GetStokResponse{}, err
	}
	return GetStokResponse{Quantity: qty}, nil
}
