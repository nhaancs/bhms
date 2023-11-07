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
    address_level_1_id  SERIAL      NOT NULL INDEX idx_address_level_1_id,
    address_level_2_id  SERIAL      NOT NULL INDEX idx_address_level_2_id,
    address_level_3_id  SERIAL      NOT NULL INDEX idx_address_level_3_id,
	street              TEXT        NOT NULL,
    manager_id          UUID        NOT NULL INDEX idx_manager_id,
    status              TEXT        NOT NULL INDEX idx_status,
    created_at          TIMESTAMP   NOT NULL,
	updated_at          TIMESTAMP   NOT NULL
);

-- Version: 1.03
-- Description: Create table blocks
CREATE TABLE blocks (
	id                  UUID        NOT NULL PRIMARY KEY,
	name                TEXT        NOT NULL,
    property_id         UUID        NOT NULL INDEX idx_property_id,
    created_at          TIMESTAMP   NOT NULL,
	updated_at          TIMESTAMP   NOT NULL
);

-- Version: 1.04
-- Description: Create table floors
CREATE TABLE floors (
	id                  UUID        NOT NULL PRIMARY KEY,
	name                TEXT        NOT NULL,
    block_id            UUID        NOT NULL INDEX idx_block_id,
    property_id         UUID        NOT NULL INDEX idx_property_id,
    created_at          TIMESTAMP   NOT NULL,
	updated_at          TIMESTAMP   NOT NULL
);

-- Version: 1.05
-- Description: Create table units
CREATE TABLE units (
   id                  UUID        NOT NULL PRIMARY KEY,
   name                TEXT        NOT NULL,
   block_id            UUID        NOT NULL INDEX idx_block_id,
   property_id         UUID        NOT NULL INDEX idx_property_id,
   floor_id            UUID        NOT NULL INDEX idx_floor_id,
   created_at          TIMESTAMP   NOT NULL,
   updated_at          TIMESTAMP   NOT NULL
);
