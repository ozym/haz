# GeoNet API

Welcome api.geonet.org.nz The data provided here is used for the GeoNet web site and other similar services. 

The GeoNet project makes all its data and images freely available. Please ensure you have read and understood our Data Policy and Disclaimer before using any of these services.

## Versioning

API queries may be versioned via the Accept header. Please specify the `Accept` header for your request exactly as specified for the endpoint query you are using.

If you don't specify an Accept header with a version then your request will be routed to the current highest API version of the query.

Taking advantage of the API versioning will pay dividends in the future for any client that you write. We use the [jq](https://stedolan.github.io/jq/) command for JSON pretty printing etc. A curl command might look like:

    curl -H "Accept: application/vnd.geo+json;version=2" "http://...API-QUERY..." | jq .

You may also be able to find a browser plugin to help with setting the Accept header for requests.

## Compression

The response for a query can be compressed. If your client can handle a compressed response then the reduced download size is a great benifit. Gzip compression is supported. You can request a compressed response by including `gzip` in your `Accept-Encoding` header.
Bugs

## Bugs

The code that provide these services is available at https://github.com/GeoNet/haz/geonet-rest If you believe you have found a bug please raise an issue or pull request there.

# Endpoints

* [Intensity](#intensity)
* [News](#news)
* [Quake](#quake)
* [Quake History](#quakehistory)
* [Quake Stats](#quakestats)
* [Quakes](#quakes)
* [Quake CAP](#quakecap)
* [Quake CAP Feed](#quakecapfeed)

## Intensity ## {#intensity}

Retrieve shaking intensity information.

    [GET] /intensity?type=(measured|reported)[&publicID=(publicID)]

### Accept Version

    application/vnd.geo+json;version=2

### Parameters

type
:   the type of shaking information, either `reported` or `measured`.   If `publicID` is not specified the information is for the last 60 minutes.

publicID
:   *optional* - a valid quake id.  Returns shaking information in a time window around the quake time.  Only available with `type=reported`.

### Response

GeoJSON features with the following properties for both `reported` and `measured` queries.

mmi
:   the MMI measured at the point or reported in the area around the point.

For `measured` queries the addtional GeoJSON properites

count
:   the count of MMI values reported in the area of around the point.		

As well as a summary of report counts by MMI.

### Examples

* [/intensity?type=measured](/intensity?type=measured)
* [/intensity?type=reported](/intensity?type=reported)
* [/intensity?type=reported&publicID=2013p407387](/intensity?type=reported&publicID=2013p407387)

## News ## {#news}

A simple JSON feed of GeoNet news.

    [GET] /news/geonet

### Accept Version

    application/json;version=2

### Parameters

No query parameters are required.

### Response

JSON with the following properties

mlink
:   a link to a mobile version of the news story.

link
:   a link to the news story.

published
:   the date the story was published.

title
:   the title of the story.

### Examples

[/news/geonet](/news/geonet)

## Quake ## {#quake}

Information for a single quake

    [GET] /quake/(publicID)

### Accept Version

    application/vnd.geo+json;version=2

### Parameters

publicID
:   A valid publicID for a quake e.g. `2014p715167`    

### Response

 GeoJSON features with the following properties:

publicID
:   the unique public identifier for this quake.

time
:   the origin time of the quake.

depth
:    the depth of the quake in km.

magnitude
:   the summary magnitude for the quake.

locality
:    distance and direction to the nearest locality.

MMI
:   the calculated MMI shaking at the closest locality in the New Zealand region.

quality
:   the quality of this information; `best`, `good`, `caution`, `deleted`.

### Examples

[/quake/2013p407387](/quake/2013p407387)

## Quake History ## {#quakehistory}

Location history for a single quake.  Not all quakes have a location history.

    [GET] /quake/history/(publicID)

### Accept Version

    application/vnd.geo+json;version=2

### Parameters

publicID
:   A valid publicID for a quake e.g. `2014p715167`    

### Response

 GeoJSON features with the following properties:

publicID
:   the unique public identifier for this quake.

time
:   the origin time of the quake.

modificationTime
:   the modification time of this information.

depth
:    the depth of the quake in km.

magnitude
:   the summary magnitude for the quake.

locality
:    distance and direction to the nearest locality.

MMI
:   the calculated MMI shaking at the closest locality in the New Zealand region.

quality
:   the quality of this information; `best`, `good`, `caution`, `deleted`.

### Examples

[/quake/history/2013p407387](/quake/history/2013p407387)

## Quake Stats ## {#quakestats}

Quake statistics. 

 [GET] /quake/stats

### Accept Version

    application/vnd.geo+json;version=2

### Response

magnitudeCount
:  contains three members that summarise magnitude by count over the last 7, 28, and 365 days.

rate
:  contains the member perDay that gives a per day summary by count of quake occurence.

### Examples

[/quake/stats](/quake/stats)

## Quakes ## {#quakes}

Rerturns quakes possibly felt in the New Zealand region during the last 365 days up to a maximum of 100 quakes.

    [GET] /quake?MMI=(int)

### Accept Version

    application/vnd.geo+json;version=2

### Parameters

MMI
:   request quakes that may have caused shaking greater than or equal to the MMI value in the New Zealand region.  Allowable values are `-1..8` inclusive.  `-1` is used for quakes that are to small to calculate a stable MMI value for.

### Response

 GeoJSON features with the following properties:

publicID
:   the unique public identifier for this quake.

time
:   the origin time of the quake.

depth
:   the depth of the quake in km.

magnitude
:   the summary magnitude for the quake.

locality
:   distance and direction to the nearest locality.

mmi
:   the calculated MMI shaking at the closest locality in the New Zealand region.

quality
:   the quality of this information; `best`, `good`, `caution`, `deleted`.

### Examples

[/quake?MMI=3](/quake?MMI=3)

## Quake CAP ## {#quakecap}

Information in CAP format for a single quake.

    [GET] /cap/1.2/GPA1.0/quake/(ID)

### Accept Version

Queries to this endpoint are not versioned by accept header.

### Parameters

ID
:   a valid quake CAP ID e.g., `2015p653589.1440964609134673`

### Examples

[/cap/1.2/GPA1.0/quake/2015p653589.1440964609134673](/cap/1.2/GPA1.0/quake/2013p407387.1370036261549894)

## Quake CAP Feed ## {#quakecapfeed}

Quake feed with CAP links for alerting.

Feed of quakes in the last seven days of intensity moderate or higher in the New Zealand region and a suitable quality for alerting.   Links (type `application/cap+xml`) to individual quakes in the requested CAP version and profile are included in the returned feed.

    [GET] /cap/1.2/GPA1.0/feed/atom1.0/quake

### Accept Version

queries to this endpoint are not versioned by accept header.

### Examples

[/cap/1.2/GPA1.0/feed/atom1.0/quake](/cap/1.2/GPA1.0/feed/atom1.0/quake)
