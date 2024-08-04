-- USERS table definition
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    email_verified BOOLEAN DEFAULT FALSE,
    phone VARCHAR(10) UNIQUE,
    phone_verified BOOLEAN DEFAULT FALSE,
    password_hash BYTEA,
    profile_pic VARCHAR(200) DEFAULT 'default.webp',
    city INT REFERENCES cities(id),
    date_joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    listing_limit INT DEFAULT 1,
    featured_limit INT DEFAULT 0,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_email_verified ON users(email_verified);
CREATE INDEX idx_users_phone_verified ON users(phone_verified);

-- DEALERS table definition
CREATE TABLE dealers (
    user_id INT PRIMARY KEY REFERENCES users(id),
    address VARCHAR(255) UNIQUE,
    timings VARCHAR(255),
    version INTEGER NOT NULL DEFAULT 1
);

-- LISTING PLANS table definition
CREATE TABLE listing_plans (
    id SERIAL PRIMARY KEY,
    listing_limit INT NOT NULL,
    featured_limit INT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_listing_plans_listing_limit ON listing_plans(listing_limit);
CREATE INDEX idx_listing_plans_featured_limit ON listing_plans(featured_limit);

-- TOKENS table definition
CREATE TABLE IF NOT EXISTS tokens (
    hash BYTEA PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope TEXT NOT NULL
);

CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_expiry ON tokens(expiry);
CREATE INDEX idx_tokens_scope ON tokens(scope);

-- PERMISSIONS table definition
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    code TEXT NOT NULL
);

CREATE INDEX idx_permissions_code ON permissions(code);
