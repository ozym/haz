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
# SEISCOMPML_S3_LIBRATO_USER=
#
# key for Librato.
# SEISCOMPML_S3_LIBRATO_KEY=
#
# source for metrics.  Appended to host if not empty.
# SEISCOMPML_S3_LIBRATO_SOURCE=
#
# token for Logentries.
# SEISCOMPML_S3_LOGENTRIES_TOKEN=
#
# SNS region e.g., ap-southeast-2.
# SEISCOMPML_S3_SNS_AWS_REGION=ap-southeast-2
#
# S3 user access key.
# SEISCOMPML_S3_S3_ACCESS_KEY=
#
# S3 user secret.
# SEISCOMPML_S3_S3_SECRET_KEY=
#
# Check folder interval
# SEISCOMPML_S3_SEIS_INTERVAL=60
#
# Input directory
# SEISCOMPML_S3_SEIS_IN_DIR=./s3/in
#
# Output directory
# SEISCOMPML_S3_SEIS_OUT_DIR=./s3/out
#
# Unprocess directory
# SEISCOMPML_S3_SEIS_UNPROCESS_DIR=./s3/unprocessed

docker run -e "SEISCOMPML_S3_LIBRATO_USER=" -e "SEISCOMPML_S3_LIBRATO_KEY=" -e "SEISCOMPML_S3_LIBRATO_SOURCE=" -e "SEISCOMPML_S3_LOGENTRIES_TOKEN=" -e "SEISCOMPML_S3_SNS_AWS_REGION=ap-southeast-2" -e "SEISCOMPML_S3_S3_ACCESS_KEY=" -e "SEISCOMPML_S3_S3_SECRET_KEY=" -e "SEISCOMPML_S3_SEIS_INTERVAL=60" -e "SEISCOMPML_S3_SEIS_IN_DIR=./s3/in" -e "SEISCOMPML_S3_SEIS_OUT_DIR=./s3/out" -e "SEISCOMPML_S3_SEIS_UNPROCESS_DIR=./s3/unprocessed" busybox
