CREATE TABLE "roles" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

-- Insert a default role if needed, or just let users handle it.
-- For now, we just create the schema.

ALTER TABLE "users" ADD COLUMN "nik" varchar(10);
ALTER TABLE "users" ADD COLUMN "password" varchar(255);
ALTER TABLE "users" ADD COLUMN "full_name" varchar(255);
ALTER TABLE "users" ADD COLUMN "role_id" bigint;
ALTER TABLE "users" ADD COLUMN "updated_at" timestamptz NOT NULL DEFAULT (now());
ALTER TABLE "users" ADD COLUMN "deleted_at" timestamptz;

-- Adding Unique Index with Partial Filter
CREATE UNIQUE INDEX users_nik_idx ON users (nik) WHERE deleted_at IS NULL;

-- Adding Foreign Key
ALTER TABLE "users" ADD CONSTRAINT "users_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles" ("id");
