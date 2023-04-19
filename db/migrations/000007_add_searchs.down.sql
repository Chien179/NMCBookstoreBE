DROP INDEX IF EXISTS "searchs_idx";

DROP TRIGGER IF EXISTS "searchs_tsv_trigger" ON "searchs";

DROP FUNCTION IF EXISTS "searchs_tsv_trigger_func";

DROP TABLE IF EXISTS "searchs" CASCADE;