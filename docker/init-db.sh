#!/bin/sh -e

POSTGRES_HOSTNAME=${POSTGRES_HOSTNAME:-localhost}
psql -v ON_ERROR_STOP=1 -h "$POSTGRES_HOSTNAME" --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    SELECT 'CREATE DATABASE "kotak"' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname='kotak')\gexec
EOSQL