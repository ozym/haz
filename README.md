# haz

Projects related to haz messaging and the haz DB.  

[![Build Status](https://travis-ci.org/GeoNet/haz.svg?branch=master)](https://travis-ci.org/GeoNet/haz)

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

Uses postgis.  There is a script to initialise the DB e.g.,

```
docker run --name postgis -p 5432:5432 -d quay.io/geonet/postgis:latest
./database/scripts/initdb-93.sh
```

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

TODO - this can be improved with Docker compose.

Each subproject that is run in Docker has a  `docker-run.sh` that documents env var config overrides.

Container linking must be done manually e.g., find the DB container IP address and then run the `geonet-rest` container linked to the DB:

```
docker inspect --format='{{.NetworkSettings.IPAddress}}' postgis
> 172.17.0.7

docker run -e "GEONET_REST_DATABASE_HOST=172.17.0.7" -p 8080:8080 quay.io/geonet/haz:geonet-rest
```

Visit the api-docs at http://localhost:8080/api-docs 

Add in messaging containers (use `haz-aws-messaging` to generate AWS resources.  The `nobody` user in the container will need to be able to read and remove files from the spool dir:

```
mkdir /work/spool
sudo chown nobody /work/spool
```

```
docker run -e "SC3_PRODUCER_SNS_ACCESS_KEY=XXX" -e "SC3_PRODUCER_SNS_SECRET_KEY=XXX" -e "SC3_PRODUCER_SNS_TOPIC_ARN=XXX"  -v /work/spool:/work/spool quay.io/geonet/haz:haz-sc3-producer
docker run -e "HAZ_DB_CONSUMER_DATABASE_HOST=XXX"  -e "HAZ_DB_CONSUMER_SQS_QUEUE_NAME=XXX" -e "HAZ_DB_CONSUMER_SQS_ACCESS_KEY=XXX" -e "HAZ_DB_CONSUMER_SQS_SECRET_KEY=XXX" -v /work/spool:/work/spool quay.io/geonet/haz:haz-db-consumer
```

Visit http://localhost:8080/soh and check that heartbeat messages are arriving.

Send a quake and check it arrived (you should get GeoJSON for the quake):

```
rsync /work/seismcompml07/2015p494191.xml /work/spool/
curl http://localhost:8080/quake/2015p494191
```