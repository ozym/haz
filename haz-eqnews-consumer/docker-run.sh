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
# EQNEWS_SQS_AWS_REGION=ap-southeast-2
#
# SQS queue name.
# EQNEWS_SQS_QUEUE_NAME=
#
# SQS queue user access key.
# EQNEWS_SQS_ACCESS_KEY=
#
# SQS queue user secret.
# EQNEWS_SQS_SECRET_KEY=
#
# username for Librato.
# EQNEWS_LIBRATO_USER=
#
# key for Librato.
# EQNEWS_LIBRATO_KEY=
#
# source for metrics.  Appended to host if not empty.
# EQNEWS_LIBRATO_SOURCE=
#
# token for Logentries.
# EQNEWS_LOGENTRIES_TOKEN=
#
# The SMTP host.
# EQNEWS_SMTP_HOST=email-smtp.us-west-2.amazonaws.com
#
# The SMTP port.
# EQNEWS_SMTP_PORT=587
#
# The SMTP user.
# EQNEWS_SMTP_USER=
#
# The SMTP password.
# EQNEWS_SMTP_PASSWORD=
#
# The from email address
# EQNEWS_SMTP_FROM=
#
# The to email address
# EQNEWS_SMTP_TO=

docker run -e "EQNEWS_SQS_AWS_REGION=ap-southeast-2" -e "EQNEWS_SQS_QUEUE_NAME=" -e "EQNEWS_SQS_ACCESS_KEY=" -e "EQNEWS_SQS_SECRET_KEY=" -e "EQNEWS_LIBRATO_USER=" -e "EQNEWS_LIBRATO_KEY=" -e "EQNEWS_LIBRATO_SOURCE=" -e "EQNEWS_LOGENTRIES_TOKEN=" -e "EQNEWS_SMTP_HOST=email-smtp.us-west-2.amazonaws.com" -e "EQNEWS_SMTP_PORT=587" -e "EQNEWS_SMTP_USER=" -e "EQNEWS_SMTP_PASSWORD=" -e "EQNEWS_SMTP_FROM=" -e "EQNEWS_SMTP_TO=" busybox
