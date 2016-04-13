#!/bin/bash -e

# code will be compiled in this container
BUILD_CONTAINER=golang:1.6.1-alpine

DOCKER_TMP=docker-build-tmp

mkdir -p $DOCKER_TMP
chmod +s $DOCKER_TMP
mkdir -p ${DOCKER_TMP}/common/etc/ssl/certs
mkdir -p ${DOCKER_TMP}/common/usr/share

# Prefix for the logs
BUILD='-X github.com/GeoNet/haz/vendor/github.com/GeoNet/log/logentries.Prefix=git-'`git rev-parse --short HEAD`

# Build all executables in the Golang container.  Output statically linked binaries to docker-build-tmp
# Assemble common resource for ssl and timezones
docker run -e "GOBIN=/usr/src/go/src/github.com/GeoNet/haz/${DOCKER_TMP}" -e "GOPATH=/usr/src/go" -e "CGO_ENABLED=0" -e "GOOS=linux" -e "BUILD=$BUILD" --rm \
	-v "$PWD":/usr/src/go/src/github.com/GeoNet/haz \
	-w /usr/src/go/src/github.com/GeoNet/haz ${BUILD_CONTAINER} \
	go install -a  -ldflags "${BUILD}" -installsuffix cgo ./...; \
	cp /etc/ssl/certs/ca-certificates.crt ${DOCKER_TMP}/common/etc/ssl/certs; \
	cp -Ra /usr/share/zoneinfo ${DOCKER_TMP}/common/usr/share

# Assemble common resource for user.
echo "nobody:x:65534:65534:Nobody:/:" > ${DOCKER_TMP}/common/etc/passwd

# Docker images for apps
for i in *-consumer haz-sc3-producer
do
	echo "FROM scratch" > ${DOCKER_TMP}/Dockerfile
	echo "ADD common ${i} /" >> ${DOCKER_TMP}/Dockerfile
	echo "USER nobody" >> ${DOCKER_TMP}/Dockerfile
	echo "CMD [\"/${i}\"]" >> ${DOCKER_TMP}/Dockerfile
	docker build --rm=true -t quay.io/geonet/haz:$i -f ${DOCKER_TMP}/Dockerfile ${DOCKER_TMP}
done

# Docker images for web apps with an open port and a tmpl directory
for i in geonet-rest 
do
	rm -rf ${DOCKER_TMP}/common/tmpl
	rsync --archive "${i}/tmpl" ${DOCKER_TMP}/common
	rm -rf ${DOCKER_TMP}/common/docs
	rsync --archive "${i}/docs" ${DOCKER_TMP}/common
	echo "FROM scratch" > ${DOCKER_TMP}/Dockerfile
	echo "ADD common ${i} /" >> ${DOCKER_TMP}/Dockerfile
	echo "USER nobody" >> ${DOCKER_TMP}/Dockerfile
	echo "EXPOSE 8080" >> ${DOCKER_TMP}/Dockerfile
	echo "CMD [\"/${i}\"]" >> ${DOCKER_TMP}/Dockerfile
	docker build --rm=true -t quay.io/geonet/haz:$i -f ${DOCKER_TMP}/Dockerfile docker-build-tmp
done
