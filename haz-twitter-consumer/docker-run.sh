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
# SQS region e.g., ap-southeast-2.
# TWITTER_CONSUMER_SQS_AWS_REGION=ap-southeast-2
#
# SQS queue name.
# TWITTER_CONSUMER_SQS_QUEUE_NAME=XXX
#
# SQS queue user access key.
# TWITTER_CONSUMER_SQS_ACCESS_KEY=XXX
#
# SQS queue user secret.
# TWITTER_CONSUMER_SQS_SECRET_KEY=XXX
#
# username for Librato.
# TWITTER_CONSUMER_LIBRATO_USER=
#
# key for Librato.
# TWITTER_CONSUMER_LIBRATO_KEY=
#
# source for metrics.  Appended to host if not empty.
# TWITTER_CONSUMER_LIBRATO_SOURCE=
#
# token for Logentries.
# TWITTER_CONSUMER_LOGENTRIES_TOKEN=

docker run -e "TWITTER_CONSUMER_SQS_AWS_REGION=ap-southeast-2" -e "TWITTER_CONSUMER_SQS_QUEUE_NAME=XXX" -e "TWITTER_CONSUMER_SQS_ACCESS_KEY=XXX" -e "TWITTER_CONSUMER_SQS_SECRET_KEY=XXX" -e "TWITTER_CONSUMER_LIBRATO_USER=" -e "TWITTER_CONSUMER_LIBRATO_KEY=" -e "TWITTER_CONSUMER_LIBRATO_SOURCE=" -e "TWITTER_CONSUMER_LOGENTRIES_TOKEN=" busybox
