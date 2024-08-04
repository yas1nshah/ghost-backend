-- DROP TRIGGER
DROP TRIGGER IF EXISTS update_updated_at ON listings;

-- DROP FUNCTION
DROP FUNCTION IF EXISTS update_updated_at_column;

-- DROP INDEXES
DROP INDEX IF EXISTS idx_listings_created_at;
DROP INDEX IF EXISTS idx_listings_updated_at;
DROP INDEX IF EXISTS idx_listings_active;
DROP INDEX IF EXISTS idx_listings_featured;
DROP INDEX IF EXISTS idx_listings_gp_managed;
DROP INDEX IF EXISTS idx_listings_gp_certified;
DROP INDEX IF EXISTS idx_listings_gp_yard;
DROP INDEX IF EXISTS idx_listings_make;
DROP INDEX IF EXISTS idx_listings_model;
DROP INDEX IF EXISTS idx_listings_version;
DROP INDEX IF EXISTS idx_listings_year;
DROP INDEX IF EXISTS idx_listings_price;
DROP INDEX IF EXISTS idx_listings_registration;
DROP INDEX IF EXISTS idx_listings_city;
DROP INDEX IF EXISTS idx_listings_area;
DROP INDEX IF EXISTS idx_listings_fuel_type;
DROP INDEX IF EXISTS idx_listings_body_type;
DROP INDEX IF EXISTS idx_listings_color;
DROP INDEX IF EXISTS idx_listings_seller;

-- DROP TABLE
DROP TABLE IF EXISTS listings;


-- DROP TRIGGER
DROP TRIGGER IF EXISTS check_listing_limit ON listings;

-- DROP FUNCTION
DROP FUNCTION IF EXISTS check_and_deduct_listing_limit();

-- DROP TRIGGER
DROP TRIGGER IF EXISTS check_featured_limit_on_update ON listings;

-- DROP FUNCTION
DROP FUNCTION IF EXISTS check_and_deduct_featured_limit();
