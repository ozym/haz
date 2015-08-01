#!/bin/bash

mkdir -p docker-build-tmp
chmod +s docker-build-tmp

BUILD='-X github.com/GeoNet/cfg.Build -'`git rev-parse --short HEAD` 

# Build all executables in the golang-godep container.  Output statically linked binaries to docker-build-tmp
docker run -e "GOBIN=/usr/src/go/src/github.com/GeoNet/haz/docker-build-tmp"  -e "CGO_ENABLED=0" -e "GOOS=linux" -e "BUILD=$BUILD" --rm -v \
"$PWD":/usr/src/go/src/github.com/GeoNet/haz -w /usr/src/go/src/github.com/GeoNet/haz quay.io/geonet/golang-godep godep go install -a  -ldflags "${BUILD}" -installsuffix cgo ./...

# Assemble common resource for ssl, timezones, and user.
mkdir -p docker-build-tmp/common/etc/ssl/certs
mkdir -p docker-build-tmp/common/usr/share
echo "nobody:x:65534:65534:Nobody:/:" > docker-build-tmp/common/etc/passwd
cp /etc/ssl/certs/ca-certificates.crt docker-build-tmp/common/etc/ssl/certs
# An alternative is to use $GOROOT/lib/time/zoneinfo.zip
rsync --archive /usr/share/zoneinfo docker-build-tmp/common/usr/share

# Docker images for apps
for i in *-consumer haz-sc3-producer
do
	cp "${i}/${i}.json" docker-build-tmp
	echo "FROM scratch" > docker-build-tmp/Dockerfile
	echo "ADD common ${i} ${i}.json /" >> docker-build-tmp/Dockerfile
	echo "USER nobody" >> docker-build-tmp/Dockerfile
	echo "CMD [\"/${i}\"]" >> docker-build-tmp/Dockerfile
	docker build --rm=true -t quay.io/geonet/haz:$i -f docker-build-tmp/Dockerfile docker-build-tmp
done

# Docker images for web apps with an open port and a tmpl directory
for i in geonet-rest 
do
	cp "${i}/${i}.json" docker-build-tmp
	rm -rf docker-build-tmp/common/tmpl
	rsync --archive "${i}/tmpl" docker-build-tmp/common
	echo "FROM scratch" > docker-build-tmp/Dockerfile
	echo "ADD common ${i} ${i}.json /" >> docker-build-tmp/Dockerfile
	echo "USER nobody" >> docker-build-tmp/Dockerfile
	echo "EXPOSE 8080" >> docker-build-tmp/Dockerfile
	echo "CMD [\"/${i}\"]" >> docker-build-tmp/Dockerfile
	docker build --rm=true -t quay.io/geonet/haz:$i -f docker-build-tmp/Dockerfile docker-build-tmp
done
