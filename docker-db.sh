#!/bin/bash

# The DB image
docker build --rm=true -t quay.io/geonet/haz:database -f database/Dockerfile database