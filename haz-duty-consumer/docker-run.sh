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
# DUTY_CONSUMER_SQS_AWS_REGION=ap-southeast-2
#
# SQS queue name.
# DUTY_CONSUMER_SQS_QUEUE_NAME=
#
# SQS queue user access key.
# DUTY_CONSUMER_SQS_ACCESS_KEY=
#
# SQS queue user secret.
# DUTY_CONSUMER_SQS_SECRET_KEY=
#
# username for Librato.
# DUTY_CONSUMER_LIBRATO_USER=
#
# key for Librato.
# DUTY_CONSUMER_LIBRATO_KEY=
#
# source for metrics.  Appended to host if not empty.
# DUTY_CONSUMER_LIBRATO_SOURCE=
#
# token for Logentries.
# DUTY_CONSUMER_LOGENTRIES_TOKEN=
#
# PagerDuty api token as per https://developer.pagerduty.com/documentation/rest/authentication
# DUTY_CONSUMER_PAGERDUTY_TOKEN=
#
# PagerDuty service GUID as per https://developer.pagerduty.com/documentation/integration/events/trigger
# DUTY_CONSUMER_PAGERDUTY_SERVICE=

docker run -e "DUTY_CONSUMER_SQS_AWS_REGION=ap-southeast-2" -e "DUTY_CONSUMER_SQS_QUEUE_NAME=" -e "DUTY_CONSUMER_SQS_ACCESS_KEY=" -e "DUTY_CONSUMER_SQS_SECRET_KEY=" -e "DUTY_CONSUMER_LIBRATO_USER=" -e "DUTY_CONSUMER_LIBRATO_KEY=" -e "DUTY_CONSUMER_LIBRATO_SOURCE=" -e "DUTY_CONSUMER_LOGENTRIES_TOKEN=" -e "DUTY_CONSUMER_PAGERDUTY_TOKEN=" -e "DUTY_CONSUMER_PAGERDUTY_SERVICE=" busybox
