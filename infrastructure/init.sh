#!/bin/bash
set -e
for database in orbis-api-db orbis-api-db-test;
 do
initdb /var/lib/postgresql/project
psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$database'" | grep -q 1 || \
psql -U postgres <<-EOSQL
CREATE DATABASE "$database" WITH owner=postgres;
EOSQL
done
psql -U postgres -tc "CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS plpgsql;"