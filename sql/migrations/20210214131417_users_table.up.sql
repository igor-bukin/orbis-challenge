BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION trigger_set_timestamp() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS users
(
    id         UUID      NOT NULL DEFAULT uuid_generate_v4(),
    email      text      NOT NULL UNIQUE,
    password   text,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

CREATE TRIGGER set_users_timestamp
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

INSERT INTO users (email, password)
VALUES ('igor.bukin@test.com', '$2y$12$Z052.WxNKuWJ1x2b9bg9LOSJ07DTHpsYWVHm73BzZJeDcFezflYyW');

COMMIT;