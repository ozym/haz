/*
* ddl for wfs: events, origins, magnitudes
* 26/4/2013
*/

DROP SCHEMA IF EXISTS wfs CASCADE;

CREATE SCHEMA wfs;

-- numeric without any precision or scale creates a column in which numeric
-- values of any precision and scale can be stored, up to the implementation
-- limit on precision.
-- http://www.postgresql.org/docs/8.3/static/datatype-numeric.html
--- 1. event table
CREATE TABLE wfs.event (
   publicID varchar(128) PRIMARY KEY,
   eventType varchar(128),
   originTime timestamp(6) WITH TIME ZONE NOT NULL,
   modificationTime timestamp(6) WITH TIME ZONE NOT NULL,
   latitude numeric NOT NULL,
   longitude numeric NOT NULL,
   depth numeric,
   magnitude numeric,
   evaluationMethod varchar(128),
   evaluationStatus varchar(128),
   evaluationMode varchar(50),
   earthModel varchar(128),
   depthType varchar(128),
   originError  numeric,
   usedPhaseCount integer,
   usedStationCount  integer,
   minimumDistance  numeric,
   azimuthalGap  numeric,
   magnitudeType varchar(50),
   magnitudeUncertainty  numeric,
   magnitudeStationCount  integer
);

--- 1.3. gem stuff
SELECT addgeometrycolumn('wfs', 'event', 'origin_geom', 4326, 'POINT', 2);

CREATE FUNCTION wfs.update_origin_geom() RETURNS  TRIGGER AS E' BEGIN NEW.origin_geom = ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude),4326); RETURN NEW;  END; ' LANGUAGE plpgsql;

CREATE TRIGGER origin_geom_trigger BEFORE INSERT OR UPDATE ON wfs.event
  FOR EACH ROW EXECUTE PROCEDURE wfs.update_origin_geom();

CREATE INDEX event_oritime_idx ON wfs.event (originTime);

--- 2. function to add event
-- drop existing
-- DROP FUNCTION IF EXISTS wfs.add_event(publicID_n TEXT, agency_n TEXT, latitude_n NUMERIC, longitude_n NUMERIC, originTime_n TIMESTAMP(6)  WITH TIME ZONE, modificationTime_n TIMESTAMP(6)  WITH TIME ZONE, depth_n NUMERIC, usedPhaseCount_n INT, magnitude_n NUMERIC, magnitudeType_n TEXT, status_n TEXT, type_n TEXT );
-- ori_method_n TEXT, earth_model_n TEXT, depth_type_n TEXT,  ori_err_n NUMERIC,  ori_stns_n INT, dist_min_n NUMERIC, azimuthalGap_n NUMERIC, ori_agency_n TEXT,  mag_err_n NUMERIC, mag_stns_n INT, mag_agency_n TEXT
-- ori_method, earth_model, depth_type,  ori_err,  ori_stns, dist_min, azimuthalGap, ori_agency,  mag_err, mag_stns, mag_agency
 
DROP FUNCTION IF EXISTS  wfs.add_event(publicID_n TEXT, agency_n TEXT, latitude_n NUMERIC, longitude_n NUMERIC, originTime_n TIMESTAMP(6)  WITH TIME ZONE, modificationTime_n TIMESTAMP(6)  WITH TIME ZONE, depth_n NUMERIC, usedPhaseCount_n INT, magnitude_n NUMERIC, magnitudeType_n TEXT, status_n TEXT, type_n TEXT, ori_method_n TEXT, earth_model_n TEXT, depth_type_n TEXT,  ori_err_n NUMERIC,  ori_stns_n INT, dist_min_n NUMERIC, azimuthalGap_n NUMERIC, ori_agency_n TEXT,  mag_err_n NUMERIC, mag_stns_n INT, mag_agency_n TEXT);

CREATE OR REPLACE FUNCTION wfs.add_event(publicID_n TEXT, eventType_n TEXT, originTime_n TIMESTAMP(6)  WITH TIME ZONE,  modificationTime_n TIMESTAMP(6)  WITH TIME ZONE, latitude_n NUMERIC, longitude_n NUMERIC, depth_n NUMERIC, magnitude_n NUMERIC, evaluationMethod_n TEXT, evaluationStatus_n TEXT, evaluationMode_n TEXT, earthModel_n TEXT, depthType_n TEXT, originError_n NUMERIC, usedPhaseCount_n INT, usedStationCount_n INT,minimumDistance_n NUMERIC, azimuthalGap_n NUMERIC, magnitudeType_n TEXT, magnitudeUncertainty_n NUMERIC, magnitudeStationCount_n INT ) RETURNS VOID AS
$$
DECLARE
  tries       INTEGER = 0;
  longitude_n numeric := longitude_n;
BEGIN
  LOOP
  IF longitude_n > 180.0
  THEN
    longitude_n = longitude_n - 360.0;
  END IF;
  UPDATE wfs.event
  SET eventType = eventType_n, originTime = originTime_n, modificationTime = modificationTime_n, latitude = latitude_n, longitude = longitude_n, depth = depth_n, magnitude = magnitude_n, evaluationMethod = evaluationMethod_n, evaluationStatus = evaluationStatus_n, evaluationMode = evaluationMode_n, earthModel = earthModel_n, depthType = depthType_n, originError = originError_n, usedPhaseCount = usedPhaseCount_n, usedStationCount = usedStationCount_n, minimumDistance = minimumDistance_n, azimuthalGap = azimuthalGap_n, magnitudeType = magnitudeType_n, magnitudeUncertainty = magnitudeUncertainty_n, magnitudeStationCount = magnitudeStationCount_n
  WHERE publicID = publicID_n and modificationTime_n > modificationTime;
  IF found
  THEN
    RETURN;
  END IF;

  BEGIN
    INSERT INTO wfs.event (publicID, eventType, originTime, modificationTime, latitude, longitude, depth, magnitude, evaluationMethod, evaluationStatus, evaluationMode, earthModel, depthType, originError, usedPhaseCount, usedStationCount, minimumDistance, azimuthalGap, magnitudeType, magnitudeUncertainty, magnitudeStationCount) VALUES (publicID_n, eventType_n, originTime_n, modificationTime_n, latitude_n, longitude_n, depth_n, magnitude_n, evaluationMethod_n, evaluationStatus_n, evaluationMode_n, earthModel_n, depthType_n, originError_n, usedPhaseCount_n, usedStationCount_n, minimumDistance_n, azimuthalGap_n, magnitudeType_n, magnitudeUncertainty_n, magnitudeStationCount_n);
    RETURN;
    EXCEPTION WHEN unique_violation
    THEN
--  If we get to here the event update is probably old (modificationTime_n <= modificationTime).
--  Loop once more to see if a different insert happend after the update but before
--  our insert.
      tries = tries + 1;
      if tries > 1
      THEN
        RETURN;
      END IF;
  END;
  END LOOP;
END;
$$
LANGUAGE plpgsql;

CREATE TABLE wfs.gt_pk_metadata_table (table_schema varchar(255), table_name varchar(255), pk_column varchar(255), pk_column_idx integer, pk_policy varchar(255), pk_sequence varchar(255));

--- 3. create view
create or replace view wfs.quake_search_v1
  as select * from wfs.event
  where eventType != 'not existing' or eventType is null order by origintime desc;

INSERT INTO wfs.gt_pk_metadata_table(table_schema, table_name, pk_column, pk_column_idx, pk_policy, pk_sequence) VALUES ('wfs', 'quake_search_v1', 'publicid', null, null, null);

INSERT INTO geometry_columns(f_table_catalog, f_table_schema, f_table_name, f_geometry_column, coord_dimension, srid, "type") VALUES ('', 'wfs', 'quake_search_v1', 'origin_geom', 2, 4326, 'POINT');

