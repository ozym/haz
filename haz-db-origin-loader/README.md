# haz-db-origin-loader

* A tool to bulk load quake data from SC3ML into the qrt schema from the `geonet` origin web
server project.


Initialize the database from the `geonet` project.

## Quake Bulk Load

Loads the database with quake information from a directory of SC3ML.

Download some SC3ML using the aws cli (the bucket should be publicly accessible for read) e.g.,

```
aws s3 sync s3://seiscompml07 /work/seismcompml07 --exclude "*"  --include "2015p*"
```

Build and run the db loader.

```
godep go build && ./haz-db-origin-loader
```
