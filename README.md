# haz

Projects related to haz messaging and the haz DB.  

## Working with Forks

To work with a fork of this repo you will need to clone the fork in such a way to preserve the import paths.  Fork and clone to preserve the organization name e.g.,  

```
 git clone git@github.com:YOUR_FORK/haz.git ${GOPATH}/src/github.com/GeoNet/haz
```

This could possibly include reseting GOPATH if you need to keep origin and a fork

## Sub Projects

### Messaging Applications

#### Producers

* haz-sc3-producer - produces haz messges from SeisComPML.  Need file system access in the container.  See Container Testing below.

#### Consumers

Subprojects `*-consumer` consume `msg.Haz` or `msg.Impact` messages from SQS and process them.

### Web Applications

* geonet-rest - the server for api.geonet.org.nz
* sc3ml-to-quakeml - web services to return QuakeML from SeisComPML.

### Support Applications

* haz-aws-messaging - creates AWS resources for the haz messaging.  `impact-intensity-consumer` has a CFN template for it's resources. 
* haz-db-loader - used to load SeisComPML into the db.  See below.

## Protobufs

Compile protobufs 

## Database

Uses postgis.  Pull and run the image (which already has the hazard db initialised and ready to use):

```
docker run --name hazdb -p 5432:5432 -d 862640294325.dkr.ecr.ap-southeast-2.amazonaws.com/haz-db:9.5
``` 

A Postgres 9.4 with Postgis 2.2 image can be built and pushed using:

```
docker build --rm=true -t 862640294325.dkr.ecr.ap-southeast-2.amazonaws.com/haz-db:9.5 -f database/Dockerfile database
docker push 862640294325.dkr.ecr.ap-southeast-2.amazonaws.com/haz-db:9.5
```

There is also a script to (re)initialise the DB  `./database/scripts/initdb-93.sh`

### Loading Quake Data

Quake data can be back loaded from SeisComPML.  Download SeisComPML from the S3 bucket and then load it to the DB using `haz-db-loader`:

```
aws s3 sync s3://seiscompml07 /work/seismcompml07 --exclude "*"  --include "2015p*"
cd haz-db-loader
go run haz-db-loader.go
```

### Loading Quake Data - Origin Web Servers

Init the database with the qrt schema from the `geonet` project.

Quake data can be back loaded from SeisComPML.  Download SeisComPML from the S3 bucket and then load it to the DB using `haz-db-loader`:

```
aws s3 sync s3://seiscompml07 /work/seismcompml07 --exclude "*"  --include "2015p*"
cd haz-db-origin-loader
go run haz-db-origin-loader.go
```

## Tests

With the DB up run all tests

```
./all.sh
```

## Docker Builds

`./build.sh proj-name [proj-name]...`

`./build-push proj-name [proj-name]...`

Also refer to `.travis.yml`.


### Container Testing

docker-compose can be used to run a number of the containers together to test messaging.

Edit `sc3-producer.env`, `db-consumer.env`, and `geonet-rest.env` adding the outputs from running `haz-aws-messaging`.

```
mkdir /work/spool
sudo chown nobody /work/spool
```

Run the containers with `docker-compose up`.

Visit http://localhost:8080/soh/esb and check that heartbeat messages are arriving.

Send a quake and check it arrived (you should get GeoJSON for the quake):

```
rsync /work/seismcompml07/2015p494191.xml /work/spool/
curl http://localhost:8080/quake/2015p494191
```

### Deployment

#### AWS Elastic Beanstalk

`geonet-rest`, `quakesearch`, `sc3ml-to-quakeml`, `wfs`.

There are files for EB - both to deploy the application and also set
up logging from the container (application) to CloudWatch Logs.  Create a zip file and then upload the 
zip to EB.

```
cd deploy
zip wfs.zip Dockerrun.aws.json .ebextensions/*
```
