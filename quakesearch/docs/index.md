# GeoNet Quakesearch API

Welcome to the GeoNet QuakeSearch API.

The GeoNet project makes all its data and images freely available.
Please ensure you have read and understood our Data Policy and
Disclaimer before using any of these services.  The QuakeSearch
API provides access the New Zealand earthquake catalogue, allows
the user to search quakes using temporal, spatial, depth and
magnitude constraints.

## Compression

The response for a query can be compressed. If your client can handle a
compressed response then the reduced download size is a great benefit.
Gzip compression is supported. You can request a compressed response by
including gzip in your Accept-Encoding header.

## Bugs

The code that provide these services is available at
[https://github.com/GeoNet/quakesearch-go](https://github.com/GeoNet/quakesearch-go)
If you believe you have found a bug please raise an issue or pull request there.

## Endpoints

The following endpoints are available:

* [Csv](#csv)
* [GeoJSON](#geojson)
* [Gml](#gml)
* [Kml](#kml)

## Csv ## {#csv}

Get Quakes in CSV format for specified time, location, magnitude and depth.

    [GET] /csv?bbox=(bbox)&minmag=(minmag)&maxmag=(maxmag)&mindepth=(mindepth)&maxdepth=(maxdepth)&startdate=(startdate)&enddate=(enddate)

### Accept Version

    text/csv

### Parameters

bbox
:   *optional* - The bbox for querying quakes, e.g., 170.29907,-45.22074,175.14404,-41.12075.

limit
:   *optional* - The maximum number of quakes to return.

mindepth
:   *optional* - The minimum depth in km for querying quakes.

maxdepth
:   *optional* - The max depth in km for querying quakes.

minmag
:   *optional* - The min magnitude for querying quakes.

maxmag
:   *optional* - The max magnitude for querying quakes.

region
:   *optional* - The region for querying quakes, e.g., canterbury

startdate
:   *optional* - The Start date for querying quakes.

enddate
:   *optional* - The End date for querying quakes.

### Response

CSV text with the following fields for all earthquakes matching the query

publicid
eventtype
origintime
modificationtime
longitude
latitude
magnitude
depth
magnitudetype
depthtype
evaluationmethod
evaluationstatus
earthmodel
usedphasecount
usedstationcount
magnitudestationcount
minimumdistance
azimuthalgap
originerror
magnitudeuncertainty


### Examples

* [/csv?bbox=163.60840,-49.18170,182.98828,-32.28713](/csv?bbox=163.60840,-49.18170,182.98828,-32.28713)
* [/csv?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=2&maxmag=8&mindepth=0&maxdepth=100](/csv?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=2&maxmag=8&mindepth=0&maxdepth=100)
* [/csv?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2000-3-4T2:00:00&enddate=2016-4-4T4:00:00](/csv?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2000-3-4T2:00:00&enddate=2016-4-4T4:00:00)

## GeoJSON ## {#geojson}

Get Quakes in GeoJSON format for specified time, location, magnitude and depth.

    [GET] /geojson?bbox=(bbox)&minmag=(minmag)&maxmag=(maxmag)&mindepth=(mindepth)&maxdepth=(maxdepth)&startdate=(startdate)&enddate=(enddate)

### Accept Version

    application/vnd.geo+json

### Parameters

bbox
:   *optional* - The bbox for querying quakes, e.g., 170.29907,-45.22074,175.14404,-41.12075.

limit
:   *optional* - The maximum number of quakes to return.

mindepth
:   *optional* - The minimum depth in km for querying quakes.

maxdepth
:   *optional* - The max depth in km for querying quakes.

minmag
:   *optional* - The min magnitude for querying quakes.

maxmag
:   *optional* - The max magnitude for querying quakes.

region
:   *optional* - The region for querying quakes, e.g., canterbury

startdate
:   *optional* - The Start date for querying quakes.

enddate
:   *optional* - The End date for querying quakes.

### Response

GeoJSON with the following example output for all earthquakes matching the query,
in this case two matching quakes.

    {
      "features": [
        {
          "properties": {
            "magnitudestationcount": 3,
            "originerror": 0.14312,
            "magnitudetype": "ML",
            "magnitude": 7.1,
            "depthtype": "operator assigned",
            "depth": 11.0426,
            "modificationtime": "2013-01-15T09:32:00.000Z",
            "origintime": "2010-09-03T16:35:41.000Z",
            "eventtype": "earthquake",
            "publicid": "3366146",
            "evaluationmethod": "GROPE",
            "evaluationstatus": "reviewed",
            "earthmodel": "nz1d",
            "usedphasecount": 21,
            "usedstationcount": 16,
            "minimumdistance": 0.06,
            "azimuthalgap": 82,
            "magnitudeuncertainty": 0.13
          },
          "geometry": {
            "coordinates": [
              172.16794,
              -43.52731
            ],
            "type": "Point"
          },
          "type": "Feature"
        },
        {
          "properties": {
            "magnitudestationcount": 8,
            "originerror": 0.79328,
            "magnitudetype": "ML",
            "magnitude": 5.003,
            "depthtype": "operator assigned",
            "depth": 20,
            "modificationtime": "2010-01-13T22:14:00.000Z",
            "origintime": "2003-08-21T16:50:04.000Z",
            "eventtype": "earthquake",
            "publicid": "2190619",
            "evaluationmethod": "GROPE",
            "evaluationstatus": "reviewed",
            "earthmodel": "nz1d",
            "usedphasecount": 11,
            "usedstationcount": 8,
            "minimumdistance": 0.59,
            "azimuthalgap": 274,
            "magnitudeuncertainty": 0.159
          },
          "geometry": {
            "coordinates": [
              167.09016,
              -45.03076
            ],
            "type": "Point"
          },
          "type": "Feature"
        }
      ],
      "type": "FeatureCollection"
    }


### Examples

* [/geojson?bbox=163.60840,-49.18170,182.98828,-32.28713](/geojson?bbox=163.60840,-49.18170,182.98828,-32.28713)
* [/geojson?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=2&maxmag=8&mindepth=0&maxdepth=100](/geojson?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=2&maxmag=8&mindepth=0&maxdepth=100)
* [/geojson?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2000-3-4T2:00:00&enddate=2016-4-4T4:00:00](/geojson?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2000-3-4T2:00:00&enddate=2016-4-4T4:00:00)

## Gml ## {#gml}

Get Quakes in GML format for specified time, location, magnitude and depth.

    [GET] /gml?bbox=(bbox)&minmag=(minmag)&maxmag=(maxmag)&mindepth=(mindepth)&maxdepth=(maxdepth)&startdate=(startdate)&enddate=(enddate)

### Accept Version

    application/xml

### Parameters

bbox
:   *optional* - The bbox for querying quakes, e.g., 170.29907,-45.22074,175.14404,-41.12075.

limit
:   *optional* - The maximum number of quakes to return.

mindepth
:   *optional* - The minimum depth in km for querying quakes.

maxdepth
:   *optional* - The max depth in km for querying quakes.

minmag
:   *optional* - The min magnitude for querying quakes.

maxmag
:   *optional* - The max magnitude for querying quakes.

region
:   *optional* - The region for querying quakes, e.g., canterbury

startdate
:   *optional* - The Start date for querying quakes.

enddate
:   *optional* - The End date for querying quakes.

### Response

GML formatted XML as in the following example

    <wfs:FeatureCollection xmlns:wfs="http://www.opengis.net/wfs" xmlns:gml="http://www.opengis.net/gml" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:geonet="http://geonet.org.nz" xsi:schemaLocation="http://geonet.org.nz http://geonet.org.nz/quakes http://www.opengis.net/wfs http://schemas.opengis.net/wfs/1.0.0/WFS-basic.xsd">
    <gml:boundedBy>
    <gml:Box srsName="http://www.opengis.net/gml/srs/epsg.xml#4326">
    <gml:coordinates decimal="." cs="," ts=" ">163.60840,-49.18170 182.98828,-32.28713</gml:coordinates>
    </gml:Box>
    </gml:boundedBy>
    <gml:featureMember>
    <geonet:quake fid="quake.3366146">
    <gml:boundedBy>
    <gml:Point srsName="EPSG:4326">
    <gml:coordinates>172.167939999999987,-43.52731</gml:coordinates>
    </gml:Point>
    </gml:boundedBy>
    <geonet:publicid>3366146</geonet:publicid>
    <geonet:origintime>2010-09-03T16:35:41.000Z</geonet:origintime>
    <geonet:latitude>-43.52731</geonet:latitude>
    <geonet:longitude>172.16794</geonet:longitude>
    <geonet:eventtype>earthquake</geonet:eventtype>
    <geonet:modificationtime>2013-01-15T09:32:00.000Z</geonet:modificationtime>
    <geonet:depth>11.0426</geonet:depth>
    <geonet:depthtype>operator assigned</geonet:depthtype>
    <geonet:magnitude>7.1</geonet:magnitude>
    <geonet:magnitudetype>ML</geonet:magnitudetype>
    <geonet:evaluationmethod>GROPE</geonet:evaluationmethod>
    <geonet:evaluationstatus>reviewed</geonet:evaluationstatus>
    <geonet:earthmodel>nz1d</geonet:earthmodel>
    <geonet:usedphasecount>21</geonet:usedphasecount>
    <geonet:usedstationcount>16</geonet:usedstationcount>
    <geonet:minimumdistance>0.06</geonet:minimumdistance>
    <geonet:azimuthalgap>82</geonet:azimuthalgap>
    <geonet:magnitudeuncertainty>0.13</geonet:magnitudeuncertainty>
    <geonet:originerror>0.14312</geonet:originerror>
    <geonet:magnitudestationcount>3</geonet:magnitudestationcount>
    <geonet:origin_geom>
    <gml:Point srsName="EPSG:4326">
    <gml:coordinates>172.167939999999987,-43.52731</gml:coordinates>
    </gml:Point>
    </geonet:origin_geom>
    </geonet:quake>
    </gml:featureMember>
    <gml:featureMember>
    <geonet:quake fid="quake.2190619">
    <gml:boundedBy>
    <gml:Point srsName="EPSG:4326">
    <gml:coordinates>167.090159999999997,-45.030760000000001</gml:coordinates>
    </gml:Point>
    </gml:boundedBy>
    <geonet:publicid>2190619</geonet:publicid>
    <geonet:origintime>2003-08-21T16:50:04.000Z</geonet:origintime>
    <geonet:latitude>-45.03076</geonet:latitude>
    <geonet:longitude>167.09016</geonet:longitude>
    <geonet:eventtype>earthquake</geonet:eventtype>
    <geonet:modificationtime>2010-01-13T22:14:00.000Z</geonet:modificationtime>
    <geonet:depth>20</geonet:depth>
    <geonet:depthtype>operator assigned</geonet:depthtype>
    <geonet:magnitude>5.003</geonet:magnitude>
    <geonet:magnitudetype>ML</geonet:magnitudetype>
    <geonet:evaluationmethod>GROPE</geonet:evaluationmethod>
    <geonet:evaluationstatus>reviewed</geonet:evaluationstatus>
    <geonet:earthmodel>nz1d</geonet:earthmodel>
    <geonet:usedphasecount>11</geonet:usedphasecount>
    <geonet:usedstationcount>8</geonet:usedstationcount>
    <geonet:minimumdistance>0.59</geonet:minimumdistance>
    <geonet:azimuthalgap>274</geonet:azimuthalgap>
    <geonet:magnitudeuncertainty>0.159</geonet:magnitudeuncertainty>
    <geonet:originerror>0.79328</geonet:originerror>
    <geonet:magnitudestationcount>8</geonet:magnitudestationcount>
    <geonet:origin_geom>
    <gml:Point srsName="EPSG:4326">
    <gml:coordinates>167.090159999999997,-45.030760000000001</gml:coordinates>
    </gml:Point>
    </geonet:origin_geom>
    </geonet:quake>
    </gml:featureMember>
    <gml:featureMember>
    <geonet:quake fid="quake.1542894">
    <gml:boundedBy>
    <gml:Point srsName="EPSG:4326">
    <gml:coordinates>174.199999999999989,-40.590000000000003</gml:coordinates>
    </gml:Point>
    </gml:boundedBy>
    <geonet:publicid>1542894</geonet:publicid>
    <geonet:origintime>1946-08-25T20:24:46.000Z</geonet:origintime>
    <geonet:latitude>-40.59</geonet:latitude>
    <geonet:longitude>174.2</geonet:longitude>
    <geonet:eventtype>earthquake</geonet:eventtype>
    <geonet:modificationtime>1982-01-28T00:00:00.000Z</geonet:modificationtime>
    <geonet:depth>12</geonet:depth>
    <geonet:depthtype>operator assigned</geonet:depthtype>
    <geonet:magnitude>3.9</geonet:magnitude>
    <geonet:magnitudetype>ML</geonet:magnitudetype>
    <geonet:evaluationmethod>LOCAL</geonet:evaluationmethod>
    <geonet:evaluationstatus>reviewed</geonet:evaluationstatus>
    <geonet:earthmodel>nz1d</geonet:earthmodel>
    <geonet:usedphasecount>5</geonet:usedphasecount>
    <geonet:usedstationcount>3</geonet:usedstationcount>
    <geonet:minimumdistance>0.809</geonet:minimumdistance>
    <geonet:azimuthalgap>208</geonet:azimuthalgap>
    <geonet:magnitudeuncertainty>0.283</geonet:magnitudeuncertainty>
    <geonet:originerror>0.2</geonet:originerror>
    <geonet:magnitudestationcount>2</geonet:magnitudestationcount>
    <geonet:origin_geom>
    <gml:Point srsName="EPSG:4326">
    <gml:coordinates>174.199999999999989,-40.590000000000003</gml:coordinates>
    </gml:Point>
    </geonet:origin_geom>
    </geonet:quake>
    </gml:featureMember>
    </wfs:FeatureCollection>


### Examples

* [/gml?bbox=163.60840,-49.18170,182.98828,-32.28713](/gml?bbox=163.60840,-49.18170,182.98828,-32.28713)
* [/gml?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=2&maxmag=8&mindepth=0&maxdepth=100](/gml?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=2&maxmag=8&mindepth=0&maxdepth=100)
* [/gml?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2000-3-4T2:00:00&enddate=2016-4-4T4:00:00](/gml?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2000-3-4T2:00:00&enddate=2016-4-4T4:00:00)

## Kml ## {#kml}

Get Quakes in KML format for specified time, location, magnitude and depth.

    [GET] /kml?bbox=(bbox)&minmag=(minmag)&maxmag=(maxmag)&mindepth=(mindepth)&maxdepth=(maxdepth)&startdate=(startdate)&enddate=(enddate)

### Accept Version

    application/vnd.google-earth.kml+xml

### Parameters

bbox
:   *optional* - The bbox for querying quakes, e.g., 170.29907,-45.22074,175.14404,-41.12075.

limit
:   *optional* - The maximum number of quakes to return.

mindepth
:   *optional* - The minimum depth in km for querying quakes.

maxdepth
:   *optional* - The max depth in km for querying quakes.

minmag
:   *optional* - The min magnitude for querying quakes.

maxmag
:   *optional* - The max magnitude for querying quakes.

region
:   *optional* - The region for querying quakes, e.g., canterbury

startdate
:   *optional* - The Start date for querying quakes.

enddate
:   *optional* - The End date for querying quakes.

### Response

KML formatted XML with the example format for two matching quakes

    <?xml version="1.0" encoding="UTF-8"?>
    <kml xmlns="http://www.opengis.net/kml/2.2"
            xmlns:atom="http://www.w3.org/2005/Atom"
            xmlns:gx="http://www.google.com/kml/ext/2.2"
            xmlns:ns5="http://www.w3.org/2001/XMLSchema-instance"
            xmlns:xal="urn:oasis:names:tc:ciq:xsdschema:xAL:2.0">
    <Document ><name>3 New Zealand Earthquakes</name>
    <open>1</open>
    <description>New Zealand earthquake as located by the GeoNet project.</description>
    <StyleMap id="0-15"><Pair><key>normal</key><styleUrl>#inactive-0-15</styleUrl></Pair><Pair><key>highlight</key><styleUrl>#active-0-15</styleUrl></Pair></StyleMap>
    <Style id="active-0-15">
    <IconStyle>
    <scale>0.0</scale>
    <heading>0.0</heading>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/depth-0-15.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <hotSpot x="0.5" y="0.5" xunits="fraction" yunits="fraction"/>
    </IconStyle>
    <LabelStyle><scale>1.0</scale></LabelStyle>
    </Style>
    <Style id="inactive-0-15">
    <IconStyle>
    <scale>0.0</scale>
    <heading>0.0</heading>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/depth-0-15.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <hotSpot x="0.5" y="0.5" xunits="fraction" yunits="fraction"/>
    </IconStyle>
    <LabelStyle><scale>0.0</scale></LabelStyle>
    </Style>
    <StyleMap id="15-40"><Pair><key>normal</key><styleUrl>#inactive-15-40</styleUrl></Pair><Pair><key>highlight</key><styleUrl>#active-15-40</styleUrl></Pair></StyleMap>
    <Style id="active-15-40">
    <IconStyle>
    <scale>0.0</scale>
    <heading>0.0</heading>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/depth-15-40.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <hotSpot x="0.5" y="0.5" xunits="fraction" yunits="fraction"/>
    </IconStyle>
    <LabelStyle><scale>1.0</scale></LabelStyle>
    </Style>
    <Style id="inactive-15-40">
    <IconStyle>
    <scale>0.0</scale>
    <heading>0.0</heading>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/depth-15-40.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <hotSpot x="0.5" y="0.5" xunits="fraction" yunits="fraction"/>
    </IconStyle>
    <LabelStyle><scale>0.0</scale></LabelStyle>
    </Style>
    <StyleMap id="40-100"><Pair><key>normal</key><styleUrl>#inactive-40-100</styleUrl></Pair><Pair><key>highlight</key><styleUrl>#active-40-100</styleUrl></Pair></StyleMap>
    <Style id="active-40-100">
    <IconStyle>
    <scale>0.0</scale>
    <heading>0.0</heading>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/depth-40-100.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <hotSpot x="0.5" y="0.5" xunits="fraction" yunits="fraction"/>
    </IconStyle>
    <LabelStyle><scale>1.0</scale></LabelStyle>
    </Style>
    <Style id="inactive-40-100">
    <IconStyle>
    <scale>0.0</scale>
    <heading>0.0</heading>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/depth-40-100.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <hotSpot x="0.5" y="0.5" xunits="fraction" yunits="fraction"/>
    </IconStyle>
    <LabelStyle><scale>0.0</scale></LabelStyle>
    </Style>
    <StyleMap id="100-200"><Pair><key>normal</key><styleUrl>#inactive-100-200</styleUrl></Pair><Pair><key>highlight</key><styleUrl>#active-100-200</styleUrl></Pair></StyleMap>
    <Style id="active-100-200">
    <IconStyle>
    <scale>0.0</scale>
    <heading>0.0</heading>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/depth-100-200.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <hotSpot x="0.5" y="0.5" xunits="fraction" yunits="fraction"/>
    </IconStyle>
    <LabelStyle><scale>1.0</scale></LabelStyle>
    </Style>
    <Style id="inactive-100-200">
    <IconStyle>
    <scale>0.0</scale>
    <heading>0.0</heading>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/depth-100-200.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <hotSpot x="0.5" y="0.5" xunits="fraction" yunits="fraction"/>
    </IconStyle>
    <LabelStyle><scale>0.0</scale></LabelStyle>
    </Style>
    <StyleMap id="200-more"><Pair><key>normal</key><styleUrl>#inactive-200-more</styleUrl></Pair><Pair><key>highlight</key><styleUrl>#active-200-more</styleUrl></Pair></StyleMap>
    <Style id="active-200-more">
    <IconStyle>
    <scale>0.0</scale>
    <heading>0.0</heading>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/depth-200-more.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <hotSpot x="0.5" y="0.5" xunits="fraction" yunits="fraction"/>
    </IconStyle>
    <LabelStyle><scale>1.0</scale></LabelStyle>
    </Style>
    <Style id="inactive-200-more">
    <IconStyle>
    <scale>0.0</scale>
    <heading>0.0</heading>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/depth-200-more.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <hotSpot x="0.5" y="0.5" xunits="fraction" yunits="fraction"/>
    </IconStyle>
    <LabelStyle><scale>0.0</scale></LabelStyle>
    </Style>
    <Folder ><name>Screen overlays</name>
    <open>1</open>
    <ScreenOverlay ><drawOrder>0</drawOrder>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/legend-depth.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <overlayXY x="0.0" y="0.0" xunits="pixels" yunits="pixels"/>
    <screenXY  x="10.0" y="20.0" xunits="pixels" yunits="pixels"/>
    <rotationXY x="0.0" y="0.0" xunits="pixels" yunits="pixels"/>
    <rotation>0.0</rotation>
    </ScreenOverlay><FolderSimpleExtensionGroup ns5:type="ScreenOverlayType"><drawOrder>0</drawOrder>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.4/images/kml/legend-depth.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <overlayXY x="0.0" y="0.0" xunits="pixels" yunits="pixels"/>
    <screenXY x="10.0" y="20.0" xunits="pixels" yunits="pixels"/>
    <rotationXY x="0.0" y="0.0" xunits="pixels" yunits="pixels"/>
    <rotation>0.0</rotation>
    </FolderSimpleExtensionGroup><ScreenOverlay ><drawOrder>0</drawOrder>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.3/images/logo-gns.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <screenXY x="0.9" y="0.2" xunits="fraction" yunits="fraction"/>
    <rotation>0.0</rotation>
    </ScreenOverlay><FolderSimpleExtensionGroup ns5:type="ScreenOverlayType"><drawOrder>0</drawOrder>
    <Icon><href>http://static.geonet.org.nz/geonet-2.0.3/images/logo-gns.png</href><refreshInterval>0.0</refreshInterval>
    <viewBoundScale>0.0</viewBoundScale>
    <viewRefreshTime>0.0</viewRefreshTime>
    </Icon>
    <screenXY x="0.9" y="0.2" xunits="fraction" yunits="fraction"/>
    <rotation>0.0</rotation>
    </FolderSimpleExtensionGroup></Folder><Folder ><name>Mag 7+ earthquakes</name>
    <Placemark>
    <name>quake.3366146</name>
    <TimeStamp><when>2010-09-03T16:35:41.000Z</when></TimeStamp>
    <styleUrl>#0-15</styleUrl>
    <Style>
    <IconStyle>
    <scale>1.4</scale>
    <heading>0.0</heading>
    </IconStyle>
    </Style>
    <Point>
    <coordinates>172.16794,-43.52731</coordinates>
    </Point>
    <ExtendedData>
    <Data name="Public Id">
    <value>3366146</value>
    </Data>
    <Data name="Universal Time">
    <value>September 03 2010 at 4:35:41 pm</value>
    </Data>
    <Data name="NZ Standard Time">
    <value>Saturday, 04 September 2010 at 4:35:41 am</value>
    </Data>
    <Data name="Focal Depth (km)">
    <value>11.0426</value>
    </Data>
    <Data name="Magnitude">
    <value>7.1</value>
    </Data>
    <Data name="Magnitude Type">
    <value>ML</value>
    </Data>
    <Data name="Depth Type">
    <value>operator assigned</value>
    </Data>
    <Data name="Evaluation Method">
    <value>GROPE</value>
    </Data>
    <Data name="Evaluation Status">
    <value>reviewed</value>
    </Data>
    <Data name="Earth Model">
    <value>nz1d</value>
    </Data>
    <Data name="Used Face Count">
    <value>21</value>
    </Data>
    <Data name="Used station Count">
    <value>16</value>
    </Data>
    <Data name="Magnitude station Count">
    <value>3</value>
    </Data>
    <Data name="Minimum Distance">
    <value>0.06</value>
    </Data>
    <Data name="Azimuthal Gap">
    <value>82</value>
    </Data>
    <Data name="Origin Error">
    <value>0.14312</value>
    </Data>
    <Data name="Magnitude Uncertainty">
    <value>0.13</value>
    </Data>
    </ExtendedData>
    </Placemark>
    </Folder><Folder ><name>Mag 5 earthquakes</name>
    <Placemark>
    <name>quake.2190619</name>
    <TimeStamp><when>2003-08-21T16:50:04.000Z</when></TimeStamp>
    <styleUrl>#15-40</styleUrl>
    <Style>
    <IconStyle>
    <scale>1.0</scale>
    <heading>0.0</heading>
    </IconStyle>
    </Style>
    <Point>
    <coordinates>167.09016,-45.03076</coordinates>
    </Point>
    <ExtendedData>
    <Data name="Public Id">
    <value>2190619</value>
    </Data>
    <Data name="Universal Time">
    <value>August 21 2003 at 4:50:04 pm</value>
    </Data>
    <Data name="NZ Standard Time">
    <value>Friday, 22 August 2003 at 4:50:04 am</value>
    </Data>
    <Data name="Focal Depth (km)">
    <value>20</value>
    </Data>
    <Data name="Magnitude">
    <value>5.003</value>
    </Data>
    <Data name="Magnitude Type">
    <value>ML</value>
    </Data>
    <Data name="Depth Type">
    <value>operator assigned</value>
    </Data>
    <Data name="Evaluation Method">
    <value>GROPE</value>
    </Data>
    <Data name="Evaluation Status">
    <value>reviewed</value>
    </Data>
    <Data name="Earth Model">
    <value>nz1d</value>
    </Data>
    <Data name="Used Face Count">
    <value>11</value>
    </Data>
    <Data name="Used station Count">
    <value>8</value>
    </Data>
    <Data name="Magnitude station Count">
    <value>8</value>
    </Data>
    <Data name="Minimum Distance">
    <value>0.59</value>
    </Data>
    <Data name="Azimuthal Gap">
    <value>274</value>
    </Data>
    <Data name="Origin Error">
    <value>0.79328</value>
    </Data>
    <Data name="Magnitude Uncertainty">
    <value>0.159</value>
    </Data>
    </ExtendedData>
    </Placemark>
    </Folder><Folder ><name>Mag 3 earthquakes</name>
    <Placemark>
    <name>quake.1542894</name>
    <TimeStamp><when>1946-08-25T20:24:46.000Z</when></TimeStamp>
    <styleUrl>#0-15</styleUrl>
    <Style>
    <IconStyle>
    <scale>0.6</scale>
    <heading>0.0</heading>
    </IconStyle>
    </Style>
    <Point>
    <coordinates>174.2,-40.59</coordinates>
    </Point>
    <ExtendedData>
    <Data name="Public Id">
    <value>1542894</value>
    </Data>
    <Data name="Universal Time">
    <value>August 25 1946 at 8:24:46 pm</value>
    </Data>
    <Data name="NZ Standard Time">
    <value>Monday, 26 August 1946 at 8:24:46 am</value>
    </Data>
    <Data name="Focal Depth (km)">
    <value>12</value>
    </Data>
    <Data name="Magnitude">
    <value>3.9</value>
    </Data>
    <Data name="Magnitude Type">
    <value>ML</value>
    </Data>
    <Data name="Depth Type">
    <value>operator assigned</value>
    </Data>
    <Data name="Evaluation Method">
    <value>LOCAL</value>
    </Data>
    <Data name="Evaluation Status">
    <value>reviewed</value>
    </Data>
    <Data name="Earth Model">
    <value>nz1d</value>
    </Data>
    <Data name="Used Face Count">
    <value>5</value>
    </Data>
    <Data name="Used station Count">
    <value>3</value>
    </Data>
    <Data name="Magnitude station Count">
    <value>2</value>
    </Data>
    <Data name="Minimum Distance">
    <value>0.809</value>
    </Data>
    <Data name="Azimuthal Gap">
    <value>208</value>
    </Data>
    <Data name="Origin Error">
    <value>0.2</value>
    </Data>
    <Data name="Magnitude Uncertainty">
    <value>0.283</value>
    </Data>
    </ExtendedData>
    </Placemark>
    </Folder></Document></kml>

### Examples

* [/kml?bbox=163.60840,-49.18170,182.98828,-32.28713](/kml?bbox=163.60840,-49.18170,182.98828,-32.28713)
* [/kml?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=2&maxmag=8&mindepth=0&maxdepth=100](/kml?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=2&maxmag=8&mindepth=0&maxdepth=100)
* [/kml?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2000-3-4T2:00:00&enddate=2016-4-4T4:00:00](/kml?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2000-3-4T2:00:00&enddate=2016-4-4T4:00:00)
