-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2024-01-01T10:56:18.680Z

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "age" int NOT NULL,
  "sex" varchar NOT NULL,
  "image" varchar NOT NULL,
  "phone_number" varchar NOT NULL,
  "role" varchar NOT NULL,
  "is_deleted" boolean NOT NULL DEFAULT false,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "is_email_verified" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar NOT NULL,
  "secret_code" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "expired_at" timestamptz NOT NULL DEFAULT (now()  + interval '15 minutes')
);

CREATE TABLE "reset_passwords" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "reset_code" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "expired_at" timestamptz NOT NULL DEFAULT (now()  + interval '15 minutes')
);

CREATE TABLE "address" (
  "id" bigserial PRIMARY KEY,
  "address" varchar NOT NULL,
  "username" varchar NOT NULL,
  "city_id" bigserial NOT NULL,
  "district_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "cities" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "districts" (
  "id" bigserial PRIMARY KEY,
  "city_id" bigserial NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "books" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "price" float NOT NULL,
  "image" varchar[] NOT NULL,
  "description" varchar NOT NULL,
  "author" varchar NOT NULL,
  "publisher" varchar NOT NULL,
  "quantity" int NOT NULL,
  "is_deleted" boolean NOT NULL DEFAULT false,
  "sale" float NOT NULL DEFAULT 0,
  "rating" float NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "genres" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "is_deleted" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "carts" (
  "id" bigserial PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "username" varchar NOT NULL,
  "amount" int NOT NULL DEFAULT 1,
  "total" float NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "books_genres" (
  "id" bigserial PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "genres_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "reviews" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "books_id" bigserial NOT NULL,
  "comments" varchar NOT NULL,
  "is_deleted" boolean NOT NULL DEFAULT false,
  "rating" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "like" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "review_id" bigserial NOT NULL,
  "is_like" boolean NOT NULL DEFAULT false
);

CREATE TABLE "dislike" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "review_id" bigserial NOT NULL,
  "is_dislike" boolean NOT NULL DEFAULT false
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT 'unpaid',
  "sub_amount" int NOT NULL DEFAULT 1,
  "sub_total" float NOT NULL DEFAULT 0,
  "sale" float NOT NULL DEFAULT 0,
  "note" varchar,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "orders_id" bigserial NOT NULL,
  "books_id" bigserial NOT NULL,
  "amount" int NOT NULL DEFAULT 1,
  "total" float NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "wishlists" (
  "id" bigserial PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "username" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT 'false',
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "searchs" (
  "id" bigserial PRIMARY KEY,
  "book_name" varchar NOT NULL,
  "price" float NOT NULL,
  "author" varchar NOT NULL,
  "publisher" varchar NOT NULL,
  "rating" float NOT NULL,
  "genres" varchar NOT NULL,
  "subgenres" varchar NOT NULL,
  "searchs_tsv" tsvector
);

CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "order_id" bigserial NOT NULL,
  "shipping_id" bigserial NOT NULL,
  "price" float NOT NULL DEFAULT 0,
  "subtotal" float NOT NULL DEFAULT 0,
  "status" varchar NOT NULL DEFAULT 'failed',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "shippings" (
  "id" bigserial PRIMARY KEY,
  "to_address" varchar NOT NULL,
  "total" float NOT NULL DEFAULT 0
);

CREATE INDEX ON "address" ("username");

CREATE INDEX ON "address" ("city_id");

CREATE INDEX ON "address" ("district_id");

CREATE INDEX ON "address" ("username", "city_id", "district_id");

CREATE INDEX ON "districts" ("city_id");

CREATE INDEX ON "carts" ("books_id");

CREATE INDEX ON "carts" ("username");

CREATE INDEX ON "carts" ("books_id", "username");

CREATE INDEX ON "books_genres" ("books_id");

CREATE INDEX ON "books_genres" ("genres_id");

CREATE INDEX ON "books_genres" ("books_id", "genres_id");

CREATE INDEX ON "reviews" ("username");

CREATE INDEX ON "reviews" ("books_id");

CREATE INDEX ON "reviews" ("username", "books_id");

CREATE INDEX ON "like" ("username");

CREATE INDEX ON "like" ("review_id");

CREATE INDEX ON "like" ("username", "review_id");

CREATE INDEX ON "dislike" ("username");

CREATE INDEX ON "dislike" ("review_id");

CREATE INDEX ON "dislike" ("username", "review_id");

CREATE INDEX ON "orders" ("username");

CREATE INDEX ON "transactions" ("books_id");

CREATE INDEX ON "transactions" ("orders_id");

CREATE INDEX ON "transactions" ("books_id", "orders_id");

CREATE INDEX ON "wishlists" ("books_id");

CREATE INDEX ON "wishlists" ("username");

CREATE INDEX ON "wishlists" ("books_id", "username");

CREATE INDEX ON "payments" ("username");

CREATE INDEX ON "payments" ("order_id");

CREATE INDEX ON "payments" ("shipping_id");

CREATE INDEX ON "payments" ("username", "order_id", "shipping_id");

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "reset_passwords" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "address" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "address" ADD FOREIGN KEY ("city_id") REFERENCES "cities" ("id");

ALTER TABLE "address" ADD FOREIGN KEY ("district_id") REFERENCES "districts" ("id");

ALTER TABLE "districts" ADD FOREIGN KEY ("city_id") REFERENCES "cities" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "books_genres" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "books_genres" ADD FOREIGN KEY ("genres_id") REFERENCES "genres" ("id");

ALTER TABLE "reviews" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "reviews" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "like" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "like" ADD FOREIGN KEY ("review_id") REFERENCES "reviews" ("id");

ALTER TABLE "dislike" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "dislike" ADD FOREIGN KEY ("review_id") REFERENCES "reviews" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "transactions" ADD FOREIGN KEY ("orders_id") REFERENCES "orders" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "wishlists" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "wishlists" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "payments" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "orders" ADD FOREIGN KEY ("id") REFERENCES "payments" ("order_id");

ALTER TABLE "shippings" ADD FOREIGN KEY ("id") REFERENCES "payments" ("shipping_id");
