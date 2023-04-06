ALTER TABLE "carts" ADD COLUMN "amount" int NOT NULL DEFAULT 1;

ALTER TABLE "carts" ADD COLUMN "total" int NOT NULL DEFAULT 0;

ALTER TABLE "transactions" ADD COLUMN "amount" int NOT NULL DEFAULT 1;

ALTER TABLE "transactions" ADD COLUMN "total" int NOT NULL DEFAULT 0;

ALTER TABLE "orders" ADD COLUMN "status" varchar NOT NULL DEFAULT 'unpaid';

ALTER TABLE "orders" ADD COLUMN "sub_amount" int NOT NULL DEFAULT 1;

ALTER TABLE "orders" ADD COLUMN "sub_total" int NOT NULL DEFAULT 0;