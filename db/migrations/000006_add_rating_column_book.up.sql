ALTER TABLE "books" ADD COLUMN "rating" FLOAT NOT NULL DEFAULT 0;

CREATE OR REPLACE FUNCTION rating_trigger_func()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
	WITH rating_avg AS (
        SELECT b.id, COALESCE(AVG(r.rating), 0)::NUMERIC(10,1) AS ravg
		FROM books b
		INNER JOIN reviews r
		ON r.books_id = b.id
		GROUP BY b.id
    )

	UPDATE "books"
	SET rating = rating_avg.ravg
	FROM rating_avg
	WHERE books.id = rating_avg.id;
	
	WITH rating_avg AS (
        SELECT b.id, COALESCE(AVG(r.rating), 0)::NUMERIC(10,1) AS ravg
		FROM books b
		INNER JOIN reviews r
		ON r.books_id = b.id
		GROUP BY b.id
    )
	
	UPDATE "searchs"
	SET rating = rating_avg.ravg
	FROM rating_avg
	WHERE searchs.id = rating_avg.id;
	
	RETURN NEW;
END $$;

CREATE TRIGGER "rating_trigger" AFTER INSERT OR UPDATE
OF rating ON reviews FOR EACH ROW
EXECUTE PROCEDURE rating_trigger_func();