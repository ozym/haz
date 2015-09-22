DROP SCHEMA IF EXISTS haz CASCADE;

CREATE SCHEMA haz;

-- holds quake history.  This should be restricted to only quakes in the last 365 days.
CREATE TABLE haz.quakehistory (
    -- properties from msg.Quake
    PublicID              TEXT NOT NULL,
    Type                  TEXT NOT NULL,
    AgencyID              TEXT NOT NULL,
    ModificationTime      TIMESTAMP(6) WITH TIME ZONE NOT NULL,
    Time                  TIMESTAMP(6) WITH TIME ZONE NOT NULL,
    Latitude              NUMERIC NOT NULL,
    Longitude             NUMERIC NOT NULL,
    Depth                 NUMERIC NOT NULL,
    DepthType             TEXT NOT NULL,
    MethodID              TEXT NOT NULL,
    EarthModelID          TEXT NOT NULL,
    EvaluationMode        TEXT NOT NULL,
    EvaluationStatus      TEXT NOT NULL,
    UsedPhaseCount        INTEGER NOT NULL,
    UsedStationCount      INTEGER NOT NULL,
    StandardError         NUMERIC NOT NULL,
    AzimuthalGap          NUMERIC NOT NULL,
    MinimumDistance       NUMERIC NOT NULL,
    Magnitude             NUMERIC NOT NULL,
    MagnitudeUncertainty  NUMERIC NOT NULL,
    MagnitudeType         TEXT NOT NULL,
    MagnitudeStationCount INTEGER NOT NULL,
    Site                  TEXT NOT NULL,
    -- everything below here is calculated from the message
    -- ModificationTimeUnixMicro is used as a key for history for CAP. 
    ModificationTimeUnixMicro  BIGINT NOT NULL DEFAULT 0, 
    BackupSite BOOLEAN NOT NULL DEFAULT false,
    Deleted BOOLEAN NOT NULL DEFAULT false,
    Locality TEXT NOT NULL DEFAULT 'unknown',
    Geom GEOGRAPHY(POINT, 4326) NOT NULL,
    Status TEXT NOT NULL DEFAULT 'unknown',
    Quality TEXT NOT NULL DEFAULT 'unknown',
    MMI INTEGER NOT NULL DEFAULT -1,
    Intensity TEXT NOT NULL DEFAULT  'unnoticeable',
    MMID_newzealand INTEGER NOT NULL DEFAULT -1,
    Intensity_newzealand TEXT NOT NULL DEFAULT 'unnoticeable',
    In_newzealand BOOLEAN NOT NULL DEFAULT false
);

-- holds latest information for all quakes.
create table haz.quake AS SELECT * FROM haz.quakehistory;

-- holds latest information for all quakes in the last interval 365 days
create table haz.quakeapi AS SELECT * FROM haz.quakehistory;

ALTER TABLE haz.quake ADD CONSTRAINT quake_publicid_key UNIQUE (PublicID);
ALTER TABLE haz.quakeapi ADD CONSTRAINT quakeapi_publicid_key UNIQUE (PublicID);
ALTER TABLE haz.quakehistory ADD CONSTRAINT quakehistory_publicid_modificationtime_key UNIQUE (PublicID,  ModificationTimeUnixMicro);

CREATE FUNCTION haz.quake_geom() 
RETURNS  TRIGGER AS 
$$
BEGIN 
NEW.geom = ST_GeogFromWKB(st_AsEWKB(st_setsrid(st_makepoint(NEW.longitude, NEW.latitude), 4326)));
NEW.In_newzealand = ST_Covers((select geom from haz.quakeregion where regionname = 'newzealand'),NEW.geom); 
RETURN NEW;  END; 
$$
LANGUAGE plpgsql;

CREATE TRIGGER quake_geom_trigger BEFORE INSERT OR UPDATE ON haz.quake
  FOR EACH ROW EXECUTE PROCEDURE haz.quake_geom();

  CREATE TRIGGER quakeapi_geom_trigger BEFORE INSERT OR UPDATE ON haz.quakeapi
  FOR EACH ROW EXECUTE PROCEDURE haz.quake_geom();

CREATE TRIGGER quakehistory_geom_trigger BEFORE INSERT OR UPDATE ON haz.quakehistory
  FOR EACH ROW EXECUTE PROCEDURE haz.quake_geom();

CREATE INDEX quake_publicid_idx ON haz.quake (PublicID);
CREATE INDEX quake_time_idx ON haz.quake (Time);
CREATE INDEX quakeapi_time_idx ON haz.quakeapi (Time);
CREATE INDEX quakehistory_publicid_idx ON haz.quakehistory (PublicID);

CREATE TABLE haz.quakeregion (
    regionname varchar(255) PRIMARY KEY, 
    title TEXT NOT NULL, 
    groupname TEXT NOT NULL,
    Geom GEOGRAPHY(POLYGON, 4326) NOT NULL
);

INSERT INTO haz.quakeregion VALUES ('newzealand', 'New Zealand', 'region', st_geomfromtext('POLYGON((190 -20, 182 -37, 184 -44, 167 -49, 160 -54, 164 -47, 165 -44, 170 -35, 174 -32, 190 -20))'::text, 4326));

CREATE TABLE haz.soh (
    serverID TEXT PRIMARY KEY, 
    timeReceived timestamp(6)  WITH TIME ZONE NOT NULL
);

CREATE OR REPLACE VIEW haz.quake_search_v1
  AS SELECT
  publicID,
   Type AS eventType,
   time AS originTime,
   modificationTime,
   latitude,
   longitude,
   depth,
   magnitude,
   MethodID AS evaluationMethod,
   evaluationStatus,
   evaluationMode,
   EarthModelID AS earthModel,
   depthType,
   StandardError AS originError, 
   usedPhaseCount,
   usedStationCount, 
   minimumDistance,
   azimuthalGap,  
   magnitudeType,
   magnitudeUncertainty,
   magnitudeStationCount ,
   geom::geometry AS origin_geom
FROM haz.quake
WHERE Deleted != true AND BackupSite != true ORDER BY time DESC;
