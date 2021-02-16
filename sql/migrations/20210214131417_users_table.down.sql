BEGIN;

DROP TRIGGER IF EXISTS set_users_timestamp on users;
DROP FUNCTION IF EXISTS trigger_set_timestamp();
DROP TABLE IF EXISTS users;

COMMIT;


