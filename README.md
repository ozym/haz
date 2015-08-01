# haz

Projects related to haz messaging and the haz DB.  

[![Build Status](https://travis-ci.org/GeoNet/haz.svg?branch=master)](https://travis-ci.org/GeoNet/haz)

## Working with Forks

To work with a fork of this repo you will need to clone the fork in such a way to preserve the import paths.  Fork and clone to preserve the organization name e.g.,  

```
 git clone git@github.com:YOUR_FORK/haz.git ${GOPATH}/src/github.com/GeoNet/haz
```

This could possibly include reseting GOPATH if you need to keep origin and a fork

## Sub Projects

### Messaging Applications

#### Producers

* haz-sc3-producer - produces haz messges from SeisComPML.  Need sfile system access in the container.  See Container Testing below.

#### Consumers

Subprojects `*-consumer` consume `msg.Haz` or `msg.Impact` messages from SQS and process them.

### Web Applications

* geonet-rest - the server for api.geonet.org.nz

### Support Applications

* haz-aws-messaging - creates AWS resources for the haz messaging.  `impact-intensity-consumer` has a CFN template for it's resources. 
* haz-db-loader - used to load SeisComPML into the db.  See below.

## Database

Uses postgis.  Pull and run an image that initializes the database on startup (without further configuration):

```
docker run --name hazdb -p 5432:5432 -d  quay.io/geonet/haz:database
``` 

On start the hazard DB is initialized which will delay DB availability.  The image can also be built using `./docker-db.sh`

There is also a script to (re)initialise the DB  `./database/scripts/initdb-93.sh`

### Loading Quake Data

Quake data can be back loaded from SeisComPML.  Download SeisComPML from the S3 bucket and then load it to the DB using `haz-db-loader`:

```
aws s3 sync s3://seiscompml07 /work/seismcompml07 --exclude "*"  --include "2015p*"
cd haz-db-loader
godep go run haz-db-loader.go
```

## Tests

As well as a running database a small amount of impact test data must be added:

```
./geonet-rest/scripts/initdb-93.sh
```

Then run all tests

```
godep go test ./...
```

## Docker Builds

Build static Go binaries that run as the `nobody` user in minimal Docker containers using `docker.sh`.  There is no need to add Dockerfiles to each subproject.  All Docker images are tags in the `haz` repo.

Remove the `docker-build-tmp` dir first to ensure a clean build:

```
rm -rf docker-build-tmp
./docker.sh
```

Once an image has been tested it can be pushed to the repo e.g.,

```
docker push quay.io/geonet/haz:geonet-rest
```

### Container Testing

docker-compose can be used to run a number of the containers together to test messaging.

Copy `dev.env` to `secret.env` (ignored by Git) and edit `secret.env` adding the outputs from running `haz-aws-messaging`.

```
mkdir /work/spool
sudo chown nobody /work/spool
```

Run the containers with `docker-compose up`.

Visit http://localhost:8080/soh and check that heartbeat messages are arriving.

Send a quake and check it arrived (you should get GeoJSON for the quake):

```
rsync /work/seismcompml07/2015p494191.xml /work/spool/
curl http://localhost:8080/quake/2015p494191
```