-- name: CreateProduct :one
INSERT INTO products (name, sku_code, uom)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE deleted_at IS NULL
  AND (name ILIKE $1 OR sku_code ILIKE $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountProducts :one
SELECT COUNT(*) FROM products
WHERE deleted_at IS NULL
  AND (name ILIKE $1 OR sku_code ILIKE $1);

-- name: UpdateProduct :one
UPDATE products
SET name = $2, sku_code = $3, uom = $4, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteProduct :exec
UPDATE products
SET deleted_at = now()
WHERE id = $1;

-- name: CreateLocation :one
INSERT INTO locations (name, code)
VALUES ($1, $2)
RETURNING *;

-- name: GetLocation :one
SELECT * FROM locations
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListLocations :many
SELECT * FROM locations
WHERE deleted_at IS NULL
  AND (name ILIKE $1 OR code ILIKE $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountLocations :one
SELECT COUNT(*) FROM locations
WHERE deleted_at IS NULL
  AND (name ILIKE $1 OR code ILIKE $1);

-- name: UpdateLocation :one
UPDATE locations
SET name = $2, code = $3, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteLocation :exec
UPDATE locations
SET deleted_at = now()
WHERE id = $1;

-- name: BulkUpsertInventories :exec
INSERT INTO inventories (product_id, location_id, quantity, updated_at)
SELECT unnest(@product_ids::bigint[]), unnest(@location_ids::bigint[]), unnest(@quantities::numeric[]), now()
ON CONFLICT (product_id, location_id) DO UPDATE SET
  quantity = EXCLUDED.quantity,
  updated_at = now();

-- name: GetInventoriesByLocation :many
SELECT i.*, p.name as product_name, p.sku_code as product_sku_code, p.uom as product_uom
FROM inventories i
JOIN products p ON i.product_id = p.id
WHERE i.location_id = $1 AND i.deleted_at IS NULL AND p.deleted_at IS NULL;

-- name: GetInventoriesByProduct :many
SELECT i.*, l.name as location_name, l.code as location_code
FROM inventories i
JOIN locations l ON i.location_id = l.id
WHERE i.product_id = $1 AND i.deleted_at IS NULL AND l.deleted_at IS NULL;

-- name: ListInventories :many
SELECT i.*, 
       p.name as product_name, p.sku_code as product_sku_code, p.uom as product_uom,
       l.name as location_name, l.code as location_code
FROM inventories i
JOIN products p ON i.product_id = p.id
JOIN locations l ON i.location_id = l.id
WHERE i.deleted_at IS NULL AND p.deleted_at IS NULL AND l.deleted_at IS NULL
ORDER BY i.updated_at DESC
LIMIT $1 OFFSET $2;

-- name: CountInventories :one
SELECT COUNT(*) FROM inventories i
JOIN products p ON i.product_id = p.id
JOIN locations l ON i.location_id = l.id
WHERE i.deleted_at IS NULL AND p.deleted_at IS NULL AND l.deleted_at IS NULL;

-- name: GetInventoryStock :one
SELECT quantity FROM inventories
WHERE product_id = $1 AND location_id = $2 AND deleted_at IS NULL;
