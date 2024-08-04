-- CITIES
CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    name_ur VARCHAR(50),
    pk_id INT UNIQUE,
    popular BOOLEAN
);
CREATE INDEX idx_cities_name ON cities(name);
CREATE INDEX idx_cities_name_ur ON cities(name_ur);


-- AREAS
CREATE TABLE areas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    name_ur VARCHAR(100) NOT NULL,
    city INT REFERENCES cities(id),
    popular BOOLEAN,
    url_slug VARCHAR(100),
    zone INT DEFAULT NULL,
    pk_id INT UNIQUE NOT NULL
);
CREATE INDEX idx_areas_name ON areas(name);
CREATE INDEX idx_areas_name_ur ON areas(name_ur);
CREATE INDEX idx_areas_city ON areas(city);
CREATE INDEX idx_areas_popular ON areas(popular);
CREATE INDEX idx_areas_url_slug ON areas(url_slug);



-- REGISTRATIONS
CREATE TABLE registrations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    name_ur VARCHAR(50),
    type VARCHAR(10) CHECK (type IN ('City', 'Province', 'Un-Registered'))
);
CREATE INDEX idx_registrations_name ON registrations(name);
CREATE INDEX idx_registrations_name_ur ON registrations(name_ur);
CREATE INDEX idx_registrations_type ON registrations(type);
