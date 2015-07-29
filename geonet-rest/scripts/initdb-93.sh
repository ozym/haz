#!/bin/bash

ddl_dir=$(dirname $0)/../ddl

user=postgres
db_user=${1:-$user}
export PGPASSWORD=$2

# A script to initialise the test data in the database.
#
psql --host=127.0.0.1 --quiet --username=$db_user hazard -f ${ddl_dir}/impact-test-data.ddl
