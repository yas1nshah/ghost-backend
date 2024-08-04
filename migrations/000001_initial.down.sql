-- DROP INDEXES
DROP INDEX IF EXISTS idx_data_colors_name;
DROP INDEX IF EXISTS idx_data_colors_hex_code;
DROP INDEX IF EXISTS idx_data_colors_version_id;

DROP INDEX IF EXISTS idx_data_details_version_id;
DROP INDEX IF EXISTS idx_data_details_model_id;
DROP INDEX IF EXISTS idx_data_details_make_id;

DROP INDEX IF EXISTS idx_data_versions_name;
DROP INDEX IF EXISTS idx_data_versions_url_slug;
DROP INDEX IF EXISTS idx_data_versions_model_id;
DROP INDEX IF EXISTS idx_data_versions_gen_id;

DROP INDEX IF EXISTS idx_data_generations_end_year;
DROP INDEX IF EXISTS idx_data_generations_start_year;
DROP INDEX IF EXISTS idx_data_generations_url_slug;
DROP INDEX IF EXISTS idx_data_generations_model_id;

DROP INDEX IF EXISTS idx_data_models_make_id;
DROP INDEX IF EXISTS idx_data_models_name_ur;
DROP INDEX IF EXISTS idx_data_models_name;
DROP INDEX IF EXISTS idx_data_models_url_slug;

DROP INDEX IF EXISTS idx_data_makes_name_ur;
DROP INDEX IF EXISTS idx_data_makes_name;
DROP INDEX IF EXISTS idx_data_makes_url_slug;

-- DROP TABLES
DROP TABLE IF EXISTS data_colors;
DROP TABLE IF EXISTS data_details;
DROP TABLE IF EXISTS data_versions;
DROP TABLE IF EXISTS data_generations;
DROP TABLE IF EXISTS data_models;
DROP TABLE IF EXISTS data_makes;
DROP TABLE IF EXISTS data_body_types;
DROP TABLE IF EXISTS fuel_types;
DROP TABLE IF EXISTS data_transmissions;
