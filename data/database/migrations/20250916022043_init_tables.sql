-- +goose Up
-- +goose StatementBegin
CREATE TABLE "accounts" (
  "id" serial PRIMARY KEY NOT NULL,
  "name" varchar(100) NOT NULL,
  "email" varchar(100) NOT NULL, -- TODO: Add indexing
  "password" bytea NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "tasks" (
  "id" serial PRIMARY KEY NOT NULL,
  "account_id" int NOT NULL,
  "title" varchar(100) NOT NULL,
  "description" varchar,
  "status" varchar(10) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deleted_at" timestamp
);

ALTER TABLE "tasks" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "tasks";
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
