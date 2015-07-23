-- if source and time exists updates mmi.
-- otherwise insert new values.
CREATE FUNCTION impact.add_intensity_reported(source_n TEXT, longitude_n NUMERIC, latitude_n NUMERIC, time_n TIMESTAMP(6) WITH TIME ZONE, mmi_n INTEGER, comment_n VARCHAR(140)) RETURNS VOID AS
$$
DECLARE
loc GEOGRAPHY = ST_GeogFromWKB(st_AsEWKB(st_setsrid(st_makepoint(longitude_n, latitude_n), 4326)));
tries INTEGER = 0;
BEGIN
LOOP
UPDATE impact.intensity_reported 
SET mmi = mmi_n, comment = comment_n
WHERE source = source_n
AND intensity_reported.time = time_n;
IF found THEN
RETURN;
END IF;

BEGIN
INSERT INTO impact.intensity_reported(source, time, mmi, comment, geohash5, geohash6, location) 
VALUES (source_n, time_n, mmi_n, comment_n, st_geohash(loc, 5), st_geohash(loc, 6), loc);
RETURN;
EXCEPTION WHEN unique_violation THEN
--  Loop once more to see if a different insert happened after the update but before our insert.
tries = tries + 1;
if tries > 1 THEN
RETURN;
END IF;
END;
END LOOP;
END;
$$
LANGUAGE plpgsql;

-- if source does not exist this inserts the new values.
-- if source exists and mmi_n > mmi then updates mmi and time.
-- age of information is not handled here.
-- mmi_n does not have to be for a newer time.  This allows for out of order information
CREATE FUNCTION impact.add_intensity_measured(source_n TEXT, longitude_n NUMERIC, latitude_n NUMERIC, time_n TIMESTAMP(6) WITH TIME ZONE, mmi_n INTEGER) RETURNS VOID AS
$$
DECLARE
tries INTEGER = 0;
BEGIN
LOOP
UPDATE impact.intensity_measured 
SET mmi = mmi_n, time = time_n
WHERE source = source_n
AND intensity_measured.mmi < mmi_n;
IF found THEN
RETURN;
END IF;

DECLARE
loc GEOGRAPHY = ST_GeogFromWKB(st_AsEWKB(st_setsrid(st_makepoint(longitude_n, latitude_n), 4326)));
BEGIN
INSERT INTO impact.intensity_measured(source, time, mmi, location) 
VALUES (source_n, time_n, mmi_n, loc);
RETURN;
EXCEPTION WHEN unique_violation THEN
--  Loop once more to see if a different insert happened after the update but before our insert.
tries = tries + 1;
if tries > 1 THEN
RETURN;
END IF;
END;
END LOOP;
END;
$$
LANGUAGE plpgsql;