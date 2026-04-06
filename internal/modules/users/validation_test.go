package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCreateUser(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateUserRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid request",
			req: CreateUserRequest{
				Username: "jdoe",
				NIK:      "1234567890",
				Password: "securepassword",
				FullName: "John Doe",
				RoleID:   1,
			},
			wantErr: false,
		},
		{
			name: "Empty Username",
			req: CreateUserRequest{
				Username: "",
				NIK:      "1234567890",
				Password: "securepassword",
				FullName: "John Doe",
				RoleID:   1,
			},
			wantErr: true,
			errMsg:  "username cannot be empty",
		},
		{
			name: "NIK too short",
			req: CreateUserRequest{
				Username: "jdoe",
				NIK:      "123",
				Password: "securepassword",
				FullName: "John Doe",
				RoleID:   1,
			},
			wantErr: true,
			errMsg:  "NIK must be exactly 10 characters",
		},
		{
			name: "NIK too long",
			req: CreateUserRequest{
				Username: "jdoe",
				NIK:      "1234567890123",
				Password: "securepassword",
				FullName: "John Doe",
				RoleID:   1,
			},
			wantErr: true,
			errMsg:  "NIK must be exactly 10 characters",
		},
		{
			name: "Password too short",
			req: CreateUserRequest{
				Username: "jdoe",
				NIK:      "1234567890",
				Password: "short",
				FullName: "John Doe",
				RoleID:   1,
			},
			wantErr: true,
			errMsg:  "password must be at least 8 characters long",
		},
		{
			name: "Empty FullName",
			req: CreateUserRequest{
				Username: "jdoe",
				NIK:      "1234567890",
				Password: "securepassword",
				FullName: "",
				RoleID:   1,
			},
			wantErr: true,
			errMsg:  "full name cannot be empty",
		},
		{
			name: "Invalid RoleID",
			req: CreateUserRequest{
				Username: "jdoe",
				NIK:      "1234567890",
				Password: "securepassword",
				FullName: "John Doe",
				RoleID:   0,
			},
			wantErr: true,
			errMsg:  "invalid role id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateUser(&tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
