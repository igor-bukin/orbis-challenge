BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION trigger_set_timestamp() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS holdings
(
    id         UUID          NOT NULL DEFAULT uuid_generate_v4(),
    name       text          NOT NULL,
    weight     NUMERIC(5, 2) NOT NULL,

    etfs_id    UUID          NOT NULL references etfs (id),

    created_at TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP     NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

CREATE TRIGGER set_holdings_timestamp
    BEFORE UPDATE
    ON holdings
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

COMMIT;