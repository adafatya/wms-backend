package inventory

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	codeRegex = regexp.MustCompile(`^[A-Z0-9-]+$`)
)

func isValidCode(code string) bool {
	return codeRegex.MatchString(code)
}

// Product Validations
func ValidateCreateProduct(req *CreateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("product name cannot be empty")
	}
	if strings.TrimSpace(req.SKUCode) == "" {
		return errors.New("sku code cannot be empty")
	}
	if !isValidCode(req.SKUCode) {
		return errors.New("sku code can only contain capital letters, numbers, and dashes")
	}
	if strings.TrimSpace(req.UOM) == "" {
		return errors.New("uom cannot be empty")
	}

	for _, inv := range req.Inventory {
		if err := ValidateInventoryInput(&inv); err != nil {
			return err
		}
		if inv.LocationID <= 0 {
			return errors.New("invalid location id in inventory")
		}
	}

	return nil
}

func ValidateUpdateProduct(req *UpdateProductRequest) error {
	if req.Name != "" && strings.TrimSpace(req.Name) == "" {
		return errors.New("product name cannot be empty")
	}
	if req.SKUCode != "" {
		if strings.TrimSpace(req.SKUCode) == "" {
			return errors.New("sku code cannot be empty")
		}
		if !isValidCode(req.SKUCode) {
			return errors.New("sku code can only contain capital letters, numbers, and dashes")
		}
	}
	if req.UOM != "" && strings.TrimSpace(req.UOM) == "" {
		return errors.New("uom cannot be empty")
	}

	for _, inv := range req.Inventory {
		if err := ValidateInventoryInput(&inv); err != nil {
			return err
		}
		if inv.LocationID <= 0 {
			return errors.New("invalid location id in inventory")
		}
	}

	return nil
}

// Location Validations
func ValidateCreateLocation(req *CreateLocationRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("location name cannot be empty")
	}
	if strings.TrimSpace(req.Code) == "" {
		return errors.New("location code cannot be empty")
	}
	if !isValidCode(req.Code) {
		return errors.New("location code can only contain capital letters, numbers, and dashes")
	}

	for _, inv := range req.Inventory {
		if err := ValidateInventoryInput(&inv); err != nil {
			return err
		}
		if inv.ProductID <= 0 {
			return errors.New("invalid product id in inventory")
		}
	}

	return nil
}

func ValidateUpdateLocation(req *UpdateLocationRequest) error {
	if req.Name != "" && strings.TrimSpace(req.Name) == "" {
		return errors.New("location name cannot be empty")
	}
	if req.Code != "" {
		if strings.TrimSpace(req.Code) == "" {
			return errors.New("location code cannot be empty")
		}
		if !isValidCode(req.Code) {
			return errors.New("location code can only contain capital letters, numbers, and dashes")
		}
	}

	for _, inv := range req.Inventory {
		if err := ValidateInventoryInput(&inv); err != nil {
			return err
		}
		if inv.ProductID <= 0 {
			return errors.New("invalid product id in inventory")
		}
	}

	return nil
}

// Inventory Validations
func ValidateInventoryInput(req *InventoryInput) error {
	if req.Quantity.LessThan(decimal.Zero) {
		return fmt.Errorf("quantity cannot be negative: %s", req.Quantity.String())
	}
	return nil
}

func ValidateUpsertInventories(req []InventoryInput) error {
	if len(req) == 0 {
		return errors.New("inventory data cannot be empty")
	}
	for _, inv := range req {
		if inv.ProductID <= 0 {
			return errors.New("invalid product id")
		}
		if inv.LocationID <= 0 {
			return errors.New("invalid location id")
		}
		if err := ValidateInventoryInput(&inv); err != nil {
			return err
		}
	}
	return nil
}
