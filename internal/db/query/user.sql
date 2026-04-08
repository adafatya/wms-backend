-- name: CreateUser :one
INSERT INTO users (username, nik, password, full_name, role_id) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT * FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;
