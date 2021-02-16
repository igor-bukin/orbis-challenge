BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION trigger_set_timestamp() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS etfs
(
    id         UUID      NOT NULL DEFAULT uuid_generate_v4(),
    name       text      NOT NULL UNIQUE,
    ticker     text      NOT NULL,
    fund_uri   text      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

CREATE TRIGGER set_efts_timestamp
    BEFORE UPDATE
    ON etfs
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

COMMIT;