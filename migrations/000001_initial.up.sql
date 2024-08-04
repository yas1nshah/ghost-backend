-- MAKES
CREATE TABLE data_makes (
    id SERIAL PRIMARY KEY,
    url_slug TEXT NOT NULL,
    name VARCHAR(255) NOT NULL,
    name_ur VARCHAR(255),
    published BOOLEAN,
    pk_id INT UNIQUE
);
CREATE INDEX idx_data_makes_url_slug ON data_makes (url_slug);
CREATE INDEX idx_data_makes_name ON data_makes (name);
CREATE INDEX idx_data_makes_name_ur ON data_makes (name_ur);

-- MODELS
CREATE TABLE data_models (
    id SERIAL PRIMARY KEY,
    url_slug TEXT NOT NULL,
    name VARCHAR(255) NOT NULL,
    name_ur VARCHAR(255),
    active BOOLEAN,
    published BOOLEAN,
    popular BOOLEAN,
    pk_id INT UNIQUE,
    make_id INT REFERENCES data_makes(id) ON DELETE CASCADE
);
CREATE INDEX idx_data_models_url_slug ON data_models (url_slug);
CREATE INDEX idx_data_models_name ON data_models (name);
CREATE INDEX idx_data_models_name_ur ON data_models (name_ur);
CREATE INDEX idx_data_models_make_id ON data_models (make_id);

-- GENERATIONS
CREATE TABLE data_generations (
    id SERIAL PRIMARY KEY,
    model_id INT REFERENCES data_models(id) ON DELETE CASCADE,
    start_year INT CHECK (start_year > 1940 AND start_year < 2025),
    end_year CHAR(4),
    name VARCHAR(100),
    is_imported BOOLEAN,
    url_slug TEXT NOT NULL,
    pk_id INT UNIQUE
);
CREATE INDEX idx_data_generations_model_id ON data_generations (model_id);
CREATE INDEX idx_data_generations_url_slug ON data_generations (url_slug);
CREATE INDEX idx_data_generations_start_year ON data_generations (start_year);
CREATE INDEX idx_data_generations_end_year ON data_generations (end_year);

-- VERSIONS
CREATE TABLE data_versions (
    id SERIAL PRIMARY KEY,
    pk_id INT UNIQUE,
    gen_id INT REFERENCES data_generations(id) ON DELETE CASCADE,
    model_id INT REFERENCES data_models(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    name_ur VARCHAR(255),

    url_slug TEXT NOT NULL,
    is_active BOOLEAN,
    is_popular BOOLEAN,
    is_published BOOLEAN,
    is_imported BOOLEAN
);
CREATE INDEX idx_data_versions_gen_id ON data_versions (gen_id);
CREATE INDEX idx_data_versions_model_id ON data_versions (model_id);
CREATE INDEX idx_data_versions_url_slug ON data_versions (url_slug);
CREATE INDEX idx_data_versions_name ON data_versions (name);


-- BODY TYPES
CREATE TABLE data_body_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) 
);

-- TRANSMISSION TYPES
CREATE TABLE data_transmissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) 
);

-- COLORS
CREATE TABLE data_colors (
    id SERIAL PRIMARY KEY,
    version_id INT REFERENCES data_versions(id) ON DELETE CASCADE,
    hex_code CHAR(7),
    name VARCHAR(50),
    name_ur VARCHAR(50),
    mapped_name VARCHAR(50)
);
CREATE INDEX idx_data_colors_version_id ON data_colors (version_id);
CREATE INDEX idx_data_colors_hex_code ON data_colors (hex_code);
CREATE INDEX idx_data_colors_name ON data_colors (name);

-- FUEL TYPES
CREATE TABLE fuel_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50)
);

INSERT INTO fuel_types (name)
VALUES ('Petrol'), ('Electric'), ('Hybrid'), ('Diesel'), ('CNG'), ('LPG');


-- DETAILS
CREATE TABLE data_details (
    id SERIAL PRIMARY KEY,
    is_active BOOLEAN,
    description TEXT,
    launch_year INT CHECK (launch_year < 2025),
    make_id INT REFERENCES data_makes(id) ON DELETE CASCADE,
    model_id INT REFERENCES data_models(id) ON DELETE CASCADE,
    version_id INT REFERENCES data_versions(id) ON DELETE CASCADE,
    imported BOOLEAN,
    transmission_type  INT REFERENCES data_transmissions(id),
    -- VARCHAR(20) CHECK (engine_type IN ('Petrol', 'Electric', 'Hybrid', 'Diesel', 'CNG', 'LPG')),
    fuel_type INT REFERENCES fuel_types(id),
    engine_capacity INT CHECK (engine_capacity < 9999),
    mileage VARCHAR(50),
    price INT,
    body_type INT REFERENCES data_body_types(id)
);
CREATE INDEX idx_data_details_make_id ON data_details (make_id);
CREATE INDEX idx_data_details_model_id ON data_details (model_id);
CREATE INDEX idx_data_details_version_id ON data_details (version_id);
