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
# QUAKE_CONSUMER_DATABASE_HOST=localhost
#
# database User password (unencrypted).
# QUAKE_CONSUMER_DATABASE_PASSWORD=test
#
# usually disable or require.
# QUAKE_CONSUMER_DATABASE_SSL_MODE=disable
#
# database connection pool.
# QUAKE_CONSUMER_DATABASE_MAX_OPEN_CONNS=2
#
# database connection pool.
# QUAKE_CONSUMER_DATABASE_MAX_IDLE_CONNS=1
#
# SQS region e.g., ap-southeast-2.
# QUAKE_CONSUMER_SQS_AWS_REGION=ap-southeast-2
#
# SQS queue name.
# QUAKE_CONSUMER_SQS_QUEUE_NAME=XXX
#
# SQS queue user access key.
# QUAKE_CONSUMER_SQS_ACCESS_KEY=XXX
#
# SQS queue user secret.
# QUAKE_CONSUMER_SQS_SECRET_KEY=XXX
#
# username for Librato.
# QUAKE_CONSUMER_LIBRATO_USER=
#
# key for Librato.
# QUAKE_CONSUMER_LIBRATO_KEY=
#
# source for metrics.  Appended to host if not empty.
# QUAKE_CONSUMER_LIBRATO_SOURCE=
#
# token for Logentries.
# QUAKE_CONSUMER_LOGENTRIES_TOKEN=

docker run -e "QUAKE_CONSUMER_DATABASE_HOST=localhost" -e "QUAKE_CONSUMER_DATABASE_PASSWORD=test" -e "QUAKE_CONSUMER_DATABASE_SSL_MODE=disable" -e "QUAKE_CONSUMER_DATABASE_MAX_OPEN_CONNS=2" -e "QUAKE_CONSUMER_DATABASE_MAX_IDLE_CONNS=1" -e "QUAKE_CONSUMER_SQS_AWS_REGION=ap-southeast-2" -e "QUAKE_CONSUMER_SQS_QUEUE_NAME=XXX" -e "QUAKE_CONSUMER_SQS_ACCESS_KEY=XXX" -e "QUAKE_CONSUMER_SQS_SECRET_KEY=XXX" -e "QUAKE_CONSUMER_LIBRATO_USER=" -e "QUAKE_CONSUMER_LIBRATO_KEY=" -e "QUAKE_CONSUMER_LIBRATO_SOURCE=" -e "QUAKE_CONSUMER_LOGENTRIES_TOKEN=" busybox
