DROP TRIGGER IF EXISTS TRIGGER "rating_trigger" ON "reviews";

DROP FUNCTION IF EXISTS rating_trigger_func();

ALTER TABLE "books" DROP COLUMN "rating";