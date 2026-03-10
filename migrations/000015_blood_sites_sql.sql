-- +goose Up
-- +goose StatementBegin
CREATE TABLE blood_banks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID REFERENCES ref_facilities(id) ON DELETE SET NULL,
    code TEXT UNIQUE,
    name TEXT NOT NULL,
    district_id UUID REFERENCES ref_districts(id),
    phone TEXT,
    email TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    location GEOGRAPHY(POINT, 4326),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_blood_banks_district_id ON blood_banks(district_id);
CREATE INDEX idx_blood_banks_location ON blood_banks USING GIST(location);

CREATE TABLE blood_inventory_sites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_type TEXT NOT NULL CHECK (site_type IN ('FACILITY', 'BLOOD_BANK')),
    facility_id UUID REFERENCES ref_facilities(id) ON DELETE CASCADE,
    blood_bank_id UUID REFERENCES blood_banks(id) ON DELETE CASCADE,
    code TEXT UNIQUE,
    name TEXT NOT NULL,
    district_id UUID REFERENCES ref_districts(id),
    contact_person_name TEXT,
    contact_phone TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    location GEOGRAPHY(POINT, 4326),
    can_issue_emergency_blood BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_blood_inventory_sites_owner CHECK (
        (facility_id IS NOT NULL AND blood_bank_id IS NULL) OR
        (facility_id IS NULL AND blood_bank_id IS NOT NULL)
    )
);

CREATE INDEX idx_blood_inventory_sites_district_id ON blood_inventory_sites(district_id);
CREATE INDEX idx_blood_inventory_sites_location ON blood_inventory_sites USING GIST(location);
CREATE INDEX idx_blood_inventory_sites_facility_id ON blood_inventory_sites(facility_id);
CREATE INDEX idx_blood_inventory_sites_blood_bank_id ON blood_inventory_sites(blood_bank_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blood_inventory_sites;
DROP TABLE IF EXISTS blood_banks;
-- +goose StatementEnd
