package users

import "time"

type User struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	NIK       string     `json:"nik"`
	FullName  string     `json:"full_name"`
	RoleID    int64      `json:"role_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	NIK      string `json:"nik" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	RoleID   int64  `json:"role_id" binding:"required"`
}
