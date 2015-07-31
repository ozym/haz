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
# UA_CONSUMER_SQS_AWS_REGION=ap-southeast-2
#
# SQS queue name.
# UA_CONSUMER_SQS_QUEUE_NAME=XXX
#
# SQS queue user access key.
# UA_CONSUMER_SQS_ACCESS_KEY=XXX
#
# SQS queue user secret.
# UA_CONSUMER_SQS_SECRET_KEY=XXX
#
# username for Librato.
# UA_CONSUMER_LIBRATO_USER=
#
# key for Librato.
# UA_CONSUMER_LIBRATO_KEY=
#
# source for metrics.  Appended to host if not empty.
# UA_CONSUMER_LIBRATO_SOURCE=
#
# token for Logentries.
# UA_CONSUMER_LOGENTRIES_TOKEN=

docker run -e "UA_CONSUMER_SQS_AWS_REGION=ap-southeast-2" -e "UA_CONSUMER_SQS_QUEUE_NAME=XXX" -e "UA_CONSUMER_SQS_ACCESS_KEY=XXX" -e "UA_CONSUMER_SQS_SECRET_KEY=XXX" -e "UA_CONSUMER_LIBRATO_USER=" -e "UA_CONSUMER_LIBRATO_KEY=" -e "UA_CONSUMER_LIBRATO_SOURCE=" -e "UA_CONSUMER_LOGENTRIES_TOKEN=" busybox
