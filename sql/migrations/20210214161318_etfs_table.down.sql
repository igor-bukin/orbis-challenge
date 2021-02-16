BEGIN;

DROP TRIGGER IF EXISTS set_efts_timestamp on efts;
DROP FUNCTION IF EXISTS trigger_set_timestamp();
DROP TABLE IF EXISTS etfs;

COMMIT;