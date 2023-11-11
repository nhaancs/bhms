-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
	id              UUID        NOT NULL PRIMARY KEY,
	first_name      TEXT        NOT NULL,
	last_name       TEXT        NOT NULL,
	phone           TEXT UNIQUE NOT NULL,
	roles           TEXT[]      NOT NULL,
	password_hash   TEXT        NOT NULL,
    status          TEXT        NOT NULL,
    created_at      TIMESTAMP   NOT NULL,
	updated_at      TIMESTAMP   NOT NULL
);

CREATE INDEX idx_status ON users (status);

-- Version: 1.02
-- Description: Create table properties
CREATE TABLE properties (
	id                  UUID        NOT NULL PRIMARY KEY,
	name                TEXT        NOT NULL,
    address_level_1_id  SERIAL      NOT NULL,
    address_level_2_id  SERIAL      NOT NULL,
    address_level_3_id  SERIAL      NOT NULL,
	street              TEXT        NOT NULL,
    manager_id          UUID        NOT NULL,
    status              TEXT        NOT NULL,
    created_at          TIMESTAMP   NOT NULL,
	updated_at          TIMESTAMP   NOT NULL
);

CREATE INDEX idx_address_level_1_id ON properties (address_level_1_id);
CREATE INDEX idx_address_level_2_id ON properties (address_level_2_id);
CREATE INDEX idx_address_level_3_id ON properties (address_level_3_id);
CREATE INDEX idx_manager_id ON properties (manager_id);
CREATE INDEX idx_status ON properties (status);

-- Version: 1.03
-- Description: Create table blocks
CREATE TABLE blocks (
	id                  UUID        NOT NULL PRIMARY KEY,
	name                TEXT        NOT NULL,
    property_id         UUID        NOT NULL,
    status              TEXT        NOT NULL,
    created_at          TIMESTAMP   NOT NULL,
	updated_at          TIMESTAMP   NOT NULL
);

CREATE INDEX idx_property_id ON blocks (property_id);
CREATE INDEX idx_status ON blocks (status);

-- Version: 1.04
-- Description: Create table floors
CREATE TABLE floors (
	id                  UUID        NOT NULL PRIMARY KEY,
	name                TEXT        NOT NULL,
    block_id            UUID        NOT NULL,
    property_id         UUID        NOT NULL,
    status              TEXT        NOT NULL,
    created_at          TIMESTAMP   NOT NULL,
	updated_at          TIMESTAMP   NOT NULL
);

CREATE INDEX idx_property_id ON floors (property_id);
CREATE INDEX idx_block_id ON floors (block_id);
CREATE INDEX idx_status ON floors (status);

-- Version: 1.05
-- Description: Create table units
CREATE TABLE units (
   id                  UUID        NOT NULL PRIMARY KEY,
   name                TEXT        NOT NULL,
   block_id            UUID        NOT NULL,
   property_id         UUID        NOT NULL,
   floor_id            UUID        NOT NULL,
   status              TEXT        NOT NULL,
   created_at          TIMESTAMP   NOT NULL,
   updated_at          TIMESTAMP   NOT NULL
);

CREATE INDEX idx_property_id ON units (property_id);
CREATE INDEX idx_block_id ON units (block_id);
CREATE INDEX idx_floor_id ON units (floor_id);
CREATE INDEX idx_status ON units (status);
