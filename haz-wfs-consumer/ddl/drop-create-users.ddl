DROP ROLE if exists geonetadmin;
DROP ROLE if exists hazard_w;
DROP ROLE if exists hazard_r;
CREATE ROLE geonetadmin WITH CREATEDB CREATEROLE LOGIN PASSWORD 'test';
CREATE ROLE hazard_w WITH LOGIN PASSWORD 'test';
CREATE ROLE hazard_r WITH LOGIN PASSWORD 'test';

