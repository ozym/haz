# GeoNet Quake Search

Provide web service to search the GeoNet quake catalog

## History
This project is migrated from ***quakesearch-go***, as ***quakesearch*** is using now the ***haz database***,
it is decided to move it into the ***haz*** project.
For previous commit history, please refer to [quakesearch-go](https://github.com/GeoNet/quakesearch-go)

## Development
Application is developed in GO

### Database
This application searches data on the GeoNet GeoNet HAZ database

## Web service api

Restful web service for search quakes in difference format( geojson, gml, kml, csv)
query parameters: date, location, depth/magnitude

## Interactive web interface
* use interactive map to define search area
* update coordinates by map extent when "Map Extent" selected as default
* allows building search query as well as showing search results on interactive map
* number of quakes to show on map limited to 2000.
* output format currently (query builder): geojson, gml, kml, csv.
* result as url(s) for intended data, also button to download data from browser.
* the maximum number of quakes for each request is limited to 20,000 (to prevent server crash), beyond that multiple requests are suggested.

## Test

Depends on haz database

```
cd ../database
./scripts//initdb-93.sh
```

Run test
```
export $(cat env.list | grep = | xargs) && go test ./...
```

This package also uses Travis CI and publishes notifications to HipChat.

## Build / Deployment


### Configuration
config parameters can be specified in ```env.list``` which are exported as environment variables


### GO
```
go build
./quakesearch
```

### Docker
```
Use the script docker.sh.  This will create a docker container containing quakesearch and echo the command to run it locally.
```

