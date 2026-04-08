-- name: CreateRole :one
INSERT INTO roles (name) VALUES ($1) RETURNING id, name, created_at, updated_at, deleted_at;

-- name: GetRole :one
SELECT id, name, created_at, updated_at, deleted_at FROM roles WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListRoles :many
SELECT id, name, created_at, updated_at, deleted_at FROM roles WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2;

-- name: CountRoles :one
SELECT COUNT(*) FROM roles WHERE deleted_at IS NULL;

-- name: UpdateRole :one
UPDATE roles SET name = $2, updated_at = now() WHERE id = $1 RETURNING id, name, created_at, updated_at, deleted_at;

-- name: DeleteRole :exec
UPDATE roles SET deleted_at = now() WHERE id = $1;
