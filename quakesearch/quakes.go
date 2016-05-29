package main

// TODO you can use b.WriteString to avoid all the casting b.Write([]byte(...)).

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/GeoNet/weft"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	empty_param_value = -1000
	GML_BBOX_NZ       = "164,-49 -176, -32"
	// TODO looks like this isn't used?  If it is then please use a protocol less url.
	GEONET_ASSET_URL  = "http://static.geonet.org.nz/"
	MAX_QUAKES_NUMBER = 20000 //for each search

	datePattern        = "(\\d{4})-(\\d{1,2})-(\\d{1,2})"
	dateHourPattern    = "(\\d{4})-(\\d{1,2})-(\\d{1,2}) (\\d{1,2})"
	dateTimePattern    = "(\\d{4})-(\\d{1,2})-(\\d{1,2}) (\\d{1,2}):(\\d{1,2}):(\\d{1,2})"
	dateTimePatternISO = "(\\d{4})-(\\d{1,2})-(\\d{1,2})T(\\d{1,2}):(\\d{1,2}):(\\d{1,2})"

	ISO_DATE_FORMAT     = "2006-01-02"
	ISO_TIME_FORMAT     = "2006-1-2T15:04:05"
	RFC3339_FORMAT      = "2006-01-02T15:04:05.999Z"
	UTC_KML_TIME_FORMAT = "January 02 2006 at 3:04:05 pm"
	NZ_KML_TIME_FORMAT  = "Monday, 02 January 2006 at 3:04:05 pm"

	CONTENT_TYPE_XML = "application/xml"
	CONTENT_TYPE_KML = "application/vnd.google-earth.kml+xml"

	CONTENT_TYPE_GeoJSON = "application/vnd.geo+json"
	CONTENT_TYPE_JSON    = "application/json"
	CONTENT_TYPE_CSV     = "text/csv"
)

var (
	NZTzLocation   *time.Location
	optionalParams = []string{"bbox",
		"enddate",
		"limit",
		"maxdepth",
		"maxmag",
		"mindepth",
		"minmag",
		"region",
		"startdate"}
)

func init() {
	//get NZ time zone location
	l, e := time.LoadLocation("NZ")
	if e == nil {
		NZTzLocation = l
	} else {
		NZTzLocation = time.Local
		log.Println("Unable to get NZ timezone, use local time instead!")
	}
}

/**
 * get the rough break point of origintime so that the number of quakes
 * in each time interval <= MAX_QUAKES_NUMBER
 * goal: to limit queries for large amount of data
 *
 */
func getBreakDates(params *QueryParams) []string {
	sql := "select to_char(origintime, 'YYYY-MM') as yrmth, count(*) as count from haz.quake_search_v1"
	sql = getSqlQueryString(sql, params)
	sql1 := sql + " group by yrmth order by yrmth desc"

	var date time.Time
	var dateStr string
	endDate := params.enddate
	if endDate != "" {
		if _, err := time.Parse(ISO_TIME_FORMAT, endDate); err == nil {
			dateStr = endDate
		} else {
			log.Println("err", err)
			date = time.Now()
			dateStr = date.Format(ISO_DATE_FORMAT)
		}
	} else {
		date = time.Now()
		dateStr = date.Format(ISO_DATE_FORMAT)
	}

	breakDates := make([]string, 0)
	breakDates = append(breakDates, dateStr)

	rows, err := db.Query(sql1)

	if err != nil {
		return breakDates
	}
	defer rows.Close()
	numberQuakes := 0

	for rows.Next() { //21 fields
		var ( //note the null values
			ymth  string
			count int
		)
		err := rows.Scan(&ymth, &count)
		if err != nil {
			return breakDates
		}
		numberQuakes += count
		dateStr = ymth + "-01"
		if numberQuakes >= MAX_QUAKES_NUMBER {
			breakDates = append(breakDates, dateStr)
			numberQuakes = 0
			dateStr = ""
		}
	}

	//add start date
	startdDate := params.startdate
	if startdDate != "" {
		breakDates = append(breakDates, startdDate)
	} else if dateStr != "" {
		if !contains(breakDates, dateStr) {
			breakDates = append(breakDates, dateStr)
		}
	}

	return breakDates
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

/**
 * get the number of quakes with breaking dates when the number is large
 */
func getQuakesCount(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {

	//1. check query parameters
	if res := weft.CheckQuery(r, []string{}, optionalParams); !res.Ok {
		return res
	}

	v := r.URL.Query()
	sqlString := `select count(*) from haz.quake_search_v1`

	params := getQueryParams(v)

	sqlString = getSqlQueryString(sqlString, params)
	var count int
	err := db.QueryRow(sqlString).Scan(&count)
	if err != nil {
		return weft.InternalServerError(err)
	}

	resp := "{\"count\":" + strconv.Itoa(count)

	if count > MAX_QUAKES_NUMBER { //get break dates
		breakDates := getBreakDates(params)
		if len(breakDates) > 0 {
			resp += ", \n\"dates\":["
			for n, date := range breakDates {
				if n > 0 {
					resp += ","
				}
				resp += "\"" + date + "\""
			}
			resp += "]\n"
		}
	}
	resp += "}"

	h.Set("Content-Type", CONTENT_TYPE_JSON)
	b.WriteString(resp)
	return &weft.StatusOK

}

/**
* ideally to use go kml library, but they are too basic, without screen overlay and style map
* so use string content instead.
* kml?region=canterbury&startdate=2010-6-29T21:00:00&enddate=2015-7-29T22:00:00
 */
func getQuakesKml(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	//1. check query parameters
	if res := weft.CheckQuery(r, []string{}, optionalParams); !res.Ok {
		return res
	}

	v := r.URL.Query()
	sqlString := `select publicid, eventtype, to_char(origintime, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') AS origintime,
     latitude, longitude, depth, depthtype, magnitude, magnitudetype, evaluationmethod, evaluationstatus,
     evaluationmode, earthmodel, usedphasecount,usedstationcount, minimumdistance, azimuthalgap, magnitudeuncertainty,
     originerror, magnitudestationcount, to_char(modificationtime, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') AS modificationtime
     from haz.quake_search_v1 `

	params := getQueryParams(v)

	sqlString = getSqlQueryString(sqlString, params)

	rows, err := db.Query(sqlString)

	if err != nil {
		return weft.InternalServerError(err)
	}
	defer rows.Close()
	count := 0

	allQuakeFolders := make(map[string]*Folder)

	for rows.Next() { //21 fields
		var ( //note the null values
			publicid              string
			origintime            string
			modificationtime      sql.NullString
			eventtype             sql.NullString
			latitude              float64
			longitude             float64
			depth                 sql.NullFloat64
			depthtype             sql.NullString
			magnitude             sql.NullFloat64
			magnitudetype         sql.NullString
			evaluationmethod      sql.NullString
			evaluationstatus      sql.NullString
			evaluationmode        sql.NullString
			earthmodel            sql.NullString
			usedphasecount        sql.NullInt64
			usedstationcount      sql.NullInt64
			minimumdistance       sql.NullFloat64
			azimuthalgap          sql.NullFloat64
			magnitudeuncertainty  sql.NullFloat64
			originerror           sql.NullFloat64
			magnitudestationcount sql.NullInt64
		)

		err := rows.Scan(&publicid, &eventtype, &origintime, &latitude, &longitude, &depth, &depthtype,
			&magnitude, &magnitudetype, &evaluationmethod, &evaluationstatus,
			&evaluationmode, &earthmodel, &usedphasecount, &usedstationcount,
			&minimumdistance, &azimuthalgap, &magnitudeuncertainty, &originerror, &magnitudestationcount,
			&modificationtime,
		)
		if err != nil {
			return weft.InternalServerError(err)
		}
		count++

		mag := 0.0
		if magnitude.Valid {
			mag = magnitude.Float64
		}
		dep := 0.0
		if depth.Valid {
			dep = depth.Float64
		}

		iconSt := NewIconStyle(getKmlIconSize(mag), 0.0)
		style := NewStyle("", iconSt, nil)
		quakePm := NewPlacemark("quake."+publicid, origintime, NewPoint(latitude, longitude))
		quakePm.SetStyleUrl(getKmlStyleUrl(dep))
		quakePm.SetStyle(style)

		exData := NewExtendedData()
		exData.AddData(NewData("Public Id", publicid))

		t, err := time.Parse(RFC3339_FORMAT, origintime)
		if err != nil {
			log.Panic("time format error", err)
			return weft.InternalServerError(err)
		}

		tu := t.In(time.UTC)
		utcTime := tu.Format(UTC_KML_TIME_FORMAT)
		exData.AddData(NewData("Universal Time", utcTime))

		tnz := t.In(NZTzLocation)
		nzTime := tnz.Format(NZ_KML_TIME_FORMAT)

		exData.AddData(NewData("NZ Standard Time", nzTime))

		if depth.Valid {
			exData.AddData(NewData("Focal Depth (km)", fmt.Sprintf("%g", depth.Float64)))
		}

		if magnitude.Valid {
			exData.AddData(NewData("Magnitude", fmt.Sprintf("%g", magnitude.Float64)))
		}

		if magnitudetype.Valid {
			exData.AddData(NewData("Magnitude Type", magnitudetype.String))
		}

		if depthtype.Valid {
			exData.AddData(NewData("Depth Type", depthtype.String))
		}

		if evaluationmethod.Valid {
			exData.AddData(NewData("Evaluation Method", evaluationmethod.String))
		}

		if evaluationstatus.Valid {
			exData.AddData(NewData("Evaluation Status", evaluationstatus.String))
		}

		if evaluationmode.Valid {
			exData.AddData(NewData("Evaluation Mode", evaluationmode.String))
		}

		if earthmodel.Valid {
			exData.AddData(NewData("Earth Model", earthmodel.String))
		}

		if usedphasecount.Valid {
			exData.AddData(NewData("Used Face Count", fmt.Sprintf("%d", usedphasecount.Int64)))
		}

		if usedstationcount.Valid {
			exData.AddData(NewData("Used station Count", fmt.Sprintf("%d", usedstationcount.Int64)))
		}

		if magnitudestationcount.Valid {
			exData.AddData(NewData("Magnitude station Count", fmt.Sprintf("%d", magnitudestationcount.Int64)))
		}

		if minimumdistance.Valid {
			exData.AddData(NewData("Minimum Distance", fmt.Sprintf("%g", minimumdistance.Float64)))
		}

		if azimuthalgap.Valid {
			exData.AddData(NewData("Azimuthal Gap", fmt.Sprintf("%g", azimuthalgap.Float64)))
		}

		if originerror.Valid {
			exData.AddData(NewData("Origin Error", fmt.Sprintf("%g", originerror.Float64)))
		}

		if magnitudeuncertainty.Valid {
			exData.AddData(NewData("Magnitude Uncertainty", fmt.Sprintf("%g", magnitudeuncertainty.Float64)))
		}

		quakePm.SetExtendedData(exData)
		if magnitude.Valid {
			quakeMagClass := getQuakeMagClass(magnitude.Float64)
			quakeFolder := allQuakeFolders[quakeMagClass[0]]

			if quakeFolder == nil {
				quakeFolder = NewFolder("Folder", "")
				quakeFolder.AddFeature(NewSimpleContentFolder("name", quakeMagClass[1]))
			}

			quakeFolder.AddFeature(quakePm)
			allQuakeFolders[quakeMagClass[0]] = quakeFolder
		}

	}

	rows.Close()

	doc := NewDocument(fmt.Sprintf("%d New Zealand Earthquakes", count), "1",
		"New Zealand earthquake as located by the GeoNet project.")
	//1. add style map and style
	for _, depth := range QUAKE_STYLE_DEPTHS {
		styleMap := NewStyleMap(depth)
		pair1 := NewPair("normal", "#inactive-"+depth)
		pair2 := NewPair("highlight", "#active-"+depth)
		styleMap.AddPair(pair1)
		styleMap.AddPair(pair2)
		doc.AddFeature(styleMap)
		doc.AddFeature(createKmlStyle("active-"+depth, depth, 1.0))
		doc.AddFeature(createKmlStyle("inactive-"+depth, depth, 0.0))
	}

	//2. add screen overlays
	screenOverLays := createGnsKmlScreenOverlays()
	//add to doc
	doc.AddFeature(screenOverLays)

	//3. add quakes folder
	for i := len(MAG_CLASSES_KEYS) - 1; i >= 0; i-- {
		folder := allQuakeFolders[MAG_CLASSES_KEYS[i]]
		if folder != nil {
			doc.AddFeature(folder)
		}
	}

	kml := NewKML(doc)
	b.WriteString(kml.Render())

	//w.Header().Set("Content-Type", "application/xml") //test!!
	h.Set("Content-Type", CONTENT_TYPE_KML)
	h.Set("Content-Disposition", `attachment; filename="earthquakes.kml"`)
	return &weft.StatusOK

}

func getQuakesGml(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	//1. check query parameters
	if res := weft.CheckQuery(r, []string{}, optionalParams); !res.Ok {
		return res
	}

	v := r.URL.Query()
	sqlString := `select publicid, eventtype, to_char(origintime, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') AS origintime,
           latitude, longitude, depth, depthtype, magnitude,  magnitudetype, evaluationmethod, evaluationstatus,
           evaluationmode, earthmodel, usedphasecount,usedstationcount, minimumdistance, azimuthalgap, magnitudeuncertainty,
           originerror, magnitudestationcount, to_char(modificationtime, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') AS modificationtime,
           ST_AsGML(origin_geom) as gml from haz.quake_search_v1 `

	params := getQueryParams(v)

	sqlString = getSqlQueryString(sqlString, params)

	rows, err := db.Query(sqlString)

	if err != nil {
		return weft.InternalServerError(err)
	}
	defer rows.Close()

	// var b bytes.Buffer
	eol := []byte("\n")
	bbox1 := getGmlBbox(params.bbox)

	if bbox1 == "" {
		bbox1 = GML_BBOX_NZ
	}
	b.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
    <wfs:FeatureCollection xmlns:wfs="http://www.opengis.net/wfs"
     xmlns:gml="http://www.opengis.net/gml"
     xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
     xmlns:geonet="http://geonet.org.nz"
     xsi:schemaLocation="http://geonet.org.nz http://geonet.org.nz/quakes http://www.opengis.net/wfs http://schemas.opengis.net/wfs/1.0.0/WFS-basic.xsd">
     <gml:boundedBy>
       <gml:Box srsName="http://www.opengis.net/gml/srs/epsg.xml#4326">
          <gml:coordinates decimal="." cs="," ts=" ">` + bbox1 + `</gml:coordinates>
       </gml:Box>
     </gml:boundedBy>`))
	b.Write(eol)

	for rows.Next() {
		var ( //note the null values
			publicid              string
			origintime            string
			modificationtime      sql.NullString
			eventtype             sql.NullString
			latitude              float64
			longitude             float64
			depth                 sql.NullFloat64
			depthtype             sql.NullString
			magnitude             sql.NullFloat64
			magnitudetype         sql.NullString
			evaluationmethod      sql.NullString
			evaluationstatus      sql.NullString
			evaluationmode        sql.NullString
			earthmodel            sql.NullString
			usedphasecount        sql.NullInt64
			usedstationcount      sql.NullInt64
			minimumdistance       sql.NullFloat64
			azimuthalgap          sql.NullFloat64
			magnitudeuncertainty  sql.NullFloat64
			originerror           sql.NullFloat64
			magnitudestationcount sql.NullInt64
			gml                   string
		)

		err := rows.Scan(&publicid, &eventtype, &origintime, &latitude, &longitude, &depth, &depthtype,
			&magnitude, &magnitudetype, &evaluationmethod, &evaluationstatus,
			&evaluationmode, &earthmodel, &usedphasecount, &usedstationcount,
			&minimumdistance, &azimuthalgap, &magnitudeuncertainty, &originerror, &magnitudestationcount,
			&modificationtime, &gml,
		)
		if err != nil {
			return weft.InternalServerError(err)
		}
		b.Write([]byte("<gml:featureMember>\n"))
		b.Write([]byte(fmt.Sprintf("<geonet:quake fid=\"quake.%s\">\n", publicid)))
		//
		b.Write([]byte(fmt.Sprintf("<gml:boundedBy>%s</gml:boundedBy>\n", gml)))

		b.Write([]byte(fmt.Sprintf("<geonet:publicid>%s</geonet:publicid>\n", publicid)))
		b.Write([]byte(fmt.Sprintf("<geonet:origintime>%s</geonet:origintime>\n", origintime)))
		b.Write([]byte(fmt.Sprintf("<geonet:latitude>%g</geonet:latitude>\n", latitude)))
		b.Write([]byte(fmt.Sprintf("<geonet:longitude>%g</geonet:longitude>\n", longitude)))
		if eventtype.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:eventtype>%s</geonet:eventtype>\n", eventtype.String)))
		}
		if modificationtime.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:modificationtime>%s</geonet:modificationtime>\n", modificationtime.String)))
		}
		if depth.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:depth>%g</geonet:depth>\n", depth.Float64)))
		}
		if depthtype.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:depthtype>%s</geonet:depthtype>\n", depthtype.String)))
		}
		if magnitude.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:magnitude>%g</geonet:magnitude>\n", magnitude.Float64)))
		}
		if magnitudetype.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:magnitudetype>%s</geonet:magnitudetype>\n", magnitudetype.String)))
		}
		if evaluationmethod.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:evaluationmethod>%s</geonet:evaluationmethod>\n", evaluationmethod.String)))
		}
		if evaluationstatus.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:evaluationstatus>%s</geonet:evaluationstatus>\n", evaluationstatus.String)))
		}
		if evaluationmode.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:evaluationmode>%s</geonet:evaluationmode>\n", evaluationmode.String)))
		}
		if earthmodel.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:earthmodel>%s</geonet:earthmodel>\n", earthmodel.String)))
		}
		if usedphasecount.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:usedphasecount>%d</geonet:usedphasecount>\n", usedphasecount.Int64)))
		}
		if usedstationcount.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:usedstationcount>%d</geonet:usedstationcount>\n", usedstationcount.Int64)))
		}
		if minimumdistance.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:minimumdistance>%g</geonet:minimumdistance>\n", minimumdistance.Float64)))
		}
		if azimuthalgap.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:azimuthalgap>%g</geonet:azimuthalgap>\n", azimuthalgap.Float64)))
		}
		if magnitudeuncertainty.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:magnitudeuncertainty>%g</geonet:magnitudeuncertainty>\n", magnitudeuncertainty.Float64)))
		}
		if originerror.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:originerror>%g</geonet:originerror>\n", originerror.Float64)))
		}
		if magnitudestationcount.Valid {
			b.Write([]byte(fmt.Sprintf("<geonet:magnitudestationcount>%d</geonet:magnitudestationcount>\n", magnitudestationcount.Int64)))
		}
		//geonet:origin_geom
		b.Write([]byte(fmt.Sprintf("<geonet:origin_geom>%s</geonet:origin_geom>\n", gml)))
		b.Write([]byte("</geonet:quake></gml:featureMember>\n"))
	}

	// TODO there is already a deferred rows.Close() and you don't use the DB again in this func so strictly this is not necessary.
	rows.Close()
	b.Write([]byte(`</wfs:FeatureCollection>`))

	// send result response
	h.Set("Content-Type", CONTENT_TYPE_XML)
	return &weft.StatusOK
}

func getQuakesCsv(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	//1. check query parameters
	if res := weft.CheckQuery(r, []string{}, optionalParams); !res.Ok {
		return res
	}

	v := r.URL.Query()
	//21  fields
	sqlString := `select format('%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s',
               publicid,eventtype,to_char(origintime, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'),
               to_char(modificationtime, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'),longitude, latitude, magnitude,
               depth,magnitudetype, depthtype, evaluationmethod, evaluationstatus, evaluationmode, earthmodel, usedphasecount,
               usedstationcount,magnitudestationcount, minimumdistance,
               azimuthalgap,originerror,magnitudeuncertainty) as csv from haz.quake_search_v1`

	params := getQueryParams(v)

	sqlString = getSqlQueryString(sqlString, params)

	rows, err := db.Query(sqlString)

	if err != nil {
		return weft.InternalServerError(err)
	}
	defer rows.Close()
	defer rows.Close()

	var (
		// b bytes.Buffer
		d string
	)
	eol := []byte("\n")

	b.Write([]byte("publicid,eventtype,origintime,modificationtime,longitude, latitude, magnitude, depth,magnitudetype,depthtype," +
		"evaluationmethod,evaluationstatus,evaluationmode,earthmodel,usedphasecount,usedstationcount,magnitudestationcount,minimumdistance," +
		"azimuthalgap,originerror,magnitudeuncertainty"))
	b.Write(eol)
	for rows.Next() {
		err := rows.Scan(&d)
		if err != nil {
			return weft.InternalServerError(err)
		}
		b.Write([]byte(d))
		b.Write(eol)
	}
	rows.Close()

	// send result response
	h.Set("Content-Disposition", `attachment; filename="earthquakes.csv"`)
	h.Set("Content-Type", CONTENT_TYPE_CSV)
	return &weft.StatusOK
}

//http://hutl14681.gns.cri.nz:8081/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2015-6-27T22:00:00&enddate=2015-7-27T23:00:00
//(r *http.Request, h http.Header, b *bytes.Buffer) *result
func getQuakesGeoJson(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	//1. check query parameters
	if res := weft.CheckQuery(r, []string{}, optionalParams); !res.Ok {
		return res
	}

	v := r.URL.Query()
	sqlString := `select publicid, eventtype, to_char(origintime, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') AS origintime,
              depth, depthtype, magnitude, magnitudetype, evaluationmethod, evaluationstatus,
              evaluationmode, earthmodel, usedphasecount,usedstationcount, minimumdistance, azimuthalgap, magnitudeuncertainty,
              originerror, magnitudestationcount, to_char(modificationtime, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') AS modificationtime,
              ST_AsGeoJSON(origin_geom) as geojson from haz.quake_search_v1`

	params := getQueryParams(v)

	sqlString = getSqlQueryString(sqlString, params)

	if params.limit != empty_param_value {
		sqlString += fmt.Sprintf(" order by origintime desc limit %d", params.limit)
	}

	rows, err := db.Query(sqlString)

	if err != nil {
		return weft.InternalServerError(err)
	}
	defer rows.Close()
	allFeatures := make([]Feature, 0)
	//
	for rows.Next() {
		var ( //note the null values
			publicid              string
			origintime            string
			modificationtime      sql.NullString
			eventtype             sql.NullString
			depth                 sql.NullFloat64
			depthtype             sql.NullString
			magnitude             sql.NullFloat64
			magnitudetype         sql.NullString
			evaluationmethod      sql.NullString
			evaluationstatus      sql.NullString
			evaluationmode        sql.NullString
			earthmodel            sql.NullString
			usedphasecount        sql.NullInt64
			usedstationcount      sql.NullInt64
			minimumdistance       sql.NullFloat64
			azimuthalgap          sql.NullFloat64
			magnitudeuncertainty  sql.NullFloat64
			originerror           sql.NullFloat64
			magnitudestationcount sql.NullInt64
			geojson               string
		)

		err := rows.Scan(&publicid, &eventtype, &origintime, &depth, &depthtype,
			&magnitude, &magnitudetype, &evaluationmethod, &evaluationstatus,
			&evaluationmode, &earthmodel, &usedphasecount, &usedstationcount,
			&minimumdistance, &azimuthalgap, &magnitudeuncertainty, &originerror, &magnitudestationcount,
			&modificationtime, &geojson,
		)
		if err != nil {
			return weft.InternalServerError(err)
		}
		quakeFeature := Feature{Type: "Feature"}
		//get geometry
		var featureGeo FeatureGeometry
		err = json.Unmarshal([]byte(geojson), &featureGeo)
		if err != nil {
			log.Panic("error", err)
			return weft.InternalServerError(err)
		}
		quakeFeature.Geometry = featureGeo
		//get properties
		quakeProp := QuakeProperties{Publicid: publicid,
			Origintime: origintime,
		}
		//only get the non null values
		if eventtype.Valid {
			quakeProp.Eventtype = eventtype.String
		}
		if modificationtime.Valid {
			quakeProp.Modificationtime = modificationtime.String
		}
		if depth.Valid {
			quakeProp.Depth = depth.Float64
		}
		if depthtype.Valid {
			quakeProp.Depthtype = depthtype.String
		}
		if magnitude.Valid {
			quakeProp.Magnitude = magnitude.Float64
		}
		if magnitudetype.Valid {
			quakeProp.Magnitudetype = magnitudetype.String
		}
		if evaluationmethod.Valid {
			quakeProp.Evaluationmethod = evaluationmethod.String
		}
		if evaluationstatus.Valid {
			quakeProp.Evaluationstatus = evaluationstatus.String
		}
		if evaluationmode.Valid {
			quakeProp.Evaluationmode = evaluationmode.String
		}
		if earthmodel.Valid {
			quakeProp.Earthmodel = earthmodel.String
		}
		if usedphasecount.Valid {
			quakeProp.Usedphasecount = usedphasecount.Int64
		}
		if usedstationcount.Valid {
			quakeProp.Usedstationcount = usedstationcount.Int64
		}
		if minimumdistance.Valid {
			quakeProp.Minimumdistance = minimumdistance.Float64
		}
		if azimuthalgap.Valid {
			quakeProp.Azimuthalgap = azimuthalgap.Float64
		}
		if magnitudeuncertainty.Valid {
			quakeProp.Magnitudeuncertainty = magnitudeuncertainty.Float64
		}
		if originerror.Valid {
			quakeProp.Originerror = originerror.Float64
		}
		if magnitudestationcount.Valid {
			quakeProp.Magnitudestationcount = magnitudestationcount.Int64
		}

		quakeFeature.Properties = quakeProp
		allFeatures = append(allFeatures, quakeFeature)
	}
	rows.Close()

	outputJson := GeoJsonFeatureCollection{
		Type:     "FeatureCollection",
		Features: allFeatures,
	}

	// send result response
	h.Set("Content-Type", CONTENT_TYPE_GeoJSON)
	// h.Set("Accept", V1GeoJSON)
	jsonBytes, err := json.Marshal(outputJson)
	if err != nil {
		return weft.InternalServerError(err)
	}

	b.Write(jsonBytes)

	return &weft.StatusOK
}

func getQueryParams(v url.Values) *QueryParams {
	return &QueryParams{
		bbox:      v.Get("bbox"),
		startdate: v.Get("startdate"),
		enddate:   v.Get("enddate"),
		region:    v.Get("region"),
		limit:     parseIntVal(v.Get("limit")),
		maxdepth:  parseFloatVal(v.Get("maxdepth")),
		maxmag:    parseFloatVal(v.Get("maxmag")),
		mindepth:  parseFloatVal(v.Get("mindepth")),
		minmag:    parseFloatVal(v.Get("minmag")),
	}
}

// TODO wrt to sql injection - I think this is nearly ok but if it's != "" and there
// is an error from strconv.Atoi then you should return that error to the client as
// a bad request.
func parseFloatVal(valstring string) float64 {
	if valstring != "" {
		if f, err := strconv.ParseFloat(valstring, 64); err == nil {
			return f
		}
	}
	return empty_param_value
}

// TODO wrt to sql injection - I think this is nearly ok but if it's != "" and there
// is an error from strconv.Atoi then you should return that error to the client as
// a bad request.
func parseIntVal(valstring string) int {
	if valstring != "" {
		if val, err := strconv.Atoi(valstring); err == nil {
			return val
		}
	}
	return empty_param_value
}

/* generate sql query string based on query parameters from url*/
// TODO - how are you avoiding SQL injection?  This looks like a bad security hole to me.
// I can see why you want to do this but you need to inspect the queries first very carefully.
// You really need to check all the numbers convert to number, all the date times to date times and
// that the string values are in some set of allowable ones.
// It looks like there is some paramter checking going on but if there is an error parsing a value
// you should return that error to the client.
func getSqlQueryString(sqlPre string, params *QueryParams) string {
	condition := false
	sql := sqlPre
	startdate1 := checkDateFormat(params.startdate)
	if startdate1 != "" {
		sql += fmt.Sprintf(" WHERE origintime >= '%s'::timestamptz", startdate1)
		condition = true
	}
	enddate1 := checkDateFormat(params.enddate)
	if enddate1 != "" {
		sql += getSqlAndOrWhere(condition)
		sql += fmt.Sprintf("  origintime < '%s'::timestamptz", enddate1)
		condition = true
	}

	// region
	if params.region != "" {
		sql += getSqlAndOrWhere(condition)
		sql += fmt.Sprintf(" ST_Contains(ST_Shift_Longitude((select geom from haz.quakeregion where regionname = '%s')::geometry), origin_geom)  ", params.region)
		condition = true
	} else {
		// change bbox format
		bbox1 := getPgBbox(params.bbox)
		if bbox1 != "" {
			sql += getSqlAndOrWhere(condition)
			sql += fmt.Sprintf(" ST_Contains(ST_SetSRID(ST_Envelope('LINESTRING(%s)'::geometry),4326), origin_geom)", bbox1)
			condition = true
		}
	}

	if params.minmag != empty_param_value {
		sql += getSqlAndOrWhere(condition)
		sql += fmt.Sprintf("  magnitude >=%g", params.minmag)
		condition = true
	}

	if params.maxmag != empty_param_value {
		sql += getSqlAndOrWhere(condition)
		sql += fmt.Sprintf("  magnitude < %g", params.maxmag)
		condition = true
	}

	if params.mindepth != empty_param_value {
		sql += getSqlAndOrWhere(condition)
		sql += fmt.Sprintf("  depth >= %g", params.mindepth)
		condition = true
	}

	if params.maxdepth != empty_param_value {
		sql += getSqlAndOrWhere(condition)
		sql += fmt.Sprintf("  depth < %g", params.maxdepth)
		condition = true
	}

	return sql

}

func getPgBbox(bbox string) string {
	if bbox != "" {
		bboxarray := strings.Split(bbox, ",")
		if len(bboxarray) == 4 {
			return bboxarray[0] + " " + bboxarray[1] + "," + bboxarray[2] + " " + bboxarray[3]
		}
	}
	return ""
}

func getGmlBbox(bbox string) string {
	if bbox != "" {
		bboxarray := strings.Split(bbox, ",")
		if len(bboxarray) == 4 {
			return bboxarray[0] + "," + bboxarray[1] + " " + bboxarray[2] + "," + bboxarray[3]
		}
	}
	return ""
}

// TODO injection again - the pattern matching isn't bound to the begining and end of the string
// ^...$ so it doesn't catch any errors.
// once you've constructed the longer string then try to parse the time as well.
func checkDateFormat(date string) string {
	if date == "" {
		return ""
	} else if patternMatch(dateTimePattern, date) || patternMatch(dateTimePatternISO, date) {
		return date + "UTC"
	} else if patternMatch(dateHourPattern, date) {
		return date + ":00:00UTC"
	} else if patternMatch(datePattern, date) { // add hour
		return date + " 00:00:00UTC"
	}
	return ""
}

func patternMatch(pattern string, str string) bool {
	match, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return match
}

func getSqlAndOrWhere(hasWhere bool) string {
	if hasWhere {
		return " and "
	} else {
		return " where "
	}
}

type GeoJsonFeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string          `json:"type"`
	Geometry   FeatureGeometry `json:"geometry"`
	Properties QuakeProperties `json:"properties"`
}

type QuakeProperties struct {
	Publicid              string  `json:"publicid"`
	Eventtype             string  `json:"eventtype,omitempty"`
	Origintime            string  `json:"origintime"`
	Modificationtime      string  `json:"modificationtime,omitempty"`
	Depth                 float64 `json:"depth"`
	Depthtype             string  `json:"depthtype,omitempty"`
	Magnitude             float64 `json:"magnitude,omitempty"`
	Magnitudetype         string  `json:"magnitudetype,omitempty"`
	Evaluationmethod      string  `json:"evaluationmethod,omitempty"`
	Evaluationstatus      string  `json:"evaluationstatus,omitempty"`
	Evaluationmode        string  `json:"evaluationmode,omitempty"`
	Earthmodel            string  `json:"earthmodel,omitempty"`
	Usedphasecount        int64   `json:"usedphasecount,omitempty"`
	Usedstationcount      int64   `json:"usedstationcount,omitempty"`
	Minimumdistance       float64 `json:"minimumdistance,omitempty"`
	Azimuthalgap          float64 `json:"azimuthalgap,omitempty"`
	Magnitudeuncertainty  float64 `json:"magnitudeuncertainty,omitempty"`
	Originerror           float64 `json:"originerror,omitempty"`
	Magnitudestationcount int64   `json:"magnitudestationcount,omitempty"`
}

type FeatureGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type QueryParams struct {
	limit     int
	bbox      string
	startdate string
	enddate   string
	maxdepth  float64
	maxmag    float64
	mindepth  float64
	minmag    float64
	region    string
}
