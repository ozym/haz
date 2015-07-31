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
# INTENSITY_DATABASE_HOST=localhost
#
# database User password (unencrypted).
# INTENSITY_DATABASE_PASSWORD=test
#
# usually disable or require.
# INTENSITY_DATABASE_SSL_MODE=disable
#
# database connection pool.
# INTENSITY_DATABASE_MAX_OPEN_CONNS=2
#
# database connection pool.
# INTENSITY_DATABASE_MAX_IDLE_CONNS=1
#
# SQS region e.g., ap-southeast-2.
# INTENSITY_SQS_AWS_REGION=ap-southeast-2
#
# SQS queue name.
# INTENSITY_SQS_QUEUE_NAME=XXX
#
# SQS queue user access key.
# INTENSITY_SQS_ACCESS_KEY=XXX
#
# SQS queue user secret.
# INTENSITY_SQS_SECRET_KEY=XXX
#
# username for Librato.
# INTENSITY_LIBRATO_USER=
#
# key for Librato.
# INTENSITY_LIBRATO_KEY=
#
# source for metrics.  Appended to host if not empty.
# INTENSITY_LIBRATO_SOURCE=
#
# token for Logentries.
# INTENSITY_LOGENTRIES_TOKEN=

docker run -e "INTENSITY_DATABASE_HOST=localhost" -e "INTENSITY_DATABASE_PASSWORD=test" -e "INTENSITY_DATABASE_SSL_MODE=disable" -e "INTENSITY_DATABASE_MAX_OPEN_CONNS=2" -e "INTENSITY_DATABASE_MAX_IDLE_CONNS=1" -e "INTENSITY_SQS_AWS_REGION=ap-southeast-2" -e "INTENSITY_SQS_QUEUE_NAME=XXX" -e "INTENSITY_SQS_ACCESS_KEY=XXX" -e "INTENSITY_SQS_SECRET_KEY=XXX" -e "INTENSITY_LIBRATO_USER=" -e "INTENSITY_LIBRATO_KEY=" -e "INTENSITY_LIBRATO_SOURCE=" -e "INTENSITY_LOGENTRIES_TOKEN=" busybox
