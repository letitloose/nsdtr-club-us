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

CREATE TABLE stateCode (
    code VARCHAR(5) UNIQUE,
    display VARCHAR(100)
);

CREATE TABLE countryCode (
    code VARCHAR(5) UNIQUE,
    display VARCHAR(100)
);

CREATE TABLE address (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	address1 VARCHAR(255) NOT NULL,
	address2 VARCHAR(255),
	city VARCHAR(100) NOT NULL,
	stateProvince VARCHAR(5),
	zipCode VARCHAR(10),
	country VARCHAR(5)
);

ALTER TABLE address ADD FOREIGN KEY (stateProvince) REFERENCES stateCode(code);
ALTER TABLE address ADD FOREIGN KEY (country) REFERENCES countryCode(code);

-- drop table members;
Create table members (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    phonenumber VARCHAR(40),
	addressID INTEGER,    
    email VARCHAR(100) UNIQUE,
    website VARCHAR(100),
    region int,
    created DATETIME NOT NULL,
    joined DATETIME
);

ALTER TABLE members ADD FOREIGN KEY (addressID) REFERENCES address(id);



