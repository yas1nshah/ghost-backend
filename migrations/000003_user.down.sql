-- DROP INDEXES
DROP INDEX IF EXISTS idx_users_name;
DROP INDEX IF EXISTS idx_users_email_verified;
DROP INDEX IF EXISTS idx_users_phone_verified;
DROP INDEX IF EXISTS idx_users_city;

DROP INDEX IF EXISTS idx_listing_plans_listing_limit;
DROP INDEX IF EXISTS idx_listing_plans_featured_limit;

DROP INDEX IF EXISTS idx_tokens_user_id;
DROP INDEX IF EXISTS idx_tokens_expiry;
DROP INDEX IF EXISTS idx_tokens_scope;

DROP INDEX IF EXISTS idx_permissions_code;

-- DROP TABLES
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS listing_plans;
DROP TABLE IF EXISTS dealers;
DROP TABLE IF EXISTS users;
