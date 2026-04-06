package roles

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCreateRoleRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateRoleRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid request",
			req: CreateRoleRequest{
				Name: "Admin",
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			req: CreateRoleRequest{
				Name: "",
			},
			wantErr: true,
			errMsg:  "role name cannot be empty",
		},
		{
			name: "Only whitespace",
			req: CreateRoleRequest{
				Name: "   ",
			},
			wantErr: true,
			errMsg:  "role name cannot be empty",
		},
		{
			name: "Name too long",
			req: CreateRoleRequest{
				Name: strings.Repeat("a", 256),
			},
			wantErr: true,
			errMsg:  "role name exceeds 255 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRoleRequest(&tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
