-- LISTINGS
CREATE TABLE listings (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    active BOOLEAN DEFAULT true,
    featured BOOLEAN DEFAULT false,

    gp_managed BOOLEAN DEFAULT false,
    gp_certified BOOLEAN DEFAULT false,
    gp_yard BOOLEAN DEFAULT false,
    
    gallery JSON,
    
    make INT REFERENCES data_makes(id),
    model INT REFERENCES data_models(id),
    version INT REFERENCES data_versions(id),
    year INT CHECK (year > 1940 AND year < 2025),
    price INT,
    
    registration INT REFERENCES registrations(id),
    city INT REFERENCES cities(id),
    area INT REFERENCES areas(id),
    
    mileage VARCHAR(8),
    transmission INT REFERENCES data_transmissions(id),
    fuel_type INT REFERENCES fuel_types(id),
    engine_capacity VARCHAR(10), -- update it to int
    body_type INT REFERENCES data_body_types(id),
    
    color INT REFERENCES data_colors(id),
    details TEXT,
    
    seller INT REFERENCES users(id),
    upversion INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_listings_created_at ON listings(created_at);
CREATE INDEX idx_listings_updated_at ON listings(updated_at);
CREATE INDEX idx_listings_active ON listings(active);
CREATE INDEX idx_listings_featured ON listings(featured);
CREATE INDEX idx_listings_gp_managed ON listings(gp_managed);
CREATE INDEX idx_listings_gp_certified ON listings(gp_certified);
CREATE INDEX idx_listings_gp_yard ON listings(gp_yard);
CREATE INDEX idx_listings_make ON listings(make);
CREATE INDEX idx_listings_model ON listings(model);
CREATE INDEX idx_listings_version ON listings(version);
CREATE INDEX idx_listings_year ON listings(year);
CREATE INDEX idx_listings_price ON listings(price);
CREATE INDEX idx_listings_registration ON listings(registration);
CREATE INDEX idx_listings_city ON listings(city);
CREATE INDEX idx_listings_area ON listings(area);
CREATE INDEX idx_listings_fuel_type ON listings(fuel_type);
CREATE INDEX idx_listings_body_type ON listings(body_type);
CREATE INDEX idx_listings_color ON listings(color);
CREATE INDEX idx_listings_seller ON listings(seller);

-- TRIGGER FUNCTION for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- TRIGGER for updated_at
CREATE TRIGGER update_updated_at
BEFORE UPDATE ON listings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- TRIGGER FUNCTION to check and deduct listing_limit on insert
CREATE OR REPLACE FUNCTION check_and_deduct_listing_limit()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the current user has listing_limit > 0
    IF (SELECT listing_limit FROM users WHERE id = NEW.seller) <= 0 THEN
        RAISE EXCEPTION 'User does not have enough listing limit.';
    ELSE
        -- Deduct listing_limit by 1
        UPDATE users SET listing_limit = listing_limit - 1 WHERE id = NEW.seller;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- TRIGGER to check and deduct listing_limit on insert
CREATE TRIGGER check_listing_limit
BEFORE INSERT ON listings
FOR EACH ROW
EXECUTE FUNCTION check_and_deduct_listing_limit();

-- TRIGGER FUNCTION to check and deduct featured_limit on update
CREATE OR REPLACE FUNCTION check_and_deduct_featured_limit()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if 'featured' is being updated to true
    IF NEW.featured AND NOT OLD.featured THEN
        -- Check if the user has featured_limit > 0
        IF (SELECT featured_limit FROM users WHERE id = NEW.seller) <= 0 THEN
            RAISE EXCEPTION 'User does not have enough featured limit.';
        ELSE
            -- Deduct featured_limit by 1
            UPDATE users SET featured_limit = featured_limit - 1 WHERE id = NEW.seller;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- TRIGGER to check and deduct featured_limit on update
CREATE TRIGGER check_featured_limit_on_update
BEFORE UPDATE ON listings
FOR EACH ROW
EXECUTE FUNCTION check_and_deduct_featured_limit();
