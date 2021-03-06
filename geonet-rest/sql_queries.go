package main

const quakeProtoSQL = `SELECT time, modificationTime, depth, magnitude, locality,
				floor(mmid_newzealand) as "mmi",
				quality,
				ST_X(geom::geometry) as longitude,
				ST_Y(geom::geometry) as latitude
			FROM haz.quake WHERE publicID = $1`

const quakeHistoryProtoSQL = `SELECT time, modificationTime, depth, magnitude, locality,
				floor(mmid_newzealand) as "mmi",
				quality,
				ST_X(geom::geometry) as longitude,
				ST_Y(geom::geometry) as latitude
			FROM haz.quakehistory WHERE publicid = $1 ORDER BY modificationtime DESC`

const quakesProtoSQL = `SELECT publicid, time, modificationTime, depth, magnitude, locality,
				floor(mmid_newzealand) as "mmi",
				quality,
				ST_X(geom::geometry) as longitude,
				ST_Y(geom::geometry) as latitude
			FROM haz.quakeapi where mmid_newzealand >= $1
			AND In_newzealand = true
			ORDER BY time DESC  limit 100`

const quakeV2SQL = `SELECT row_to_json(fc)
FROM ( SELECT 'FeatureCollection' as type, array_to_json(array_agg(f)) as features
	FROM (SELECT 'Feature' as type,
		ST_AsGeoJSON(q.geom)::json as geometry,
		row_to_json((SELECT l FROM 
			(
				SELECT 
				publicid AS "publicID",
				to_char(time, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as "time",
				depth, 
				magnitude, 
				locality,
				floor(mmid_newzealand) as "mmi",
				quality
				) as l
)) as properties FROM haz.quake as q where publicid = $1 ) As f )  as fc`

const quakesV2SQL = `SELECT row_to_json(fc)
FROM ( SELECT 'FeatureCollection' as type, COALESCE(array_to_json(array_agg(f)), '[]') as features
	FROM (SELECT 'Feature' as type,
		ST_AsGeoJSON(q.geom)::json as geometry,
		row_to_json((SELECT l FROM
			(
				SELECT
				publicid AS "publicID",
				to_char(time, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as "time",
				depth,
				magnitude,
				locality,
				floor(mmid_newzealand) as "mmi",
				quality
				) as l
)) as properties FROM haz.quakeapi as q where mmid_newzealand >= $1
AND In_newzealand = true
ORDER BY time DESC  limit 100 ) as f ) as fc`

const quakeHistoryV2SQL = `SELECT row_to_json(fc)
FROM ( SELECT 'FeatureCollection' as type, array_to_json(array_agg(f)) as features
	FROM (SELECT 'Feature' as type,
		ST_AsGeoJSON(q.geom)::json as geometry,
		row_to_json((SELECT l FROM 
			(
				SELECT 
				publicid AS "publicID",
				to_char(time, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as "time",
				to_char(modificationtime, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as "modificationTime",
				depth, 
				magnitude, 
				locality,
				floor(mmid_newzealand) as "mmi",
				quality
				) as l
)) as properties FROM haz.quakehistory as q where publicid = $1 order by modificationtime desc ) As f )  as fc`

const intensityMeasuredLatestV2SQL = `SELECT row_to_json(fc)
FROM ( SELECT 'FeatureCollection' as type, COALESCE(array_to_json(array_agg(f)), '[]') as features
	FROM (SELECT 'Feature' as type,
		ST_AsGeoJSON(s.location)::json as geometry,
		row_to_json(( select l from 
			( 
				select mmi
				) as l )) 
as properties from (select location, mmi 
	FROM impact.intensity_measured) as s 
) As f )  as fc`

const intenstityReportedLatestV2SQL = `WITH features as (
	select COALESCE(array_to_json(array_agg(fs)), '[]') as features from (SELECT 'Feature' as type,
		ST_AsGeoJSON(s.location)::json as geometry,
		row_to_json(( select l from 
			( 
				select mmi,
				count
				) as l )) 
as properties from (select st_pointfromgeohash(geohash6) as location, 
	max(mmi) as mmi, 
	count(mmi) as count 
	FROM impact.intensity_reported 
	WHERE time >= (now() - interval '60 minutes')
	group by (geohash6)) as s) as fs
), summary as (
	select COALESCE(json_object_agg(summ.mmi, summ.count), '{}') as count_mmi, COALESCE(sum(count), 0) as count
	from (select mmi as mmi, count(*) as count from impact.intensity_reported 
		WHERE time >= (now() - interval '60 minutes')
		group by mmi
		) as summ
)
SELECT row_to_json(fc)
FROM ( SELECT 'FeatureCollection' as type, 
	features.features, 
	summary.count_mmi,
	summary.count
	FROM features, summary )  as fc`

const intenstityReportedWindowV2SQL = `WITH features as (
	select COALESCE(array_to_json(array_agg(fs)), '[]') as features from (SELECT 'Feature' as type,
		ST_AsGeoJSON(s.location)::json as geometry,
		row_to_json(( select l from 
			( 
				select mmi,
				count
				) as l )) 
as properties from (select st_pointfromgeohash(geohash6) as location, 
	max(mmi) as mmi, 
	count(mmi) as count 
	FROM impact.intensity_reported 
	WHERE time >= $1
	AND time <= $2
	group by (geohash6)) as s) as fs
), summary as (
	select COALESCE(json_object_agg(summ.mmi, summ.count), '{}') as count_mmi, COALESCE(sum(count), 0) as count
	from (select mmi as mmi, count(*) as count from impact.intensity_reported 
		WHERE time >= $1
		AND time <= $2
		group by mmi
		) as summ
)
SELECT row_to_json(fc)
FROM ( SELECT 'FeatureCollection' as type, 
	features.features, 
	summary.count_mmi,
	summary.count
	FROM features, summary )  as fc`

/*
	There needs to be an Atom feed entity for every CAP message.
	A CAP message ID is unique for each CAP message (and is not just the quake PublicID).
	1. Find the first time any quake was reviewed within an hour of the quake
	and was strong enough for a CAP message.  This is the first CAP message.
	2. Select the first CAP message and any subsequent reviews or deletes that happened with an hour
	of the quake.  Each of this is a CAP message and gets an entity in the feed.
*/
const capQuakeFeedSQL = `with first_review as 
	(select publicid, min(modificationtimeunixmicro) as modificationtimeunixmicro 
		from haz.quakehistory 
		where status = 'reviewed' 
		AND modificationTime - time < interval '1 hour' 
		AND MMID_newzealand >= $1 
		AND now() - time < interval '48 hours' group by publicid)
select h.publicid, h.modificationtimeunixmicro, h.modificationTime 
from haz.quakehistory h, first_review 
where h.publicid = first_review.publicid 
and h.modificationtimeunixmicro >= first_review.modificationtimeunixmicro 
and status in ('reviewed','deleted') 
AND modificationTime - time < interval '1 hour' ORDER BY time DESC, modificationTime DESC`

const quakeStatsV2SQL = `with mags as (
	select floor(magnitude) as magnitude, time, date_trunc('day', time) as day
	from haz.quakeapi 
	where in_newzealand and not deleted
	), year as (
		select COALESCE(json_object_agg(summ.magnitude, summ.count), '{}') as count_mags 
		from (select magnitude, count(magnitude) as count 
			from mags where  time >= (now() - interval '364 days') group by magnitude) as summ
),
month as (
	select COALESCE(json_object_agg(summ.magnitude, summ.count), '{}') as count_mags 
	from (select magnitude, count(magnitude) as count 
		from mags where  time >= (now() - interval '28 days') group by magnitude) as summ
),
week as (
	select COALESCE(json_object_agg(summ.magnitude, summ.count), '{}') as count_mags 
	from (select magnitude, count(magnitude) as count 
		from mags where  time >= (now() - interval '7 days') group by magnitude) as summ
),
perday as (
	select COALESCE(json_object_agg(summ.day, summ.count), '{}') as "perDay" 
	from (select day, count(day) as count 
		from mags group by day order by day) as summ
)
select row_to_json(f) from (
	select row_to_json(fc) as "magnitudeCount", row_to_json(perday) as "rate" FROM perday, (
		SELECT 
		year.count_mags as "days365",
		month.count_mags as "days28",
		week.count_mags as "days7"
		FROM year, month, week) as fc) as f`

const quakesPerDaySQL = `WITH perday AS (
SELECT date_trunc('day', time) as day
FROM haz.quakeapi
WHERE in_newzealand AND NOT deleted
)
SELECT day, count(day)
FROM perday GROUP BY day ORDER BY day
`

// use this query with fmt.Sprintf to set the days interval e.g.
//   if rows, err = db.Query(fmt.Sprintf(sumMagsSQL, 365)); err != nil {
const sumMagsSQL = `WITH mags AS (
SELECT time, floor(magnitude) AS magnitude
FROM haz.quakeapi
WHERE in_newzealand AND NOT deleted
)
SELECT magnitude, count(magnitude)
FROM mags where  time >= (now() - interval '%d days') group by magnitude
`

const quakesNZWWWSQL = `SELECT row_to_json(fc)
FROM ( SELECT 'FeatureCollection' as type,
COALESCE(array_to_json(array_agg(f)), '[]') as features,
row_to_json(
(SELECT n FROM (SELECT 'EPSG' as type,
	row_to_json((
		SELECT m FROM (SELECT '4326' as "code") as m)) as properties) as n)
) as crs
	FROM (SELECT 'Feature' as type,
		'quake.' || publicid as "id",
		ST_AsGeoJSON(q.geom)::json as geometry,
		'origin_geom' as "geometry_name",
		row_to_json((SELECT l FROM
			(
				SELECT
				publicid AS "publicid",
				to_char(time, 'YYYY-MM-DD HH24:MI:SS.US') as "origintime",
				depth,
				magnitude,
				intensity,
				status,
				agencyid as "agency",
				to_char(modificationtime, 'YYYY-MM-DD HH24:MI:SS.US') AS "updatetime"
				) as l
)) as properties FROM haz.quakeapi as q where mmid_newzealand >= $1
AND In_newzealand = true
ORDER BY time DESC  limit $2 ) as f) as fc`

const quakesWWWSQL = `SELECT row_to_json(fc)
FROM ( SELECT 'FeatureCollection' as type,
COALESCE(array_to_json(array_agg(f)), '[]') as features,
row_to_json(
(SELECT n FROM (SELECT 'EPSG' as type,
	row_to_json((
		SELECT m FROM (SELECT '4326' as "code") as m)) as properties) as n)
) as crs
	FROM (SELECT 'Feature' as type,
		'quake.' || publicid as "id",
		ST_AsGeoJSON(q.geom)::json as geometry,
		'origin_geom' as "geometry_name",
		row_to_json((SELECT l FROM
			(
				SELECT
				publicid AS "publicid",
				to_char(time, 'YYYY-MM-DD HH24:MI:SS.US') as "origintime",
				depth,
				magnitude,
				intensity,
				status,
				agencyid as "agency",
				to_char(modificationtime, 'YYYY-MM-DD HH24:MI:SS.US') AS "updatetime"
				) as l
)) as properties FROM haz.quakeapi as q where mmi >= $1
AND In_newzealand = true
ORDER BY time DESC  limit $2 ) as f) as fc`

const quakeWWWSQL = `SELECT row_to_json(fc)
FROM ( SELECT 'FeatureCollection' as type,
COALESCE(array_to_json(array_agg(f)), '[]') as features,
row_to_json(
(SELECT n FROM (SELECT 'EPSG' as type,
	row_to_json((
		SELECT m FROM (SELECT '4326' as "code") as m)) as properties) as n)
) as crs
	FROM (SELECT 'Feature' as type,
		'quake.' || publicid as "id",
		ST_AsGeoJSON(q.geom)::json as geometry,
		'origin_geom' as "geometry_name",
		row_to_json((SELECT l FROM
			(
				SELECT
				publicid AS "publicid",
				to_char(time, 'YYYY-MM-DD HH24:MI:SS.US') as "origintime",
				depth,
				magnitude,
				intensity,
				status,
				agencyid as "agency",
				to_char(modificationtime, 'YYYY-MM-DD HH24:MI:SS.US') AS "updatetime"
				) as l
)) as properties FROM haz.quakeapi as q where publicid = $1
AND In_newzealand = true
ORDER BY time DESC  limit 100 ) as f) as fc`
