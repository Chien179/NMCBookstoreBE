CREATE EXTENSION unaccent;

CREATE OR REPLACE FUNCTION create_searchs_table_func()
RETURNS VOID
LANGUAGE plpgsql AS $$
BEGIN
	DROP TABLE IF EXISTS "searchs" CASCADE;
	EXECUTE format('
		CREATE TABLE "searchs" AS
		SELECT
			b.id AS id,
			b."name" AS name,
			b.price AS price,
			b.image AS image,
			b.description AS description,
			b.author AS author,
			b.publisher AS publisher,
			b.quantity AS quantity,
			b.rating AS rating,
			b.created_at AS created_at,
			g."name" AS genres,
			s."name" AS subgenres
		FROM
			books b
			INNER JOIN books_genres bg ON b.id = bg.id
			INNER JOIN books_subgenres bs ON b.id = bs.books_id
			INNER JOIN genres g ON bg.genres_id = g.id
			INNER JOIN subgenres s ON bs.subgenres_id = s.id
	');
	
	ALTER TABLE searchs ADD COLUMN searchs_tsv tsvector;
END $$;

CREATE TABLE "searchs" AS
		SELECT
			b.id AS id,
			b."name" AS name,
			b.price AS price,
			b.image AS image,
			b.description AS description,
			b.author AS author,
			b.publisher AS publisher,
			b.quantity AS quantity,
			b.rating AS rating,
			b.created_at AS created_at,
			g."name" AS genres,
			s."name" AS subgenres
		FROM
			books b
			INNER JOIN books_genres bg ON b.id = bg.id
			INNER JOIN books_subgenres bs ON b.id = bs.books_id
			INNER JOIN genres g ON bg.genres_id = g.id
			INNER JOIN subgenres s ON bs.subgenres_id = s.id;

ALTER TABLE searchs ADD COLUMN searchs_tsv tsvector;

CREATE OR REPLACE FUNCTION update_searchs_table_trigger_func()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
	PERFORM create_searchs_table_func();
	
	UPDATE searchs
	SET
		searchs_tsv =
			setweight(to_tsvector(coalesce(unaccent(book_names))), 'A') ||
			setweight(to_tsvector(coalesce(unaccent(genres))), 'A') || 
			setweight(to_tsvector(coalesce(unaccent(subgenres))), 'A') ||
			setweight(to_tsvector(coalesce(unaccent(author))), 'B') || 
			setweight(to_tsvector(coalesce(unaccent(publisher))), 'B');
			
	CREATE INDEX "searchs_idx" ON "searchs" USING GIN(searchs_tsv);
	
	RETURN NEW;
END $$;

CREATE TRIGGER update_searchs_table_trigger AFTER INSERT OR UPDATE
OF name, price, author, publisher, rating ON books FOR EACH ROW
EXECUTE PROCEDURE update_searchs_table_trigger_func();

CREATE TRIGGER update_searchs_table_trigger AFTER INSERT OR UPDATE
OF name ON genres FOR EACH ROW
EXECUTE PROCEDURE update_searchs_table_trigger_func();

CREATE TRIGGER update_searchs_table_trigger AFTER INSERT OR UPDATE
OF genres_id,name ON subgenres FOR EACH ROW
EXECUTE PROCEDURE update_searchs_table_trigger_func();

CREATE TRIGGER update_searchs_table_trigger AFTER INSERT OR UPDATE
OF books_id,genres_id ON books_genres FOR EACH ROW
EXECUTE PROCEDURE update_searchs_table_trigger_func();

CREATE TRIGGER update_searchs_table_trigger AFTER INSERT OR UPDATE
OF books_id,subgenres_id ON books_subgenres FOR EACH ROW
EXECUTE PROCEDURE update_searchs_table_trigger_func();