CREATE TABLE "payments" (
  "id" varchar PRIMARY KEY,
  "username" varchar NOT NULL,
  "order_id" bigserial NOT NULL,
  "shipping_id" bigserial NOT NULL,
  "subtotal" float NOT NULL DEFAULT 0,
  "status" varchar NOT NULL DEFAULT 'failed',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "shippings" (
  "id" bigserial PRIMARY KEY,
  "to_address" varchar NOT NULL,
  "total" float NOT NULL DEFAULT 0
);

CREATE INDEX ON "payments" ("username");

CREATE INDEX ON "payments" ("order_id");

CREATE INDEX ON "payments" ("shipping_id");

CREATE INDEX ON "payments" ("username", "order_id", "shipping_id");

ALTER TABLE "payments" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "payments" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("shipping_id") REFERENCES "shippings" ("id");