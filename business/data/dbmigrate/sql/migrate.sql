-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
    id              UUID NOT NULL,
	username        TEXT NOT NULL,
	email           TEXT UNIQUE NOT NULL,
	bio             TEXT NOT NULL,
	image           TEXT NOT NULL,
	password_hash   TEXT NOT NULL,
	created_at      TIMESTAMP NOT NULL,
	updated_at      TIMESTAMP NOT NULL,

	PRIMARY KEY (id)
);

