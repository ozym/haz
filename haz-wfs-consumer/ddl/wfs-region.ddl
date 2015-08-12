-- define regions, should be consistent with qrt.region
DROP TABLE IF EXISTS wfs.region;

CREATE TABLE wfs.region (regionname varchar(255) PRIMARY KEY, title varchar(255), groupname varchar(255));
SELECT addgeometrycolumn('wfs', 'region', 'geom', 4326, 'POLYGON', 2);
