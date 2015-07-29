# geonet-rest

Rest API for GeoNet web site data.

## Development 

Requires Go 1.2.1 or newer (for db.SetMaxOpenConns(n)).

### Dependencies and Compilation

Dependencies are included in this repo using godep vendoring.  There should be no need to `go get` the dependencies 
separately unless you are updating them.

* Install godep (you will need Git and Mercurial installed to do this). https://github.com/tools/godep
* Prefix go commands with godep.

Run:

```godep go build && ./geonet-rest```

Run all tests (including those in sub dirs):

```godep go test ./...```

### API Design

* URIs should return a resource and the query parameters should be used to filter (search) for them.
* Use ISO8601 date times in UTC e.g., `2013-05-30T15:15:37.812Z`
* Use http methods in routes (`GET`, `PUT` etc).
* Use camelCase for query and property names.  Be consistent with SeisComPML or QuakeML e.g., `publicID` not `publicId` or `publicid`.
* The  http `Accept-Header` should be used to determine which data version and format to return.

### API Documentation

API documentation is generated from doc{} structs in the code.  Run the application and visit `http://localhost:8080/api-docs`.

### API Changes

#### Non Breaking Changes

* Make non breaking **additions** as required.
* Add to the tests.
* Add Markdown documention to the tests and regenerate the API docs.

#### Breaking Changes

* Are you really sure you have to.  Discuss widely.
* Copy the current API verion code to the next API version (so as to support all queries at the new version)
* Monotonically increment the `Accept` constant e.g., `application/vnd.geo+json;version=1 -> application/vnd.geo+json;version=2`
* Change the tests.  
* Update the documentation.  
* Make the changes.  
* Update the routes.  


### Database

Use the database from the haz-db project.

Use `./scripts/init-db.sh` to initialise impact test data.

### Properties

Either or both of: 
1. Copy an appropriately edited version of `geonet-rest.json` to `/etc/sysconfig/geonet-rest.json`  This should include write access credentials for accessing the impact database.
2. Refer to docker-run.sh for overriding from env var.

## Deployment

### Properties 

Copy an appropriately edited version of `geonet-rest.json` to `/etc/sysconfig/geonet-rest.json`  This should include read only credentials for accessing the hazard database.  Properties can also be set from env var.

### Monitoring

There are state of health pages available for montoring with web probes:

* http://.../soh - this will return a 500 error if any HeartBeat messages in the DB are old.
* http://.../soh/impact - this will return a 500 error if the measured shaking intensity messages fall below 50.  Not all servers may be receiving these messages.

### Logging and Metrics

Can be sent to Logentries and Librato Metrics respectively by setting the appropriate credentials in the config.

Fatal application errors, 4xx and 5xx requests are syslogged.

### Regions

Regions change very rarely and are served with a long surrogate cache time.  If the regions are changed the regions will need to be
purged from CDN.

### Database

Use the haz-db-load tool in the haz-db project to sync a new DB.
