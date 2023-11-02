-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
	id              UUID        NOT NULL PRIMARY KEY,
	first_name      TEXT        NOT NULL,
	last_name       TEXT        NOT NULL,
	phone           TEXT UNIQUE NOT NULL,
	roles           TEXT[]      NOT NULL,
	password_hash   TEXT        NOT NULL,
    status          TEXT        NOT NULL INDEX idx_status,
    created_at      TIMESTAMP   NOT NULL,
	updated_at      TIMESTAMP   NOT NULL
);

-- Version: 1.02
-- Description: Create table properties
CREATE TABLE properties (
	id                  UUID        NOT NULL PRIMARY KEY,
	name                TEXT        NOT NULL,
    address_level1_id   TEXT        NOT NULL INDEX idx_address_level1_id,
    address_level2_id   TEXT        NOT NULL INDEX idx_address_level2_id,
    address_level3_id   TEXT        NOT NULL INDEX idx_address_level3_id,
	street              TEXT        NOT NULL,
    manager_id          UUID        NOT NULL INDEX idx_manager_id,
    status              TEXT        NOT NULL INDEX idx_status,
    created_at          TIMESTAMP   NOT NULL,
	updated_at          TIMESTAMP   NOT NULL
);
