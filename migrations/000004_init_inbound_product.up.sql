CREATE TABLE "incoming_schedules" (
  "id" bigserial PRIMARY KEY,
  "location_id" bigint NOT NULL REFERENCES "locations" ("id"),
  "po_number" varchar(255) NOT NULL,
  "expected_date" date NOT NULL,
  "status" varchar(50) NOT NULL DEFAULT 'pending',
  "note" text,
  "received_quantity" numeric(12, 3) NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "incoming_schedule_items" (
  "id" bigserial PRIMARY KEY,
  "incoming_schedule_id" bigint NOT NULL REFERENCES "incoming_schedules" ("id"),
  "product_id" bigint NOT NULL REFERENCES "products" ("id"),
  "quantity" numeric(12, 3) NOT NULL,
  "received_quantity" numeric(12, 3) NOT NULL DEFAULT 0,
  "status" varchar(50) NOT NULL DEFAULT 'pending',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz,
  UNIQUE ("incoming_schedule_id", "product_id")
);

CREATE TABLE "product_receipts" (
  "id" bigserial PRIMARY KEY,
  "incoming_schedule_id" bigint REFERENCES "incoming_schedules" ("id"),
  "location_id" bigint NOT NULL REFERENCES "locations" ("id"),
  "received_date" date NOT NULL,
  "received_by" bigint NOT NULL REFERENCES "users" ("id"),
  "note" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "product_receipt_items" (
  "id" bigserial PRIMARY KEY,
  "product_receipt_id" bigint NOT NULL REFERENCES "product_receipts" ("id"),
  "product_id" bigint NOT NULL REFERENCES "products" ("id"),
  "quantity" numeric(12, 3) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE INDEX ON "incoming_schedules" ("location_id");
CREATE INDEX ON "incoming_schedules" ("po_number");
CREATE INDEX ON "incoming_schedule_items" ("incoming_schedule_id");
CREATE INDEX ON "incoming_schedule_items" ("product_id");
CREATE INDEX ON "product_receipts" ("incoming_schedule_id");
CREATE INDEX ON "product_receipts" ("location_id");
CREATE INDEX ON "product_receipt_items" ("product_receipt_id");
CREATE INDEX ON "product_receipt_items" ("product_id");
