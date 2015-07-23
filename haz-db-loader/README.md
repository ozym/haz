# haz-db 

* Database for haz and impact schema.
* A tool to bulk load quake data from SC3ML.

Needs a recent postgres+postgis DB server (Docker is fine).

Initialize the database with `scripts/initdb-93.sh`.  *Caution* this drops the `hazard` db.

## Quake Bulk Load

Loads the database with quake information from a directory of SC3ML.

If you choose not to use `/work/seismcompml07` as your work dir then change the path in `haz-db.json`.

Download some SC3ML using the aws cli (the bucket should be publicly accessible for read) e.g.,

```
aws s3 sync s3://seiscompml07 /work/seismcompml07 --exclude "*"  --include "2015p*"
```

Build and run the db loader.

```
godep go build && ./haz-db
```
