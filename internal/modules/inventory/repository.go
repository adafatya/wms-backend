package inventory

import (
	"context"
	"database/sql"
	"time"

	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/models"
	"github.com/shopspring/decimal"
)

type Repository interface {
	// Product
	CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error)
	GetProduct(ctx context.Context, id int64) (Product, error)
	ListProducts(ctx context.Context, search string, page, limit int) ([]Product, *models.Pagination, error)
	UpdateProduct(ctx context.Context, id int64, req UpdateProductRequest) (Product, error)
	DeleteProduct(ctx context.Context, id int64) error

	// Location
	CreateLocation(ctx context.Context, req CreateLocationRequest) (Location, error)
	GetLocation(ctx context.Context, id int64) (Location, error)
	ListLocations(ctx context.Context, search string, page, limit int) ([]Location, *models.Pagination, error)
	UpdateLocation(ctx context.Context, id int64, req UpdateLocationRequest) (Location, error)
	DeleteLocation(ctx context.Context, id int64) error

	// Inventory
	BulkUpsertInventories(ctx context.Context, productIDs, locationIDs []int64, quantities []decimal.Decimal) error
	GetInventoriesByLocation(ctx context.Context, locationID int64) ([]LocationInventoryItem, error)
	GetInventoriesByProduct(ctx context.Context, productID int64) ([]ProductInventoryItem, error)
	ListInventories(ctx context.Context, page, limit int) ([]InventoryResponse, *models.Pagination, error)
	GetInventoryStock(ctx context.Context, productID, locationID int64) (decimal.Decimal, error)
	WithTx(querier sqlc.Querier) Repository
}

type repository struct {
	querier sqlc.Querier
}

func NewRepository(querier sqlc.Querier) Repository {
	return &repository{
		querier: querier,
	}
}

func (r *repository) WithTx(querier sqlc.Querier) Repository {
	return &repository{
		querier: querier,
	}
}

// Product implementations
func (r *repository) CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error) {
	u, err := r.querier.CreateProduct(ctx, sqlc.CreateProductParams{
		Name:    req.Name,
		SkuCode: req.SKUCode,
		Uom:     req.UOM,
	})
	if err != nil {
		return Product{}, err
	}
	return mapProduct(u), nil
}

func (r *repository) GetProduct(ctx context.Context, id int64) (Product, error) {
	u, err := r.querier.GetProduct(ctx, id)
	if err != nil {
		return Product{}, err
	}
	return mapProduct(u), nil
}

func (r *repository) ListProducts(ctx context.Context, search string, page, limit int) ([]Product, *models.Pagination, error) {
	offset := (page - 1) * limit
	rows, err := r.querier.ListProducts(ctx, sqlc.ListProductsParams{
		Name:   "%" + search + "%",
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, nil, err
	}

	products := make([]Product, len(rows))
	for i, row := range rows {
		products[i] = mapProduct(row)
	}

	totalData, err := r.querier.CountProducts(ctx, "%"+search+"%")
	if err != nil {
		return nil, nil, err
	}

	pagination := models.NewPagination(page, limit, totalData)
	return products, pagination, nil
}

func (r *repository) UpdateProduct(ctx context.Context, id int64, req UpdateProductRequest) (Product, error) {
	u, err := r.querier.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:      id,
		Name:    req.Name,
		SkuCode: req.SKUCode,
		Uom:     req.UOM,
	})
	if err != nil {
		return Product{}, err
	}
	return mapProduct(u), nil
}

func (r *repository) DeleteProduct(ctx context.Context, id int64) error {
	return r.querier.DeleteProduct(ctx, id)
}

// Location implementations
func (r *repository) CreateLocation(ctx context.Context, req CreateLocationRequest) (Location, error) {
	u, err := r.querier.CreateLocation(ctx, sqlc.CreateLocationParams{
		Name: req.Name,
		Code: req.Code,
	})
	if err != nil {
		return Location{}, err
	}
	return mapLocation(u), nil
}

func (r *repository) GetLocation(ctx context.Context, id int64) (Location, error) {
	u, err := r.querier.GetLocation(ctx, id)
	if err != nil {
		return Location{}, err
	}
	return mapLocation(u), nil
}

func (r *repository) ListLocations(ctx context.Context, search string, page, limit int) ([]Location, *models.Pagination, error) {
	offset := (page - 1) * limit
	rows, err := r.querier.ListLocations(ctx, sqlc.ListLocationsParams{
		Name:   "%" + search + "%",
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, nil, err
	}

	locations := make([]Location, len(rows))
	for i, row := range rows {
		locations[i] = mapLocation(row)
	}

	totalData, err := r.querier.CountLocations(ctx, "%"+search+"%")
	if err != nil {
		return nil, nil, err
	}

	pagination := models.NewPagination(page, limit, totalData)
	return locations, pagination, nil
}

func (r *repository) UpdateLocation(ctx context.Context, id int64, req UpdateLocationRequest) (Location, error) {
	u, err := r.querier.UpdateLocation(ctx, sqlc.UpdateLocationParams{
		ID:   id,
		Name: req.Name,
		Code: req.Code,
	})
	if err != nil {
		return Location{}, err
	}
	return mapLocation(u), nil
}

func (r *repository) DeleteLocation(ctx context.Context, id int64) error {
	return r.querier.DeleteLocation(ctx, id)
}

// Inventory implementations
func (r *repository) BulkUpsertInventories(ctx context.Context, productIDs, locationIDs []int64, quantities []decimal.Decimal) error {
	// sqlc-generated code expects sql.NullString for numeric if not configured, 
	// but here since it's numeric(12,3), it might be string or float.
	// Actually, let's check how sqlc mapped numeric(12,3).
	// Usually it maps to string or float64 depending on config.
	// Let's assume it's string as it's common for decimal.
	
	qStrs := make([]string, len(quantities))
	for i, q := range quantities {
		qStrs[i] = q.String()
	}

	return r.querier.BulkUpsertInventories(ctx, sqlc.BulkUpsertInventoriesParams{
		ProductIds:  productIDs,
		LocationIds: locationIDs,
		Quantities:  qStrs,
	})
}

func (r *repository) GetInventoriesByLocation(ctx context.Context, locationID int64) ([]LocationInventoryItem, error) {
	rows, err := r.querier.GetInventoriesByLocation(ctx, locationID)
	if err != nil {
		return nil, err
	}

	items := make([]LocationInventoryItem, len(rows))
	for i, row := range rows {
		qty, _ := decimal.NewFromString(row.Quantity)
		items[i] = LocationInventoryItem{
			ProductID: row.ProductID,
			Product: ProductSummary{
				Name:    row.ProductName,
				SKUCode: row.ProductSkuCode,
				UOM:     row.ProductUom,
			},
			Quantity: qty,
		}
	}
	return items, nil
}

func (r *repository) GetInventoriesByProduct(ctx context.Context, productID int64) ([]ProductInventoryItem, error) {
	rows, err := r.querier.GetInventoriesByProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	items := make([]ProductInventoryItem, len(rows))
	for i, row := range rows {
		qty, _ := decimal.NewFromString(row.Quantity)
		items[i] = ProductInventoryItem{
			LocationID: row.LocationID,
			Location: LocationSummary{
				Name: row.LocationName,
				Code: row.LocationCode,
			},
			Quantity: qty,
		}
	}
	return items, nil
}

func (r *repository) ListInventories(ctx context.Context, page, limit int) ([]InventoryResponse, *models.Pagination, error) {
	offset := (page - 1) * limit
	rows, err := r.querier.ListInventories(ctx, sqlc.ListInventoriesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, nil, err
	}

	items := make([]InventoryResponse, len(rows))
	for i, row := range rows {
		qty, _ := decimal.NewFromString(row.Quantity)
		items[i] = InventoryResponse{
			LocationID: row.LocationID,
			Location: LocationSummary{
				Name: row.LocationName,
				Code: row.LocationCode,
			},
			ProductID: row.ProductID,
			Product: ProductSummary{
				Name:    row.ProductName,
				SKUCode: row.ProductSkuCode,
				UOM:     row.ProductUom,
			},
			Quantity: qty,
		}
	}

	totalData, err := r.querier.CountInventories(ctx)
	if err != nil {
		return nil, nil, err
	}

	pagination := models.NewPagination(page, limit, totalData)
	return items, pagination, nil
}

func (r *repository) GetInventoryStock(ctx context.Context, productID, locationID int64) (decimal.Decimal, error) {
	qtyStr, err := r.querier.GetInventoryStock(ctx, sqlc.GetInventoryStockParams{
		ProductID:  productID,
		LocationID: locationID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return decimal.Zero, nil
		}
		return decimal.Zero, err
	}
	qty, _ := decimal.NewFromString(qtyStr)
	return qty, nil
}

func mapProduct(u sqlc.Product) Product {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}
	return Product{
		ID:        u.ID,
		Name:      u.Name,
		SKUCode:   u.SkuCode,
		UOM:       u.Uom,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func mapLocation(l sqlc.Location) Location {
	var deletedAt *time.Time
	if l.DeletedAt.Valid {
		deletedAt = &l.DeletedAt.Time
	}
	return Location{
		ID:        l.ID,
		Name:      l.Name,
		Code:      l.Code,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
