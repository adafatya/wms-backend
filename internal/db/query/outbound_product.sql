-- name: CreateCustomer :one
INSERT INTO customers (
  name, address, contact_name, contact_info
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListCustomers :many
SELECT * FROM customers
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountCustomers :one
SELECT COUNT(*) FROM customers
WHERE deleted_at IS NULL;

-- name: UpdateCustomer :one
UPDATE customers
SET name = $2, address = $3, contact_name = $4, contact_info = $5, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteCustomer :execresult
UPDATE customers
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateDeliveryOrder :one
INSERT INTO delivery_orders (
  customer_id, order_number, delivery_date, note
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDeliveryOrder :one
SELECT dlo.*, c.name as customer_name, c.address as customer_address
FROM delivery_orders dlo
JOIN customers c ON dlo.customer_id = c.id
WHERE dlo.id = $1 AND dlo.deleted_at IS NULL LIMIT 1;

-- name: ListDeliveryOrders :many
SELECT dlo.*, c.name as customer_name
FROM delivery_orders dlo
JOIN customers c ON dlo.customer_id = c.id
WHERE dlo.deleted_at IS NULL
ORDER BY dlo.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountDeliveryOrders :one
SELECT COUNT(*) FROM delivery_orders
WHERE deleted_at IS NULL;

-- name: UpdateDeliveryOrder :one
UPDATE delivery_orders
SET customer_id = $2, order_number = $3, delivery_date = $4, status = $5, note = $6, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateDeliveryOrderStatus :one
UPDATE delivery_orders
SET status = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDeliveryOrder :execresult
UPDATE delivery_orders
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: BulkCreateDeliveryOrderItems :exec
INSERT INTO delivery_order_items (
  delivery_order_id, product_id, quantity
)
SELECT $1, unnest($2::bigint[]), unnest($3::numeric[]);

-- name: GetDeliveryOrderItems :many
SELECT doi.*, p.name as product_name, p.sku_code as product_sku_code, p.uom as product_uom
FROM delivery_order_items doi
JOIN products p ON doi.product_id = p.id
WHERE doi.delivery_order_id = $1 AND doi.deleted_at IS NULL;

-- name: UpsertDeliveryOrderItem :one
INSERT INTO delivery_order_items (
  delivery_order_id, product_id, quantity, delivered_quantity
) VALUES ($1, $2, $3, $4)
ON CONFLICT (delivery_order_id, product_id) DO UPDATE SET
  quantity = EXCLUDED.quantity,
  delivered_quantity = EXCLUDED.delivered_quantity,
  updated_at = now(),
  deleted_at = NULL
RETURNING *;

-- name: DeleteDeliveryOrderItems :exec
UPDATE delivery_order_items
SET deleted_at = now()
WHERE delivery_order_id = $1 AND deleted_at IS NULL;

-- name: IncrementDeliveryOrderItemDeliveredQty :exec
UPDATE delivery_order_items
SET delivered_quantity = delivered_quantity + $3,
    updated_at = now()
WHERE delivery_order_id = $1 AND product_id = $2 AND deleted_at IS NULL;

-- name: CreateDelivery :one
INSERT INTO deliveries (
  delivery_order_id, user_id, location_id, delivered_at, vehicle_number, note
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetDelivery :one
SELECT d.*, dlo.order_number, u.username as user_name, l.name as location_name
FROM deliveries d
JOIN delivery_orders dlo ON d.delivery_order_id = dlo.id
JOIN users u ON d.user_id = u.id
JOIN locations l ON d.location_id = l.id
WHERE d.id = $1 AND d.deleted_at IS NULL LIMIT 1;

-- name: ListDeliveries :many
SELECT d.*, dlo.order_number
FROM deliveries d
JOIN delivery_orders dlo ON d.delivery_order_id = dlo.id
WHERE d.deleted_at IS NULL
ORDER BY d.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountDeliveries :one
SELECT COUNT(*) FROM deliveries
WHERE deleted_at IS NULL;

-- name: UpdateDelivery :one
UPDATE deliveries
SET delivery_order_id = $2, user_id = $3, location_id = $4, delivered_at = $5, vehicle_number = $6, note = $7, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDelivery :execresult
UPDATE deliveries
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: BulkCreateDeliveryItems :exec
INSERT INTO delivery_items (
  delivery_id, product_id, quantity
)
SELECT $1, unnest($2::bigint[]), unnest($3::numeric[]);

-- name: GetDeliveryItems :many
SELECT di.*, p.name as product_name, p.sku_code as product_sku_code, p.uom as product_uom
FROM delivery_items di
JOIN products p ON di.product_id = p.id
WHERE di.delivery_id = $1 AND di.deleted_at IS NULL;
