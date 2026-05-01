package outboundproduct

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateCustomer(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		req := CreateCustomerRequest{
			Name:    "Customer A",
			Address: "Address A",
		}
		err := ValidateCreateCustomer(req)
		assert.NoError(t, err)
	})

	t.Run("Empty Name", func(t *testing.T) {
		req := CreateCustomerRequest{
			Address: "Address A",
		}
		err := ValidateCreateCustomer(req)
		assert.Error(t, err)
	})
}

func TestValidateCreateDeliveryOrder(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		req := CreateDeliveryOrderRequest{
			CustomerID:   1,
			OrderNumber:  "DO-001",
			DeliveryDate: "2024-05-20",
			Items: []CreateDeliveryOrderItem{
				{ProductID: 1, Quantity: decimal.NewFromInt(10)},
			},
		}
		err := ValidateCreateDeliveryOrder(req)
		assert.NoError(t, err)
	})

	t.Run("Invalid Date", func(t *testing.T) {
		req := CreateDeliveryOrderRequest{
			CustomerID:   1,
			OrderNumber:  "DO-001",
			DeliveryDate: "20-05-2024",
			Items: []CreateDeliveryOrderItem{
				{ProductID: 1, Quantity: decimal.NewFromInt(10)},
			},
		}
		err := ValidateCreateDeliveryOrder(req)
		assert.Error(t, err)
	})

	t.Run("Empty Items", func(t *testing.T) {
		req := CreateDeliveryOrderRequest{
			CustomerID:   1,
			OrderNumber:  "DO-001",
			DeliveryDate: "2024-05-20",
			Items:        []CreateDeliveryOrderItem{},
		}
		err := ValidateCreateDeliveryOrder(req)
		assert.Error(t, err)
	})
}

func TestValidateCreateDelivery(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		req := CreateDeliveryRequest{
			DeliveryOrderID: 1,
			LocationID:      1,
			UserID:          1,
			DeliveredAt:     "2024-05-20 10:00:00",
			Items: []CreateDeliveryItemDTO{
				{ProductID: 1, Quantity: decimal.NewFromInt(5)},
			},
		}
		err := ValidateCreateDelivery(req)
		assert.NoError(t, err)
	})

	t.Run("Invalid DateTime", func(t *testing.T) {
		req := CreateDeliveryRequest{
			DeliveryOrderID: 1,
			LocationID:      1,
			UserID:          1,
			DeliveredAt:     "2024-05-20",
			Items: []CreateDeliveryItemDTO{
				{ProductID: 1, Quantity: decimal.NewFromInt(5)},
			},
		}
		err := ValidateCreateDelivery(req)
		assert.Error(t, err)
	})
}
