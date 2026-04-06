package users

import (
	"errors"
	"strings"
)

// ValidateCreateUser validates the user creation request
func ValidateCreateUser(req *CreateUserRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(req.NIK) != 10 {
		return errors.New("NIK must be exactly 10 characters")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(req.Password) > 72 {
		return errors.New("password must not exceed 72 characters")
	}
	if strings.TrimSpace(req.FullName) == "" {
		return errors.New("full name cannot be empty")
	}
	if req.RoleID <= 0 {
		return errors.New("invalid role id")
	}
	return nil
}
