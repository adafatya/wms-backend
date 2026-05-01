package outboundproduct

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusShipped    = "shipped"
	StatusCancelled  = "cancelled"
	StatusCompleted  = "completed"
)

// Customer
type CreateCustomerRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	ContactName string `json:"contact_name"`
	ContactInfo string `json:"contact_info"`
}

type UpdateCustomerRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	ContactName string `json:"contact_name"`
	ContactInfo string `json:"contact_info"`
}

type Customer struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	ContactName *string    `json:"contact_name,omitempty"`
	ContactInfo *string    `json:"contact_info,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// Delivery Order
type CreateDeliveryOrderRequest struct {
	CustomerID   int64                      `json:"customer_id" binding:"required"`
	OrderNumber  string                     `json:"order_number" binding:"required"`
	DeliveryDate string                     `json:"delivery_date" binding:"required"` // YYYY-MM-DD
	Note         string                     `json:"note"`
	Items        []CreateDeliveryOrderItem `json:"items" binding:"required,min=1"`
}

type CreateDeliveryOrderItem struct {
	ProductID int64           `json:"product_id" binding:"required"`
	Quantity  decimal.Decimal `json:"quantity" binding:"required"`
}

type UpdateDeliveryOrderRequest struct {
	CustomerID   int64                      `json:"customer_id" binding:"required"`
	OrderNumber  string                     `json:"order_number" binding:"required"`
	DeliveryDate string                     `json:"delivery_date" binding:"required"`
	Status       string                     `json:"status" binding:"required"`
	Note         string                     `json:"note"`
	Items        []UpdateDeliveryOrderItem `json:"items" binding:"required,min=1"`
}

type UpdateDeliveryOrderItem struct {
	ProductID         int64           `json:"product_id" binding:"required"`
	Quantity          decimal.Decimal `json:"quantity" binding:"required"`
	DeliveredQuantity decimal.Decimal `json:"delivered_quantity"`
}

type DeliveryOrder struct {
	ID           int64               `json:"id"`
	CustomerID   int64               `json:"customer_id"`
	Customer     *CustomerSummary    `json:"customer,omitempty"`
	OrderNumber  string              `json:"order_number"`
	DeliveryDate time.Time           `json:"delivery_date"`
	Status       string              `json:"status"`
	Note         *string             `json:"note,omitempty"`
	Items        []DeliveryOrderItem `json:"items,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	DeletedAt    *time.Time          `json:"deleted_at,omitempty"`
}

type CustomerSummary struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type DeliveryOrderItem struct {
	ID                int64           `json:"id"`
	DeliveryOrderID   int64           `json:"delivery_order_id"`
	ProductID         int64           `json:"product_id"`
	Product           *ProductSummary `json:"product,omitempty"`
	Quantity          decimal.Decimal `json:"quantity"`
	DeliveredQuantity decimal.Decimal `json:"delivered_quantity"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         *time.Time      `json:"deleted_at,omitempty"`
}

type ProductSummary struct {
	Name    string `json:"name"`
	SKUCode string `json:"sku_code"`
	UOM     string `json:"uom"`
}

// Delivery
type CreateDeliveryRequest struct {
	DeliveryOrderID int64                  `json:"delivery_order_id" binding:"required"`
	LocationID      int64                  `json:"location_id" binding:"required"`
	UserID          int64                  `json:"user_id" binding:"required"`
	DeliveredAt     string                 `json:"delivered_at" binding:"required"` // YYYY-MM-DD HH:mm:ss
	VehicleNumber   string                 `json:"vehicle_number"`
	Note            string                 `json:"note"`
	Items           []CreateDeliveryItemDTO `json:"items" binding:"required,min=1"`
}

type CreateDeliveryItemDTO struct {
	ProductID int64           `json:"product_id" binding:"required"`
	Quantity  decimal.Decimal `json:"quantity" binding:"required"`
}

type UpdateDeliveryRequest struct {
	DeliveryOrderID int64                  `json:"delivery_order_id" binding:"required"`
	LocationID      int64                  `json:"location_id" binding:"required"`
	UserID          int64                  `json:"user_id" binding:"required"`
	DeliveredAt     string                 `json:"delivered_at" binding:"required"`
	VehicleNumber   string                 `json:"vehicle_number"`
	Note            string                 `json:"note"`
	Items           []UpdateDeliveryItemDTO `json:"items" binding:"required,min=1"`
}

type UpdateDeliveryItemDTO struct {
	ProductID int64           `json:"product_id" binding:"required"`
	Quantity  decimal.Decimal `json:"quantity" binding:"required"`
}

type Delivery struct {
	ID              int64           `json:"id"`
	DeliveryOrderID int64           `json:"delivery_order_id"`
	OrderNumber     string          `json:"order_number,omitempty"`
	UserID          int64           `json:"user_id"`
	UserName        string          `json:"user_name,omitempty"`
	LocationID      int64           `json:"location_id"`
	LocationName    string          `json:"location_name,omitempty"`
	DeliveredAt     time.Time       `json:"delivered_at"`
	VehicleNumber   *string         `json:"vehicle_number,omitempty"`
	Note            *string         `json:"note,omitempty"`
	Items           []DeliveryItem `json:"items,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       *time.Time      `json:"deleted_at,omitempty"`
}

type DeliveryItem struct {
	ID         int64           `json:"id"`
	DeliveryID int64           `json:"delivery_id"`
	ProductID  int64           `json:"product_id"`
	Product    *ProductSummary `json:"product,omitempty"`
	Quantity   decimal.Decimal `json:"quantity"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  *time.Time      `json:"deleted_at,omitempty"`
}
