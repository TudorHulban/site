CREATE TABLE countries (
    id VARCHAR(23) PRIMARY KEY, -- e.g., 'ROU'
    name VARCHAR(100) NOT NULL,
    is_eu BOOLEAN DEFAULT FALSE
);

CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    district VARCHAR(100),       -- e.g., 'Municipiul Iasi' or 'Stockholms län'
    country_id VARCHAR(3) REFERENCES countries(id),
    CONSTRAINT unique_city_per_district_country UNIQUE (name, district, country_id)
);

-- The Parent Entities (e.g., Google LLC, Tele2 Sverige AB, DIGI ROMANIA S.A.)
CREATE TABLE asn_entities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL -- The organization name
);

-- The Specific Network Routes (e.g., AS1257, AS8708)
CREATE TABLE asn_numbers (
    id SERIAL PRIMARY KEY,
    asn VARCHAR(20) UNIQUE NOT NULL, -- e.g., 'AS8708'
    entity_id INT REFERENCES asn_entities(id),
    country_id VARCHAR(2) REFERENCES countries(id)
);

-- ASN Table
CREATE TABLE asn_info (
    id SERIAL PRIMARY KEY,
    asn VARCHAR(20) UNIQUE NOT NULL,
    entity VARCHAR(255) NOT NULL,
    country_id VARCHAR(2) REFERENCES countries(id)
);