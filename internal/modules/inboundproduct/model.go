package inboundproduct

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	StatusPending   = "pending"
	StatusReceived  = "received"
	StatusCancelled = "cancelled"
)

// IncomingSchedule represents a scheduled product arrival
type IncomingSchedule struct {
	ID               int64                  `json:"id"`
	LocationID       int64                  `json:"location_id"`
	PONumber         string                 `json:"po_number"`
	ExpectedDate     time.Time              `json:"expected_date"`
	Status           string                 `json:"status"`
	Note             *string                `json:"note"`
	ReceivedQuantity decimal.Decimal        `json:"received_quantity"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	DeletedAt        *time.Time             `json:"deleted_at"`
	Items            []IncomingScheduleItem `json:"items"`
}

// IncomingScheduleItem represents an item within a schedule
type IncomingScheduleItem struct {
	ID                 int64           `json:"id"`
	IncomingScheduleID int64           `json:"incoming_schedule_id"`
	ProductID          int64           `json:"product_id"`
	Product            *ProductSummary `json:"product,omitempty"`
	Quantity           decimal.Decimal `json:"quantity"`
	ReceivedQuantity   decimal.Decimal `json:"received_quantity"`
	Status             string          `json:"status"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt          *time.Time      `json:"deleted_at"`
}

type ProductSummary struct {
	Name    string `json:"name"`
	SKUCode string `json:"sku_code"`
	UOM     string `json:"uom"`
}

// ProductReceipt represents a physical receipt of products
type ProductReceipt struct {
	ID                 int64                `json:"id"`
	IncomingScheduleID *int64               `json:"incoming_schedule_id"`
	IncomingSchedule   *ScheduleSummary     `json:"incoming_schedule,omitempty"`
	LocationID         int64                `json:"location_id"`
	Location           *LocationSummary     `json:"location,omitempty"`
	ReceivedDate       time.Time            `json:"received_date"`
	ReceivedBy         int64                `json:"received_by"`
	Note               *string              `json:"note"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	DeletedAt          *time.Time           `json:"deleted_at"`
	Items              []ProductReceiptItem `json:"items"`
}

type ScheduleSummary struct {
	PONumber     string    `json:"po_number"`
	ExpectedDate time.Time `json:"expected_date"`
}

type LocationSummary struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// ProductReceiptItem represents an item received
type ProductReceiptItem struct {
	ID               int64           `json:"id"`
	ProductReceiptID int64           `json:"product_receipt_id"`
	ProductID        int64           `json:"product_id"`
	Product          *ProductSummary `json:"product,omitempty"`
	Quantity         decimal.Decimal `json:"quantity"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        *time.Time      `json:"deleted_at"`
}

// Request DTOs
type CreateScheduleRequest struct {
	LocationID   int64                   `json:"location_id" binding:"required"`
	PONumber     string                  `json:"po_number" binding:"required"`
	ExpectedDate string                  `json:"expected_date" binding:"required"` // Format: YYYY-MM-DD
	Note         string                  `json:"note"`
	Items        []CreateScheduleItemDTO `json:"items"`
}

type CreateScheduleItemDTO struct {
	ProductID int64           `json:"product_id" binding:"required"`
	Quantity  decimal.Decimal `json:"quantity" binding:"required"`
}

type UpdateScheduleRequest struct {
	LocationID   int64                   `json:"location_id" binding:"required"`
	PONumber     string                  `json:"po_number" binding:"required"`
	ExpectedDate string                  `json:"expected_date" binding:"required"`
	Status       string                  `json:"status" binding:"required"`
	Note         string                  `json:"note"`
	Items        []UpdateScheduleItemDTO `json:"items"`
}

type UpdateScheduleItemDTO struct {
	ProductID        int64           `json:"product_id" binding:"required"`
	Quantity         decimal.Decimal `json:"quantity" binding:"required"`
	ReceivedQuantity decimal.Decimal `json:"received_quantity"`
	Status           string          `json:"status"`
}

type CreateReceiptRequest struct {
	LocationID         int64                  `json:"location_id" binding:"required"`
	IncomingScheduleID *int64                 `json:"incoming_schedule_id"`
	ReceivedDate       string                 `json:"received_date" binding:"required"`
	ReceivedBy         int64                  `json:"received_by" binding:"required"`
	Note               string                 `json:"note"`
	Items              []CreateReceiptItemDTO `json:"items" binding:"required,min=1"`
}

type CreateReceiptItemDTO struct {
	ProductID int64           `json:"product_id" binding:"required"`
	Quantity  decimal.Decimal `json:"quantity" binding:"required"`
}
