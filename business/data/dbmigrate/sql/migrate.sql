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
