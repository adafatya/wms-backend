package roles

import (
	"errors"
	"strings"
)

// ValidateCreateRoleRequest validates the role name
func ValidateCreateRoleRequest(req *CreateRoleRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("role name cannot be empty")
	}
	if len(req.Name) > 255 {
		return errors.New("role name exceeds 255 characters")
	}
	return nil
}

// ValidateUpdateRoleRequest validates the updated role name
func ValidateUpdateRoleRequest(req *UpdateRoleRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("role name cannot be empty")
	}
	if len(req.Name) > 255 {
		return errors.New("role name exceeds 255 characters")
	}
	return nil
}
