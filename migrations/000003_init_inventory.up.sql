CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "sku_code" varchar(255) UNIQUE NOT NULL,
  "uom" varchar(50) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "locations" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "code" varchar(255) UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "inventories" (
  "product_id" bigint NOT NULL REFERENCES "products" ("id"),
  "location_id" bigint NOT NULL REFERENCES "locations" ("id"),
  "quantity" numeric(12, 3) NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz,
  PRIMARY KEY ("product_id", "location_id")
);

CREATE INDEX ON "inventories" ("product_id");
CREATE INDEX ON "inventories" ("location_id");
