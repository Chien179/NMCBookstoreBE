CREATE TABLE "searchs" AS
SELECT 
	b.id AS id,
	b."name" AS book_names,
	b.price AS price,
	b.author AS author,
	b.publisher AS publisher,
	b.rating AS rating,
	g."name" AS genres,
	s."name" AS subgenres
FROM books b
INNER JOIN books_genres bg
ON b.id = bg.id
INNER JOIN books_subgenres bs
ON b.id = bs.books_id
INNER JOIN genres g
ON bg.genres_id = g.id
INNER JOIN subgenres s
ON bs.subgenres_id = s.id;

ALTER TABLE "searchs" ADD COLUMN "searchs_tsv" tsvector NOT NULL;

UPDATE "searchs"
SET searchs_tsv = to_tsvector(book_names || ' ' || author || ' ' || publisher || ' ' || genres || ' ' || subgenres);

CREATE OR REPLACE FUNCTION searchs_tsv_trigger_func()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN NEW.searchs_tsv =
	setweight(to_tsvector(coalesce(unaccent(NEW.book_names))), 'A') ||
	setweight(to_tsvector(coalesce(unaccent(NEW.genres))), 'A') || 
	setweight(to_tsvector(coalesce(unaccent(NEW.subgenres))), 'A') ||
	setweight(to_tsvector(coalesce(unaccent(NEW.author))), 'B') || 
	setweight(to_tsvector(coalesce(unaccent(NEW.publisher))), 'B');
	RETURN NEW;
END $$;

CREATE TRIGGER searchs_tsv_trigger BEFORE INSERT OR UPDATE
OF book_names, author, publisher, genres, subgenres ON searchs FOR EACH ROW
EXECUTE PROCEDURE searchs_tsv_trigger_func();

CREATE INDEX searchs_idx ON searchs USING GIN(searchs_tsv);