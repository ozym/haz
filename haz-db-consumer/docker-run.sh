#!/bin/bash

#
# This file is auto generated.  Do not edit.
#
# It was created from the JSON config file and shows the env var that can be used to config the app.
# The docker run command will set the env vars on the container.
# You will need to adjust the image name in the Docker command.
#
# The values shown for the env var are the app defaults from the JSON file.
#
# database host name.
# HAZ_DB_CONSUMER_DATABASE_HOST=localhost
#
# database User password (unencrypted).
# HAZ_DB_CONSUMER_DATABASE_PASSWORD=test
#
# usually disable or require.
# HAZ_DB_CONSUMER_DATABASE_SSL_MODE=disable
#
# database connection pool.
# HAZ_DB_CONSUMER_DATABASE_MAX_OPEN_CONNS=2
#
# database connection pool.
# HAZ_DB_CONSUMER_DATABASE_MAX_IDLE_CONNS=1
#
# SQS region e.g., ap-southeast-2.
# HAZ_DB_CONSUMER_SQS_AWS_REGION=ap-southeast-2
#
# SQS queue name.
# HAZ_DB_CONSUMER_SQS_QUEUE_NAME=XXX
#
# SQS queue user access key.
# HAZ_DB_CONSUMER_SQS_ACCESS_KEY=XXX
#
# SQS queue user secret.
# HAZ_DB_CONSUMER_SQS_SECRET_KEY=XXX
#
# username for Librato.
# HAZ_DB_CONSUMER_LIBRATO_USER=
#
# key for Librato.
# HAZ_DB_CONSUMER_LIBRATO_KEY=
#
# source for metrics.  Appended to host if not empty.
# HAZ_DB_CONSUMER_LIBRATO_SOURCE=
#
# token for Logentries.
# HAZ_DB_CONSUMER_LOGENTRIES_TOKEN=

docker run -e "HAZ_DB_CONSUMER_DATABASE_HOST=localhost" -e "HAZ_DB_CONSUMER_DATABASE_PASSWORD=test" -e "HAZ_DB_CONSUMER_DATABASE_SSL_MODE=disable" -e "HAZ_DB_CONSUMER_DATABASE_MAX_OPEN_CONNS=2" -e "HAZ_DB_CONSUMER_DATABASE_MAX_IDLE_CONNS=1" -e "HAZ_DB_CONSUMER_SQS_AWS_REGION=ap-southeast-2" -e "HAZ_DB_CONSUMER_SQS_QUEUE_NAME=XXX" -e "HAZ_DB_CONSUMER_SQS_ACCESS_KEY=XXX" -e "HAZ_DB_CONSUMER_SQS_SECRET_KEY=XXX" -e "HAZ_DB_CONSUMER_LIBRATO_USER=" -e "HAZ_DB_CONSUMER_LIBRATO_KEY=" -e "HAZ_DB_CONSUMER_LIBRATO_SOURCE=" -e "HAZ_DB_CONSUMER_LOGENTRIES_TOKEN=" busybox
