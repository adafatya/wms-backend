-- name: CreateUser :one
INSERT INTO users (username, nik, password, full_name, role_id) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;
