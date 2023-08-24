DROP INDEX IF EXISTS "searchs_idx";

DROP TRIGGER IF EXISTS "update_searchs_table_trigger" ON "books";

DROP TRIGGER IF EXISTS "update_searchs_table_trigger" ON "genres";

DROP TRIGGER IF EXISTS "update_searchs_table_trigger" ON "subgenres";

DROP TRIGGER IF EXISTS "update_searchs_table_trigger" ON "books_genres";

DROP TRIGGER IF EXISTS "update_searchs_table_trigger" ON "books_subgenres";

DROP FUNCTION IF EXISTS "update_searchs_table_trigger_func";

DROP FUNCTION IF EXISTS "create_searchs_table_func";

DROP EXTENSION IF EXISTS "unaccent";

DROP TABLE IF EXISTS "searchs" CASCADE;