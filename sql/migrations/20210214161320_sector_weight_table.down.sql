BEGIN;

DROP TRIGGER IF EXISTS set_sector_weight_timestamp on sector_weights;
DROP FUNCTION IF EXISTS trigger_set_timestamp();
DROP TABLE IF EXISTS sector_weights;

COMMIT;