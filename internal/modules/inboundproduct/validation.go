package inboundproduct

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func ValidateCreateSchedule(req *CreateScheduleRequest) error {
	if req.LocationID <= 0 {
		return errors.New("location_id must be greater than 0")
	}
	if strings.TrimSpace(req.PONumber) == "" {
		return errors.New("po_number cannot be empty")
	}
	if req.ExpectedDate == "" {
		return errors.New("expected_date cannot be empty")
	}
	if _, err := time.Parse("2006-01-02", req.ExpectedDate); err != nil {
		return fmt.Errorf("invalid expected_date format, use YYYY-MM-DD: %w", err)
	}

	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return errors.New("product_id must be greater than 0")
		}
		if item.Quantity.IsNegative() {
			return errors.New("quantity cannot be negative")
		}
	}
	return nil
}

func ValidateUpdateSchedule(req *UpdateScheduleRequest) error {
	if req.LocationID <= 0 {
		return errors.New("location_id must be greater than 0")
	}
	if strings.TrimSpace(req.PONumber) == "" {
		return errors.New("po_number cannot be empty")
	}
	if req.ExpectedDate == "" {
		return errors.New("expected_date cannot be empty")
	}
	if _, err := time.Parse("2006-01-02", req.ExpectedDate); err != nil {
		return fmt.Errorf("invalid expected_date format, use YYYY-MM-DD: %w", err)
	}
	
	status := strings.ToLower(req.Status)
	if status != StatusPending && status != StatusReceived && status != StatusCancelled {
		return errors.New("invalid status")
	}

	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return errors.New("product_id must be greater than 0")
		}
		if item.Quantity.IsNegative() {
			return errors.New("quantity cannot be negative")
		}
		if item.ReceivedQuantity.IsNegative() {
			return errors.New("received_quantity cannot be negative")
		}
	}
	return nil
}

func ValidateCreateReceipt(req *CreateReceiptRequest) error {
	if req.LocationID <= 0 {
		return errors.New("location_id must be greater than 0")
	}
	if req.ReceivedDate == "" {
		return errors.New("received_date cannot be empty")
	}
	if _, err := time.Parse("2006-01-02", req.ReceivedDate); err != nil {
		return fmt.Errorf("invalid received_date format, use YYYY-MM-DD: %w", err)
	}
	if req.ReceivedBy <= 0 {
		return errors.New("received_by must be greater than 0")
	}
	if len(req.Items) == 0 {
		return errors.New("items cannot be empty")
	}

	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return errors.New("product_id must be greater than 0")
		}
		if item.Quantity.IsNegative() || item.Quantity.IsZero() {
			return errors.New("quantity must be greater than 0")
		}
	}
	return nil
}
