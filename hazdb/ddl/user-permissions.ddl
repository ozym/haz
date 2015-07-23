-- haz and impact are kept as separate schemas with separate write
-- users for security.

DROP ROLE if exists hazard_w;
DROP ROLE if exists hazard_r;
DROP ROLE if exists impact_w;

CREATE ROLE hazard_w WITH LOGIN PASSWORD 'test';
CREATE ROLE hazard_r WITH LOGIN PASSWORD 'test';
CREATE ROLE impact_w WITH LOGIN PASSWORD 'test';

GRANT CONNECT ON DATABASE hazard TO hazard_w;
GRANT USAGE ON SCHEMA haz TO hazard_w;
GRANT ALL ON ALL TABLES IN SCHEMA haz TO hazard_w;
GRANT ALL ON ALL SEQUENCES IN SCHEMA haz TO hazard_w;

GRANT CONNECT ON DATABASE hazard TO hazard_r;
GRANT USAGE ON SCHEMA haz TO hazard_r;
GRANT SELECT ON ALL TABLES IN SCHEMA haz TO hazard_r;
GRANT USAGE ON SCHEMA impact TO hazard_r;
GRANT SELECT ON ALL TABLES IN SCHEMA impact TO hazard_r;

GRANT CONNECT ON DATABASE hazard TO impact_w;
GRANT USAGE ON SCHEMA impact TO impact_w;
GRANT ALL ON ALL TABLES IN SCHEMA impact TO impact_w;
GRANT ALL ON ALL SEQUENCES IN SCHEMA impact TO impact_w;
