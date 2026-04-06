-- name: CreateRole :one
INSERT INTO roles (name) VALUES ($1) RETURNING id, name, created_at, updated_at, deleted_at;

-- name: GetRole :one
SELECT id, name, created_at, updated_at, deleted_at FROM roles WHERE id = $1 LIMIT 1;

-- name: ListRoles :many
SELECT id, name, created_at, updated_at, deleted_at FROM roles ORDER BY id;

-- name: UpdateRole :one
UPDATE roles SET name = $2, updated_at = now() WHERE id = $1 RETURNING id, name, created_at, updated_at, deleted_at;

-- name: DeleteRole :exec
UPDATE roles SET deleted_at = now() WHERE id = $1;
