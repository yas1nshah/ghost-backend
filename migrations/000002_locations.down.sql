-- DROP INDEXES
DROP INDEX IF EXISTS idx_cities_name;
DROP INDEX IF EXISTS idx_cities_name_ur;
DROP INDEX IF EXISTS idx_cities_pk_id;

DROP INDEX IF EXISTS idx_areas_name;
DROP INDEX IF EXISTS idx_areas_city;
DROP INDEX IF EXISTS idx_areas_popular;
DROP INDEX IF EXISTS idx_areas_url_slug;
DROP INDEX IF EXISTS idx_areas_zone;
DROP INDEX IF EXISTS idx_areas_pk_id;

DROP INDEX IF EXISTS idx_registrations_name;
DROP INDEX IF EXISTS idx_registrations_name_ur;
DROP INDEX IF EXISTS idx_registrations_type;

-- DROP TABLES
DROP TABLE IF EXISTS registrations;
DROP TABLE IF EXISTS areas;
DROP TABLE IF EXISTS cities;
