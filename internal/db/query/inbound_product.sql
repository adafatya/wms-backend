-- name: CreateIncomingSchedule :one
INSERT INTO incoming_schedules (
  location_id, po_number, expected_date, note
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetIncomingSchedule :one
SELECT * FROM incoming_schedules
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListIncomingSchedules :many
SELECT * FROM incoming_schedules
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountIncomingSchedules :one
SELECT COUNT(*) FROM incoming_schedules
WHERE deleted_at IS NULL;

-- name: UpdateIncomingSchedule :one
UPDATE incoming_schedules
SET location_id = $2, po_number = $3, expected_date = $4, status = $5, note = $6, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteIncomingSchedule :execresult
UPDATE incoming_schedules
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpsertIncomingScheduleItem :one
INSERT INTO incoming_schedule_items (
  incoming_schedule_id, product_id, quantity, received_quantity, status
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (incoming_schedule_id, product_id) DO UPDATE SET
  quantity = EXCLUDED.quantity,
  received_quantity = EXCLUDED.received_quantity,
  status = EXCLUDED.status,
  updated_at = now()
RETURNING *;

-- name: GetIncomingScheduleItems :many
SELECT i.*, p.name as product_name, p.sku_code as product_sku_code, p.uom as product_uom
FROM incoming_schedule_items i
JOIN products p ON i.product_id = p.id
WHERE i.incoming_schedule_id = $1 AND i.deleted_at IS NULL;

-- name: DeleteIncomingScheduleItems :exec
UPDATE incoming_schedule_items
SET deleted_at = now()
WHERE incoming_schedule_id = $1 AND deleted_at IS NULL;

-- name: CreateProductReceipt :one
INSERT INTO product_receipts (
  incoming_schedule_id, location_id, received_date, received_by, note
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProductReceipt :one
SELECT r.*, 
       l.name as location_name, l.code as location_code,
       s.po_number as schedule_po_number, s.expected_date as schedule_expected_date
FROM product_receipts r
LEFT JOIN locations l ON r.location_id = l.id
LEFT JOIN incoming_schedules s ON r.incoming_schedule_id = s.id
WHERE r.id = $1 AND r.deleted_at IS NULL LIMIT 1;

-- name: ListProductReceipts :many
SELECT r.*, 
       l.name as location_name, l.code as location_code,
       s.po_number as schedule_po_number, s.expected_date as schedule_expected_date
FROM product_receipts r
LEFT JOIN locations l ON r.location_id = l.id
LEFT JOIN incoming_schedules s ON r.incoming_schedule_id = s.id
WHERE r.deleted_at IS NULL
ORDER BY r.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProductReceipts :one
SELECT COUNT(*) FROM product_receipts
WHERE deleted_at IS NULL;

-- name: BulkCreateProductReceiptItems :exec
INSERT INTO product_receipt_items (
  product_receipt_id, product_id, quantity
)
SELECT $1, unnest(@product_ids::bigint[]), unnest(@quantities::numeric[])
RETURNING *;

-- name: GetProductReceiptItems :many
SELECT i.*, p.name as product_name, p.sku_code as product_sku_code, p.uom as product_uom
FROM product_receipt_items i
JOIN products p ON i.product_id = p.id
WHERE i.product_receipt_id = $1 AND i.deleted_at IS NULL;
