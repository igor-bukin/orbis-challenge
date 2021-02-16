BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION trigger_set_timestamp() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS sector_weights
(
    id         UUID          NOT NULL DEFAULT uuid_generate_v4(),
    name       text          NOT NULL,
    weight     NUMERIC(5, 2) NOT NULL,

    etfs_id    UUID          NOT NULL references etfs (id),

    created_at TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP     NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

CREATE TRIGGER set_sector_weight_timestamp
    BEFORE UPDATE
    ON sector_weights
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

COMMIT;