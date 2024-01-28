create schema nsdtrc;
use nsdtrc;

Create table members (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    phonenumber VARCHAR(40),
    email VARCHAR(100) UNIQUE,
    website VARCHAR(100),
    region int,
    created DATETIME NOT NULL,
    joined DATETIME
);

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    active boolean NOT NULL default 0,
    verification_hash CHAR(60)
);