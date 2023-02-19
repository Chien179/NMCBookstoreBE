CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "image" varchar NOT NULL,
  "phone_number" varchar NOT NULL,
  "role" boolean NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "address" (
  "id" BIGSERIAL PRIMARY KEY,
  "address" varchar NOT NULL,
  "users_id" bigserial NOT NULL,
  "district" varchar NOT NULL,
  "city" varchar NOT NULL
);

CREATE TABLE "books" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "price" float NOT NULL,
  "image" varchar NOT NULL,
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
  "users_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "books_carts" (
  "id" BIGSERIAL PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "carts_id" bigserial NOT NULL
);

CREATE TABLE "books_genres" (
  "id" BIGSERIAL PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "genres_id" bigserial NOT NULL
);

CREATE TABLE "books_subgenres" (
  "id" BIGSERIAL PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "subgenres_id" bigserial NOT NULL
);

CREATE TABLE "reviews" (
  "id" BIGSERIAL PRIMARY KEY,
  "users_id" bigserial NOT NULL,
  "books_id" bigserial NOT NULL,
  "comments" varchar NOT NULL,
  "rating" int NOT NULL
);

CREATE TABLE "orders" (
  "id" BIGSERIAL PRIMARY KEY,
  "users_id" bigserial NOT NULL
);

CREATE TABLE "transactions" (
  "id" BIGSERIAL PRIMARY KEY,
  "orders_id" bigserial NOT NULL,
  "books_id" bigserial NOT NULL
);

CREATE TABLE "wishlists" (
  "id" BIGSERIAL PRIMARY KEY,
  "users_id" bigserial NOT NULL
);

CREATE TABLE "books_wishlists" (
  "id" BIGSERIAL PRIMARY KEY,
  "books_id" bigserial NOT NULL,
  "wishlists_id" bigserial NOT NULL
);

CREATE INDEX ON "address" ("users_id");

CREATE INDEX ON "subgenres" ("genres_id");

CREATE INDEX ON "carts" ("users_id");

CREATE INDEX ON "books_carts" ("books_id");

CREATE INDEX ON "books_carts" ("carts_id");

CREATE INDEX ON "books_carts" ("books_id", "carts_id");

CREATE INDEX ON "books_genres" ("books_id");

CREATE INDEX ON "books_genres" ("genres_id");

CREATE INDEX ON "books_genres" ("books_id", "genres_id");

CREATE INDEX ON "books_subgenres" ("books_id");

CREATE INDEX ON "books_subgenres" ("subgenres_id");

CREATE INDEX ON "books_subgenres" ("books_id", "subgenres_id");

CREATE INDEX ON "reviews" ("users_id");

CREATE INDEX ON "reviews" ("books_id");

CREATE INDEX ON "reviews" ("users_id", "books_id");

CREATE INDEX ON "orders" ("users_id");

CREATE INDEX ON "transactions" ("books_id");

CREATE INDEX ON "transactions" ("orders_id");

CREATE INDEX ON "transactions" ("books_id", "orders_id");

CREATE INDEX ON "wishlists" ("users_id");

CREATE INDEX ON "books_wishlists" ("books_id");

CREATE INDEX ON "books_wishlists" ("wishlists_id");

CREATE INDEX ON "books_wishlists" ("books_id", "wishlists_id");

ALTER TABLE "address" ADD FOREIGN KEY ("users_id") REFERENCES "users" ("id");

ALTER TABLE "subgenres" ADD FOREIGN KEY ("genres_id") REFERENCES "genres" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("users_id") REFERENCES "users" ("id");

ALTER TABLE "books_carts" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "books_carts" ADD FOREIGN KEY ("carts_id") REFERENCES "carts" ("id");

ALTER TABLE "books_genres" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "books_genres" ADD FOREIGN KEY ("genres_id") REFERENCES "genres" ("id");

ALTER TABLE "books_subgenres" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "books_subgenres" ADD FOREIGN KEY ("subgenres_id") REFERENCES "subgenres" ("id");

ALTER TABLE "reviews" ADD FOREIGN KEY ("users_id") REFERENCES "users" ("id");

ALTER TABLE "reviews" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("users_id") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("orders_id") REFERENCES "orders" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "wishlists" ADD FOREIGN KEY ("users_id") REFERENCES "users" ("id");

ALTER TABLE "books_wishlists" ADD FOREIGN KEY ("books_id") REFERENCES "books" ("id");

ALTER TABLE "books_wishlists" ADD FOREIGN KEY ("wishlists_id") REFERENCES "wishlists" ("id");
