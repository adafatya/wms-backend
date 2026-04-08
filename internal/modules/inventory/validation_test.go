package inventory

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestValidation_isValidCode(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"Valid Code", "SKU-123-ABC", true},
		{"Valid Code Numbers Only", "123456", true},
		{"Valid Code Letters Only", "ABC-DEF", true},
		{"Invalid lowercase", "abc-123", false},
		{"Invalid space", "SKU 123", false},
		{"Invalid symbol", "SKU@123", false},
		{"Empty string", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isValidCode(tt.code))
		})
	}
}

func TestValidateCreateProduct(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateProductRequest
		wantErr bool
		msg     string
	}{
		{
			name: "Valid Request",
			req: CreateProductRequest{
				Name:    "Product A",
				SKUCode: "PROD-A",
				UOM:     "PCS",
			},
			wantErr: false,
		},
		{
			name: "Empty Name",
			req: CreateProductRequest{
				Name:    " ",
				SKUCode: "PROD-A",
				UOM:     "PCS",
			},
			wantErr: true,
			msg:     "product name cannot be empty",
		},
		{
			name: "Invalid SKU",
			req: CreateProductRequest{
				Name:    "Product A",
				SKUCode: "prod-a",
				UOM:     "PCS",
			},
			wantErr: true,
			msg:     "sku code can only contain capital letters, numbers, and dashes",
		},
		{
			name: "Negative Quantity",
			req: CreateProductRequest{
				Name:    "Product A",
				SKUCode: "PROD-A",
				UOM:     "PCS",
				Inventory: []InventoryInput{
					{LocationID: 1, Quantity: decimal.NewFromInt(-1)},
				},
			},
			wantErr: true,
			msg:     "quantity cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateProduct(&tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.msg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCreateLocation(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateLocationRequest
		wantErr bool
		msg     string
	}{
		{
			name: "Valid Request",
			req: CreateLocationRequest{
				Name: "Warehouse A",
				Code: "WH-A",
			},
			wantErr: false,
		},
		{
			name: "Invalid Code",
			req: CreateLocationRequest{
				Name: "Warehouse A",
				Code: "wh a",
			},
			wantErr: true,
			msg:     "location code can only contain capital letters, numbers, and dashes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateLocation(&tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.msg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
