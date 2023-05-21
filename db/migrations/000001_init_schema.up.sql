CREATE TABLE "users" (
  "username" varchar PRIMARY KEY NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "image" varchar NOT NULL,
  "phone_number" varchar NOT NULL,
  "age" int NOT NULL,
  "sex" varchar NOT NULL,
  "role" varchar NOT NULL DEFAULT 'user',
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "address" (
  "id" BIGSERIAL PRIMARY KEY,
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
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "price" float NOT NULL,
  "image" varchar[] NOT NULL,
  "description" varchar NOT NULL,
  "author" varchar NOT NULL,
  "publisher" varchar NOT NULL,
  "quantity" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "genres" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "subgenres" (
  "id" BIGSERIAL PRIMARY KEY,
  "genres_id" bigserial NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "carts" (
  "id" BIGSERIAL PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "username" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "books_genres" (
  "id" BIGSERIAL PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "genres_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "books_subgenres" (
  "id" BIGSERIAL PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "subgenres_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "reviews" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar NOT NULL,
  "books_id" bigserial NOT NULL,
  "comments" varchar NOT NULL,
  "rating" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "orders" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "transactions" (
  "id" BIGSERIAL PRIMARY KEY,
  "orders_id" bigserial NOT NULL,
  "books_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "wishlists" (
  "id" BIGSERIAL PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "username" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "address" ("username");

CREATE INDEX ON "address" ("city_id");

CREATE INDEX ON "address" ("district_id");

CREATE INDEX ON "address" ("username", "city_id", "district_id");

CREATE INDEX ON "districts" ("city_id");

CREATE INDEX ON "subgenres" ("genres_id");

CREATE INDEX ON "carts" ("books_id");

CREATE INDEX ON "carts" ("username");

CREATE INDEX ON "carts" ("books_id", "username");

CREATE INDEX ON "books_genres" ("books_id");

CREATE INDEX ON "books_genres" ("genres_id");

CREATE INDEX ON "books_genres" ("books_id", "genres_id");

CREATE INDEX ON "books_subgenres" ("books_id");

CREATE INDEX ON "books_subgenres" ("subgenres_id");

CREATE INDEX ON "books_subgenres" ("books_id", "subgenres_id");

CREATE INDEX ON "reviews" ("username");

CREATE INDEX ON "reviews" ("books_id");

CREATE INDEX ON "reviews" ("username", "books_id");

CREATE INDEX ON "orders" ("username");

CREATE INDEX ON "transactions" ("books_id");

CREATE INDEX ON "transactions" ("orders_id");

CREATE INDEX ON "transactions" ("books_id", "orders_id");

CREATE INDEX ON "wishlists" ("books_id");

CREATE INDEX ON "wishlists" ("username");

CREATE INDEX ON "wishlists" ("books_id", "username");

ALTER TABLE "address" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "address" ADD FOREIGN KEY ("city_id") REFERENCES "cities" ("id");

ALTER TABLE "address" ADD FOREIGN KEY ("district_id") REFERENCES "districts" ("id");

ALTER TABLE "districts" ADD FOREIGN KEY ("city_id") REFERENCES "cities" ("id");

ALTER TABLE "subgenres" ADD FOREIGN KEY ("genres_id") REFERENCES "genres" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "books_genres" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "books_genres" ADD FOREIGN KEY ("genres_id") REFERENCES "genres" ("id");

ALTER TABLE "books_subgenres" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "books_subgenres" ADD FOREIGN KEY ("subgenres_id") REFERENCES "subgenres" ("id");

ALTER TABLE "reviews" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "reviews" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "transactions" ADD FOREIGN KEY ("orders_id") REFERENCES "orders" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "wishlists" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "wishlists" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
