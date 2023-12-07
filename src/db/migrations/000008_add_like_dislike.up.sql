CREATE TABLE "like"(
  id bigserial PRIMARY KEY,
  username varchar NOT NULL,
  review_id bigserial NOT NULL,
  is_like boolean NOT NULL DEFAULT false
);

CREATE TABLE "dislike"(
  id bigserial PRIMARY KEY,
  username varchar NOT NULL,
  review_id bigserial NOT NULL,
  is_dislike boolean NOT NULL DEFAULT false
);

CREATE INDEX ON "like" ("username", "review_id");

CREATE INDEX ON "dislike" ("username", "review_id");

ALTER TABLE "like" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "like" ADD FOREIGN KEY ("review_id") REFERENCES "reviews" ("id");

ALTER TABLE "dislike" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "dislike" ADD FOREIGN KEY ("review_id") REFERENCES "reviews" ("id");