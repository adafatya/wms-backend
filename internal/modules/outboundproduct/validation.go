package outboundproduct

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

func ValidateCreateCustomer(req CreateCustomerRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Address == "" {
		return fmt.Errorf("address is required")
	}
	return nil
}

func ValidateUpdateCustomer(req UpdateCustomerRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Address == "" {
		return fmt.Errorf("address is required")
	}
	return nil
}

func ValidateCreateDeliveryOrder(req CreateDeliveryOrderRequest) error {
	if req.CustomerID <= 0 {
		return fmt.Errorf("customer_id is required")
	}
	if req.OrderNumber == "" {
		return fmt.Errorf("order_number is required")
	}
	if _, err := time.Parse("2006-01-02", req.DeliveryDate); err != nil {
		return fmt.Errorf("invalid delivery_date format, use YYYY-MM-DD")
	}
	if len(req.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return fmt.Errorf("invalid product_id")
		}
		if item.Quantity.LessThanOrEqual(decimal.Zero) {
			return fmt.Errorf("quantity must be greater than zero")
		}
	}
	return nil
}

func ValidateUpdateDeliveryOrder(req UpdateDeliveryOrderRequest) error {
	if req.CustomerID <= 0 {
		return fmt.Errorf("customer_id is required")
	}
	if req.OrderNumber == "" {
		return fmt.Errorf("order_number is required")
	}
	if _, err := time.Parse("2006-01-02", req.DeliveryDate); err != nil {
		return fmt.Errorf("invalid delivery_date format, use YYYY-MM-DD")
	}
	if req.Status == "" {
		return fmt.Errorf("status is required")
	}

	validStatuses := map[string]bool{
		StatusPending:    true,
		StatusProcessing: true,
		StatusShipped:    true,
		StatusCancelled:  true,
		StatusCompleted:  true,
	}
	if !validStatuses[req.Status] {
		return fmt.Errorf("invalid status: %s", req.Status)
	}
	if len(req.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return fmt.Errorf("invalid product_id")
		}
		if item.Quantity.LessThanOrEqual(decimal.Zero) {
			return fmt.Errorf("quantity must be greater than zero")
		}
	}
	return nil
}

func ValidateCreateDelivery(req CreateDeliveryRequest) error {
	if req.DeliveryOrderID <= 0 {
		return fmt.Errorf("delivery_order_id is required")
	}
	if req.LocationID <= 0 {
		return fmt.Errorf("location_id is required")
	}
	if req.UserID <= 0 {
		return fmt.Errorf("user_id is required")
	}
	if _, err := time.Parse("2006-01-02 15:04:05", req.DeliveredAt); err != nil {
		return fmt.Errorf("invalid delivered_at format, use YYYY-MM-DD HH:mm:ss")
	}
	if len(req.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return fmt.Errorf("invalid product_id")
		}
		if item.Quantity.LessThanOrEqual(decimal.Zero) {
			return fmt.Errorf("quantity must be greater than zero")
		}
	}
	return nil
}

func ValidateUpdateDelivery(req UpdateDeliveryRequest) error {
	if req.DeliveryOrderID <= 0 {
		return fmt.Errorf("delivery_order_id is required")
	}
	if req.LocationID <= 0 {
		return fmt.Errorf("location_id is required")
	}
	if req.UserID <= 0 {
		return fmt.Errorf("user_id is required")
	}
	if _, err := time.Parse("2006-01-02 15:04:05", req.DeliveredAt); err != nil {
		return fmt.Errorf("invalid delivered_at format, use YYYY-MM-DD HH:mm:ss")
	}
	if len(req.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return fmt.Errorf("invalid product_id")
		}
		if item.Quantity.LessThanOrEqual(decimal.Zero) {
			return fmt.Errorf("quantity must be greater than zero")
		}
	}
	return nil
}
