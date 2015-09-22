-- This file can be used to remove the region features from the haz DB.
--
alter table haz.quake alter column mmi type int using floor(mmi);
alter table haz.quakeapi alter column mmi type int using floor(mmi);
alter table haz.quakehistory alter column mmi type int using floor(mmi);

alter table haz.quake alter column mmid_newzealand type int using floor(mmid_newzealand);
alter table haz.quakeapi alter column mmid_newzealand type int using floor(mmid_newzealand);
alter table haz.quakehistory alter column mmid_newzealand type int using floor(mmid_newzealand);

alter table haz.quake drop column MMID_aucklandnorthland; 
alter table haz.quake drop column MMID_tongagrirobayofplenty;
alter table haz.quake drop column MMID_gisborne;
alter table haz.quake drop column MMID_hawkesbay;
alter table haz.quake drop column MMID_taranaki;
alter table haz.quake drop column MMID_wellington;
alter table haz.quake drop column MMID_nelsonwestcoast;
alter table haz.quake drop column MMID_canterbury;
alter table haz.quake drop column MMID_fiordland;
alter table haz.quake drop column MMID_otagosouthland;

alter table haz.quakeapi drop column MMID_aucklandnorthland; 
alter table haz.quakeapi drop column MMID_tongagrirobayofplenty;
alter table haz.quakeapi drop column MMID_gisborne;
alter table haz.quakeapi drop column MMID_hawkesbay;
alter table haz.quakeapi drop column MMID_taranaki;
alter table haz.quakeapi drop column MMID_wellington;
alter table haz.quakeapi drop column MMID_nelsonwestcoast;
alter table haz.quakeapi drop column MMID_canterbury;
alter table haz.quakeapi drop column MMID_fiordland;
alter table haz.quakeapi drop column MMID_otagosouthland;

alter table haz.quakehistory drop column MMID_aucklandnorthland; 
alter table haz.quakehistory drop column MMID_tongagrirobayofplenty;
alter table haz.quakehistory drop column MMID_gisborne;
alter table haz.quakehistory drop column MMID_hawkesbay;
alter table haz.quakehistory drop column MMID_taranaki;
alter table haz.quakehistory drop column MMID_wellington;
alter table haz.quakehistory drop column MMID_nelsonwestcoast;
alter table haz.quakehistory drop column MMID_canterbury;
alter table haz.quakehistory drop column MMID_fiordland;
alter table haz.quakehistory drop column MMID_otagosouthland;

alter table haz.quake drop column intensity_aucklandnorthland; 
alter table haz.quake drop column intensity_tongagrirobayofplenty;
alter table haz.quake drop column intensity_gisborne;
alter table haz.quake drop column intensity_hawkesbay;
alter table haz.quake drop column intensity_taranaki;
alter table haz.quake drop column intensity_wellington;
alter table haz.quake drop column intensity_nelsonwestcoast;
alter table haz.quake drop column intensity_canterbury;
alter table haz.quake drop column intensity_fiordland;
alter table haz.quake drop column intensity_otagosouthland;

alter table haz.quakeapi drop column intensity_aucklandnorthland; 
alter table haz.quakeapi drop column intensity_tongagrirobayofplenty;
alter table haz.quakeapi drop column intensity_gisborne;
alter table haz.quakeapi drop column intensity_hawkesbay;
alter table haz.quakeapi drop column intensity_taranaki;
alter table haz.quakeapi drop column intensity_wellington;
alter table haz.quakeapi drop column intensity_nelsonwestcoast;
alter table haz.quakeapi drop column intensity_canterbury;
alter table haz.quakeapi drop column intensity_fiordland;
alter table haz.quakeapi drop column intensity_otagosouthland;

alter table haz.quakehistory drop column intensity_aucklandnorthland; 
alter table haz.quakehistory drop column intensity_tongagrirobayofplenty;
alter table haz.quakehistory drop column intensity_gisborne;
alter table haz.quakehistory drop column intensity_hawkesbay;
alter table haz.quakehistory drop column intensity_taranaki;
alter table haz.quakehistory drop column intensity_wellington;
alter table haz.quakehistory drop column intensity_nelsonwestcoast;
alter table haz.quakehistory drop column intensity_canterbury;
alter table haz.quakehistory drop column intensity_fiordland;
alter table haz.quakehistory drop column intensity_otagosouthland;

alter table haz.quake drop column in_aucklandnorthland; 
alter table haz.quake drop column in_tongagrirobayofplenty;
alter table haz.quake drop column in_gisborne;
alter table haz.quake drop column in_hawkesbay;
alter table haz.quake drop column in_taranaki;
alter table haz.quake drop column in_wellington;
alter table haz.quake drop column in_nelsonwestcoast;
alter table haz.quake drop column in_canterbury;
alter table haz.quake drop column in_fiordland;
alter table haz.quake drop column in_otagosouthland;

alter table haz.quakeapi drop column in_aucklandnorthland; 
alter table haz.quakeapi drop column in_tongagrirobayofplenty;
alter table haz.quakeapi drop column in_gisborne;
alter table haz.quakeapi drop column in_hawkesbay;
alter table haz.quakeapi drop column in_taranaki;
alter table haz.quakeapi drop column in_wellington;
alter table haz.quakeapi drop column in_nelsonwestcoast;
alter table haz.quakeapi drop column in_canterbury;
alter table haz.quakeapi drop column in_fiordland;
alter table haz.quakeapi drop column in_otagosouthland;

alter table haz.quakehistory drop column in_aucklandnorthland; 
alter table haz.quakehistory drop column in_tongagrirobayofplenty;
alter table haz.quakehistory drop column in_gisborne;
alter table haz.quakehistory drop column in_hawkesbay;
alter table haz.quakehistory drop column in_taranaki;
alter table haz.quakehistory drop column in_wellington;
alter table haz.quakehistory drop column in_nelsonwestcoast;
alter table haz.quakehistory drop column in_canterbury;
alter table haz.quakehistory drop column in_fiordland;
alter table haz.quakehistory drop column in_otagosouthland;

drop function haz.quake_geom() cascade;

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

delete from haz.quakeregion where regionname != 'newzealand';
