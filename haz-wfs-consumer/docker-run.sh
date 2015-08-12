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
# WFS_CONSUMER_DATABASE_HOST=localhost
#
# database User password (unencrypted).
# WFS_CONSUMER_DATABASE_PASSWORD=test
#
# usually disable or require.
# WFS_CONSUMER_DATABASE_SSL_MODE=disable
#
# database connection pool.
# WFS_CONSUMER_DATABASE_MAX_OPEN_CONNS=2
#
# database connection pool.
# WFS_CONSUMER_DATABASE_MAX_IDLE_CONNS=1
#
# SQS region e.g., ap-southeast-2.
# WFS_CONSUMER_SQS_AWS_REGION=ap-southeast-2
#
# SQS queue name.
# WFS_CONSUMER_SQS_QUEUE_NAME=XXX
#
# SQS queue user access key.
# WFS_CONSUMER_SQS_ACCESS_KEY=XXX
#
# SQS queue user secret.
# WFS_CONSUMER_SQS_SECRET_KEY=XXX
#
# username for Librato.
# WFS_CONSUMER_LIBRATO_USER=
#
# key for Librato.
# WFS_CONSUMER_LIBRATO_KEY=
#
# source for metrics.  Appended to host if not empty.
# WFS_CONSUMER_LIBRATO_SOURCE=
#
# token for Logentries.
# WFS_CONSUMER_LOGENTRIES_TOKEN=

docker run -e "WFS_CONSUMER_DATABASE_HOST=localhost" -e "WFS_CONSUMER_DATABASE_PASSWORD=test" -e "WFS_CONSUMER_DATABASE_SSL_MODE=disable" -e "WFS_CONSUMER_DATABASE_MAX_OPEN_CONNS=2" -e "WFS_CONSUMER_DATABASE_MAX_IDLE_CONNS=1" -e "WFS_CONSUMER_SQS_AWS_REGION=ap-southeast-2" -e "WFS_CONSUMER_SQS_QUEUE_NAME=XXX" -e "WFS_CONSUMER_SQS_ACCESS_KEY=XXX" -e "WFS_CONSUMER_SQS_SECRET_KEY=XXX" -e "WFS_CONSUMER_LIBRATO_USER=" -e "WFS_CONSUMER_LIBRATO_KEY=" -e "WFS_CONSUMER_LIBRATO_SOURCE=" -e "WFS_CONSUMER_LOGENTRIES_TOKEN=" busybox
