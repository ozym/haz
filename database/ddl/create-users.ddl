-- haz and impact are kept as separate schemas with separate write
-- users for security.
DROP ROLE if exists geonetadmin;
DROP ROLE if exists hazard_w;
DROP ROLE if exists hazard_r;
DROP ROLE if exists impact_w;

CREATE ROLE geonetadmin WITH CREATEDB CREATEROLE LOGIN PASSWORD 'test';
CREATE ROLE hazard_w WITH LOGIN PASSWORD 'test';
CREATE ROLE hazard_r WITH LOGIN PASSWORD 'test';
CREATE ROLE impact_w WITH LOGIN PASSWORD 'test';

