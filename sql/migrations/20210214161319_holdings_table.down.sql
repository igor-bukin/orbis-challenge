BEGIN;

DROP TRIGGER IF EXISTS set_holdings_timestamp on holdings;
DROP FUNCTION IF EXISTS trigger_set_timestamp();
DROP TABLE IF EXISTS holdings;

COMMIT;