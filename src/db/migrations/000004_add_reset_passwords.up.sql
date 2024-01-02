CREATE TABLE "reset_passwords" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "reset_code" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expired_at" timestamptz NOT NULL DEFAULT (now()  + interval '15 minutes')
);

ALTER TABLE "reset_passwords" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");