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
# username for Librato.
# SC3_PRODUCER_LIBRATO_USER=
#
# key for Librato.
# SC3_PRODUCER_LIBRATO_KEY=
#
# source for metrics.  Appended to host if not empty.
# SC3_PRODUCER_LIBRATO_SOURCE=
#
# token for Logentries.
# SC3_PRODUCER_LOGENTRIES_TOKEN=
#
# SNS region e.g., ap-southeast-2.
# SC3_PRODUCER_SNS_AWS_REGION=ap-southeast-2
#
# SNS queue user access key.
# SC3_PRODUCER_SNS_ACCESS_KEY=XXX
#
# SNS queue user secret.
# SC3_PRODUCER_SNS_SECRET_KEY=XXX
#
# SNS Topic Arn.
# SC3_PRODUCER_SNS_TOPIC_ARN=XXX
#
# A service id for heartbeat messages.
# SC3_PRODUCER_SERVICE_ID=haz-sc3-producer.localhost
#
# Spool directory for SeisComPML files.
# SC3_PRODUCER_SC3_SPOOL_DIR=/work/spool

docker run -e "SC3_PRODUCER_LIBRATO_USER=" -e "SC3_PRODUCER_LIBRATO_KEY=" -e "SC3_PRODUCER_LIBRATO_SOURCE=" -e "SC3_PRODUCER_LOGENTRIES_TOKEN=" -e "SC3_PRODUCER_SNS_AWS_REGION=ap-southeast-2" -e "SC3_PRODUCER_SNS_ACCESS_KEY=XXX" -e "SC3_PRODUCER_SNS_SECRET_KEY=XXX" -e "SC3_PRODUCER_SNS_TOPIC_ARN=XXX" -e "SC3_PRODUCER_SERVICE_ID=haz-sc3-producer.localhost" -e "SC3_PRODUCER_SC3_SPOOL_DIR=/work/spool" busybox
