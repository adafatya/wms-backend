package inboundproduct

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateSchedule(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateScheduleRequest
		wantErr bool
		msg     string
	}{
		{
			name: "Valid Request",
			req: CreateScheduleRequest{
				LocationID:   1,
				PONumber:     "PO-001",
				ExpectedDate: "2024-12-31",
				Items: []CreateScheduleItemDTO{
					{ProductID: 1, Quantity: decimal.NewFromInt(10)},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid LocationID",
			req: CreateScheduleRequest{
				LocationID:   0,
				PONumber:     "PO-001",
				ExpectedDate: "2024-12-31",
			},
			wantErr: true,
			msg:     "location_id must be greater than 0",
		},
		{
			name: "Empty PONumber",
			req: CreateScheduleRequest{
				LocationID:   1,
				PONumber:     "",
				ExpectedDate: "2024-12-31",
			},
			wantErr: true,
			msg:     "po_number cannot be empty",
		},
		{
			name: "Empty ExpectedDate",
			req: CreateScheduleRequest{
				LocationID:   1,
				PONumber:     "PO-001",
				ExpectedDate: "",
			},
			wantErr: true,
			msg:     "expected_date cannot be empty",
		},
		{
			name: "Invalid ExpectedDate Format",
			req: CreateScheduleRequest{
				LocationID:   1,
				PONumber:     "PO-001",
				ExpectedDate: "31-12-2024",
			},
			wantErr: true,
			msg:     "invalid expected_date format",
		},
		{
			name: "Invalid ProductID in Items",
			req: CreateScheduleRequest{
				LocationID:   1,
				PONumber:     "PO-001",
				ExpectedDate: "2024-12-31",
				Items: []CreateScheduleItemDTO{
					{ProductID: 0, Quantity: decimal.NewFromInt(10)},
				},
			},
			wantErr: true,
			msg:     "product_id must be greater than 0",
		},
		{
			name: "Negative Quantity in Items",
			req: CreateScheduleRequest{
				LocationID:   1,
				PONumber:     "PO-001",
				ExpectedDate: "2024-12-31",
				Items: []CreateScheduleItemDTO{
					{ProductID: 1, Quantity: decimal.NewFromInt(-1)},
				},
			},
			wantErr: true,
			msg:     "quantity cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateSchedule(&tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.msg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUpdateSchedule(t *testing.T) {
	tests := []struct {
		name    string
		req     UpdateScheduleRequest
		wantErr bool
		msg     string
	}{
		{
			name: "Valid Request",
			req: UpdateScheduleRequest{
				LocationID:   1,
				PONumber:     "PO-001",
				ExpectedDate: "2024-12-31",
				Status:       StatusPending,
				Items: []UpdateScheduleItemDTO{
					{ProductID: 1, Quantity: decimal.NewFromInt(10), ReceivedQuantity: decimal.NewFromInt(0), Status: StatusPending},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Status",
			req: UpdateScheduleRequest{
				LocationID:   1,
				PONumber:     "PO-001",
				ExpectedDate: "2024-12-31",
				Status:       "unknown",
			},
			wantErr: true,
			msg:     "invalid status",
		},
		{
			name: "Negative ReceivedQuantity",
			req: UpdateScheduleRequest{
				LocationID:   1,
				PONumber:     "PO-001",
				ExpectedDate: "2024-12-31",
				Status:       StatusPending,
				Items: []UpdateScheduleItemDTO{
					{ProductID: 1, Quantity: decimal.NewFromInt(10), ReceivedQuantity: decimal.NewFromInt(-1)},
				},
			},
			wantErr: true,
			msg:     "received_quantity cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateSchedule(&tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.msg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCreateReceipt(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateReceiptRequest
		wantErr bool
		msg     string
	}{
		{
			name: "Valid Request",
			req: CreateReceiptRequest{
				LocationID:   1,
				ReceivedDate: "2024-12-31",
				ReceivedBy:   1,
				Items: []CreateReceiptItemDTO{
					{ProductID: 1, Quantity: decimal.NewFromInt(10)},
				},
			},
			wantErr: false,
		},
		{
			name: "Empty Items",
			req: CreateReceiptRequest{
				LocationID:   1,
				ReceivedDate: "2024-12-31",
				ReceivedBy:   1,
				Items:        []CreateReceiptItemDTO{},
			},
			wantErr: true,
			msg:     "items cannot be empty",
		},
		{
			name: "Zero Quantity",
			req: CreateReceiptRequest{
				LocationID:   1,
				ReceivedDate: "2024-12-31",
				ReceivedBy:   1,
				Items: []CreateReceiptItemDTO{
					{ProductID: 1, Quantity: decimal.NewFromInt(0)},
				},
			},
			wantErr: true,
			msg:     "quantity must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateReceipt(&tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.msg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
