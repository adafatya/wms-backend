package inventory

import (
	"time"

	"github.com/shopspring/decimal"
)

// Product models
type Product struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	SKUCode   string     `json:"sku_code"`
	UOM       string     `json:"uom"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type ProductDetail struct {
	Product
	Inventory []ProductInventoryItem `json:"inventory"`
}

type ProductInventoryItem struct {
	LocationID int64            `json:"location_id"`
	Location   LocationSummary  `json:"location"`
	Quantity   decimal.Decimal  `json:"quantity"`
}

type LocationSummary struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type CreateProductRequest struct {
	Name      string           `json:"name" binding:"required"`
	SKUCode   string           `json:"sku_code" binding:"required"`
	UOM       string           `json:"uom" binding:"required"`
	Inventory []InventoryInput `json:"inventory"`
}

type UpdateProductRequest struct {
	Name      string           `json:"name"`
	SKUCode   string           `json:"sku_code"`
	UOM       string           `json:"uom"`
	Inventory []InventoryInput `json:"inventory"`
}

// Location models
type Location struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type LocationDetail struct {
	Location
	Inventory []LocationInventoryItem `json:"inventory"`
}

type LocationInventoryItem struct {
	ProductID int64           `json:"product_id"`
	Product   ProductSummary  `json:"product"`
	Quantity  decimal.Decimal `json:"quantity"`
}

type ProductSummary struct {
	Name    string `json:"name"`
	SKUCode string `json:"sku_code"`
	UOM     string `json:"uom"`
}

type CreateLocationRequest struct {
	Name      string           `json:"name" binding:"required"`
	Code      string           `json:"code" binding:"required"`
	Inventory []InventoryInput `json:"inventory"`
}

type UpdateLocationRequest struct {
	Name      string           `json:"name"`
	Code      string           `json:"code"`
	Inventory []InventoryInput `json:"inventory"`
}

// Inventory models
type InventoryInput struct {
	ProductID  int64           `json:"product_id,omitempty"`
	LocationID int64           `json:"location_id,omitempty"`
	Quantity   decimal.Decimal `json:"quantity" binding:"required"`
}

type InventoryResponse struct {
	LocationID int64           `json:"location_id"`
	Location   LocationSummary `json:"location"`
	ProductID  int64           `json:"product_id"`
	Product    ProductSummary  `json:"product"`
	Quantity   decimal.Decimal `json:"quantity"`
}

type GetStokResponse struct {
	Quantity decimal.Decimal `json:"quantity"`
}
