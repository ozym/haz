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
    MMI NUMERIC NOT NULL DEFAULT -1,
    Intensity TEXT NOT NULL DEFAULT  'unnoticeable',
    MMID_newzealand NUMERIC NOT NULL DEFAULT -1,
    MMID_aucklandnorthland NUMERIC NOT NULL DEFAULT -1,
    MMID_tongagrirobayofplenty NUMERIC NOT NULL DEFAULT -1,
    MMID_gisborne NUMERIC NOT NULL DEFAULT -1,
    MMID_hawkesbay NUMERIC NOT NULL DEFAULT -1,
    MMID_taranaki NUMERIC NOT NULL DEFAULT -1,
    MMID_wellington NUMERIC NOT NULL DEFAULT -1,
    MMID_nelsonwestcoast NUMERIC NOT NULL DEFAULT -1,
    MMID_canterbury NUMERIC NOT NULL DEFAULT -1,
    MMID_fiordland NUMERIC NOT NULL DEFAULT -1,
    MMID_otagosouthland NUMERIC NOT NULL DEFAULT -1,
    Intensity_newzealand TEXT NOT NULL DEFAULT 'unnoticeable',
    Intensity_aucklandnorthland TEXT NOT NULL DEFAULT 'unnoticeable',
    Intensity_tongagrirobayofplenty TEXT NOT NULL DEFAULT 'unnoticeable',
    Intensity_gisborne TEXT NOT NULL DEFAULT 'unnoticeable',
    Intensity_hawkesbay TEXT NOT NULL DEFAULT 'unnoticeable',
    Intensity_taranaki TEXT NOT NULL DEFAULT 'unnoticeable',
    Intensity_wellington TEXT NOT NULL DEFAULT 'unnoticeable',
    Intensity_nelsonwestcoast TEXT NOT NULL DEFAULT 'unnoticeable',
    Intensity_canterbury TEXT NOT NULL DEFAULT 'unnoticeable',
    Intensity_fiordland TEXT NOT NULL DEFAULT 'unnoticeable',
    Intensity_otagosouthland TEXT NOT NULL DEFAULT 'unnoticeable',
    In_newzealand BOOLEAN NOT NULL DEFAULT false,
    In_aucklandnorthland BOOLEAN NOT NULL DEFAULT false,
    In_tongagrirobayofplenty BOOLEAN NOT NULL DEFAULT false,
    In_gisborne BOOLEAN NOT NULL DEFAULT false,
    In_hawkesbay BOOLEAN NOT NULL DEFAULT false,
    In_taranaki BOOLEAN NOT NULL DEFAULT false,
    In_wellington BOOLEAN NOT NULL DEFAULT false,
    In_nelsonwestcoast BOOLEAN NOT NULL DEFAULT false,
    In_canterbury BOOLEAN NOT NULL DEFAULT false,
    In_fiordland BOOLEAN NOT NULL DEFAULT false,
    In_otagosouthland BOOLEAN NOT NULL DEFAULT false
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
NEW.In_aucklandnorthland = ST_Covers((select geom from haz.quakeregion where regionname = 'aucklandnorthland'),NEW.geom); 
NEW.In_tongagrirobayofplenty = ST_Covers((select geom from haz.quakeregion where regionname = 'tongagrirobayofplenty'),NEW.geom); 
NEW.In_gisborne = ST_Covers((select geom from haz.quakeregion where regionname = 'gisborne'),NEW.geom); 
NEW.In_hawkesbay = ST_Covers((select geom from haz.quakeregion where regionname = 'hawkesbay'),NEW.geom); 
NEW.In_taranaki = ST_Covers((select geom from haz.quakeregion where regionname = 'taranaki'),NEW.geom); 
NEW.In_wellington = ST_Covers((select geom from haz.quakeregion where regionname = 'wellington'),NEW.geom); 
NEW.In_nelsonwestcoast = ST_Covers((select geom from haz.quakeregion where regionname = 'nelsonwestcoast'),NEW.geom); 
NEW.In_canterbury = ST_Covers((select geom from haz.quakeregion where regionname = 'canterbury'),NEW.geom); 
NEW.In_fiordland = ST_Covers((select geom from haz.quakeregion where regionname = 'fiordland'),NEW.geom); 
NEW.In_otagosouthland = ST_Covers((select geom from haz.quakeregion where regionname = 'otagosouthland'),NEW.geom); 
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
INSERT INTO haz.quakeregion VALUES ('aucklandnorthland', 'Auckland and Northland', 'north', st_geomfromtext('POLYGON((173.251 -38.138, 175.583 -38.045, 176.474 -36.379, 174.285 -34.026, 171.857 -34.135, 173.251 -38.138))'::text, 4326));
INSERT INTO haz.quakeregion VALUES ('tongagrirobayofplenty', 'Tongariro and Bay of Plenty', 'north', st_geomfromtext('POLYGON((175.028 -39.526, 175.722 -39.809, 176.931 -38.688, 178.346 -36.770, 176.474 -36.379, 175.583 -38.045, 175.028 -39.526))'::text, 4326));
INSERT INTO haz.quakeregion VALUES ('gisborne', 'Gisborne', 'north', st_geomfromtext('POLYGON((176.931 -38.688, 178.561 -39.274, 179.898 -37.361, 178.346 -36.770, 176.931 -38.688))'::text, 4326));
INSERT INTO haz.quakeregion VALUES ('hawkesbay', 'Hawke''s Bay', 'north', st_geomfromtext('POLYGON((176.931 -38.688, 175.722 -39.809, 177.560 -40.638, 178.561 -39.274, 176.931 -38.688))'::text, 4326));
INSERT INTO haz.quakeregion VALUES ('taranaki', 'Taranaki', 'north', st_geomfromtext('POLYGON((172.004 -39.632, 174.156 -40.456, 175.028 -39.526, 175.583 -38.045, 173.251 -38.138, 172.004 -39.632))'::text, 4326));
INSERT INTO haz.quakeregion VALUES ('wellington', 'Wellington and Marlborough', 'north', st_geomfromtext('POLYGON((172.951 -41.767, 175.748 -42.908, 177.560 -40.638, 175.028 -39.526, 174.109 -40.462, 172.951 -41.767))'::text, 4326));
INSERT INTO haz.quakeregion VALUES ('nelsonwestcoast', 'Nelson and West Coast', 'south', st_geomfromtext('POLYGON((167.399 -43.711, 169.168 -44.668, 169.564 -44.341, 172.001 -42.832, 172.951 -41.767, 174.109 -40.462, 172.004 -39.632, 170.180 -41.892, 167.399 -43.711))'::text, 4326));
INSERT INTO haz.quakeregion VALUES ('canterbury', 'Canterbury', 'south', st_geomfromtext('POLYGON((172.951 -41.767, 172.001 -42.832, 169.564 -44.341, 172.312 -45.412, 175.748 -42.908, 172.951 -41.767))'::text, 4326));
INSERT INTO haz.quakeregion VALUES ('fiordland', 'Fiordland', 'south', st_geomfromtext('POLYGON((164.218 -46.083, 163.787 -47.212, 165.247 -47.827, 169.168 -44.668, 167.399 -43.711, 164.218 -46.083))'::text, 4326));
INSERT INTO haz.quakeregion VALUES ('otagosouthland', 'Otago and Southland', 'south', st_geomfromtext('POLYGON((165.247 -47.827, 169.148 -48.410, 172.312 -45.412, 169.564 -44.341, 169.168 -44.668, 165.247 -47.827))'::text, 4326));

CREATE TABLE haz.soh (
    serverID TEXT PRIMARY KEY, 
    timeReceived timestamp(6)  WITH TIME ZONE NOT NULL
);

CREATE TABLE haz.volcanic_alert_level (
    alert_level integer PRIMARY KEY,
    hazards TEXT NOT NULL,
    activity TEXT NOT NULL
);

CREATE TABLE haz.volcano (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    region GEOGRAPHY(POLYGON, 4326),
    alert_level integer references haz.volcanic_alert_level(alert_level)
);

INSERT INTO haz.volcanic_alert_level VALUES(0, 'Volcanic environment hazards.', 'No volcanic unrest.');
INSERT INTO haz.volcanic_alert_level VALUES(1, 'Volcanic unrest hazards.', 'Minor volcanic unrest.');
INSERT INTO haz.volcanic_alert_level VALUES(2, 'Volcanic unrest hazards, potential for eruption hazards.', 'Moderate to heightened volcanic unrest.');
INSERT INTO haz.volcanic_alert_level VALUES(3, 'Eruption hazards near vent. Note: ash, lava flow, and lahar (mudflow) hazards may impact areas distant from the volcano.', 'Minor volcanic eruption.');
INSERT INTO haz.volcanic_alert_level VALUES(4, 'Eruption hazards on and near volcano. Note: ash, lava flow, and lahar (mudflow) hazards may impact areas distant from the volcano.', 'Moderate volcanic eruption.');
INSERT INTO haz.volcanic_alert_level VALUES(5, 'Eruption hazards on and beyond volcano. Note: ash, lava flow, and lahar (mudflow) hazards may impact areas distant from the volcano.', 'Major volcanic eruption.');

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('aucklandvolcanicfield', 'Auckland Volcanic Field', ST_GeographyFromText('POINT(174.77 -36.985)'::text), 0, 
    ST_GeographyFromText('POLYGON((174.4585197 -37.16746562, 174.4585197 -36.58689239, 175.510701 -36.58689239, 175.510701 -37.16746562, 174.4585197 -37.16746562))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('kermadecislands', 'Kermadec Islands', ST_GeographyFromText('POINT(-177.914 -29.254)'::text), 0, 
    ST_GeographyFromText('POLYGON((-179.0291841 -32.93325524, -179.0291841 -25.70303694, -175.775 -25.70303694, -175.775 -32.93325524, -179.0291841 -32.93325524))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('mayorisland', 'Mayor Island', ST_GeographyFromText('POINT(176.251 -37.286)'::text), 0, 
    ST_GeographyFromText('POLYGON((175.870104 -37.53170262, 175.870104 -37.04070906, 176.6399397 -37.04070906, 176.6399397 -37.53170262, 175.870104 -37.53170262))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('ngauruhoe', 'Ngauruhoe', ST_GeographyFromText('POINT(175.632 -39.156)'::text), 0, 
    ST_GeographyFromText('POLYGON((175.5471825 -39.21615818, 175.5471825 -39.10384673, 175.728312 -39.10384673, 175.728312 -39.21615818, 175.5471825 -39.21615818))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('northland', 'Northland', ST_GeographyFromText('POINT(173.63 -35.395)'::text), 0, 
    ST_GeographyFromText('POLYGON((173.2122957 -36.25470988, 173.2122957 -34.88581459, 175.0724475 -34.88581459, 175.0724475 -36.25470988, 173.2122957 -36.25470988))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('okataina', 'Okataina', ST_GeographyFromText('POINT(176.501 -38.119)'::text), 0, 
    ST_GeographyFromText('POLYGON((176.3158211 -38.33990913, 176.3158211 -37.94704823, 176.8111052 -37.94704823, 176.8111052 -38.33990913, 176.3158211 -38.33990913))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('rotorua', 'Rotorua', ST_GeographyFromText('POINT(176.281 -38.093)'::text), 0, 
    ST_GeographyFromText('POLYGON((176.11533 -38.20135287, 176.11533 -37.97620536, 176.4250812 -37.97620536, 176.4250812 -38.20135287, 176.11533 -38.20135287))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('ruapehu', 'Ruapehu', ST_GeographyFromText('POINT(175.563 -39.281)'::text),1, 
    ST_GeographyFromText('POLYGON((175.3707552 -39.481325, 175.3707552 -39.09468564, 175.7744228 -39.09468564, 175.7744228 -39.481325, 175.3707552 -39.481325))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('taupo', 'Taupo', ST_GeographyFromText('POINT(175.896 -38.784)'::text), 0, 
    ST_GeographyFromText('POLYGON((175.564837 -39.08056833, 175.564837 -38.58664502, 176.2482749 -38.58664502, 176.2482749 -39.08056833, 175.564837 -39.08056833))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('tongariro', 'Tongariro', ST_GeographyFromText('POINT(175.641727 -39.133318)'::text),1, 
    ST_GeographyFromText('POLYGON((175.5689901 -39.17961512, 175.5689901 -39.06727363, 175.7499926 -39.06727363, 175.7499926 -39.17961512, 175.5689901 -39.17961512))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('taranakiegmont', 'Taranaki/Egmont', ST_GeographyFromText('POINT(174.061 -39.298)'::text), 0, 
    ST_GeographyFromText('POLYGON((173.6983776 -39.67527512, 173.6983776 -38.94831596, 174.4993628 -38.94831596, 174.4993628 -39.67527512, 173.6983776 -39.67527512))'::text));

INSERT INTO haz.volcano (id, title, location, alert_level, region)
VALUES ('whiteisland', 'White Island', ST_GeographyFromText('POINT(177.183 -37.521)'::text),1, 
    ST_GeographyFromText('POLYGON((176.6867564 -38.00383212, 176.6867564 -37.33926271, 177.400852 -37.33926271, 177.400852 -38.00383212, 176.6867564 -38.00383212))'::text));


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
