-- ============================================
-- File: 000024_facilities_hierarchy_sql.sql
-- ============================================
-- +goose Up
CREATE TABLE regions (
    region_uid TEXT PRIMARY KEY,
    region TEXT NOT NULL
);

CREATE TABLE districts (
    district_uid TEXT PRIMARY KEY,
    region_uid TEXT NOT NULL REFERENCES regions(region_uid),
    district TEXT NOT NULL
);

CREATE TABLE subcounties (
    subcounty_uid TEXT PRIMARY KEY,
    district_uid TEXT NOT NULL REFERENCES districts(district_uid),
    subcounty TEXT NOT NULL
);

CREATE TABLE facilities (
    facility_uid TEXT PRIMARY KEY,
    subcounty_uid TEXT NOT NULL REFERENCES subcounties(subcounty_uid),
    facility TEXT NOT NULL,
    level TEXT,
    ownership TEXT
);

-- +goose Down
DROP TABLE IF EXISTS facilities;
DROP TABLE IF EXISTS subcounties;
DROP TABLE IF EXISTS districts;
DROP TABLE IF EXISTS regions;

