#!/bin/bash

# script for initializing the db in the postgres Docker container.

export PGUSER=postgres

cd /docker-entrypoint-initdb.d

psql  -d postgres < /docker-entrypoint-initdb.d/create-users.ddl
psql  -d postgres < /docker-entrypoint-initdb.d/create-db.ddl
psql  -d hazard -c 'create extension postgis;'
psql  --quiet hazard < /docker-entrypoint-initdb.d/drop-create.ddl
psql  --quiet hazard < /docker-entrypoint-initdb.d/impact-create.ddl
psql  --quiet hazard < /docker-entrypoint-initdb.d/impact-functions.ddl
psql  --quiet hazard < /docker-entrypoint-initdb.d/user-permissions.ddl
