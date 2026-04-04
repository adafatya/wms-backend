CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
