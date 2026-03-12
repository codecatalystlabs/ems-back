-- +goose Up
CREATE TABLE ref_districts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT UNIQUE,
    name TEXT NOT NULL UNIQUE,
    region TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ref_subcounties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    district_id UUID NOT NULL REFERENCES ref_districts(id) ON DELETE CASCADE,
    code TEXT UNIQUE,
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(district_id, name)
);

CREATE INDEX idx_ref_subcounties_district_id ON ref_subcounties(district_id);

CREATE TABLE ref_facility_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    rank_no INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ref_facilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,         -- this is the CSV uid
    name TEXT NOT NULL,
    short_name TEXT,
    nhfr_id TEXT,
    district_id UUID REFERENCES ref_districts(id),
    subcounty_id UUID REFERENCES ref_subcounties(id),
    level_id UUID REFERENCES ref_facility_levels(id),
    ownership TEXT,
    phone TEXT,
    email TEXT,
    address TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    location GEOGRAPHY(POINT, 4326),
    is_dispatch_station BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(name, district_id)
);

CREATE INDEX idx_ref_facilities_district_id ON ref_facilities(district_id);
CREATE INDEX idx_ref_facilities_subcounty_id ON ref_facilities(subcounty_id);
CREATE INDEX idx_ref_facilities_level_id ON ref_facilities(level_id);
CREATE INDEX idx_ref_facilities_nhfr_id ON ref_facilities(nhfr_id);

CREATE TABLE ref_incident_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    requires_transport BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ref_priority_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    color_code TEXT NOT NULL CHECK (color_code IN ('RED', 'ORANGE', 'GREEN')),
    sort_order INT NOT NULL,
    target_response_minutes INT,
    severity_weight INT NOT NULL DEFAULT 0,
    escalation_note TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ref_severity_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    sort_order INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ref_ambulance_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    supports_maternal BOOLEAN NOT NULL DEFAULT FALSE,
    supports_neonatal BOOLEAN NOT NULL DEFAULT FALSE,
    supports_trauma BOOLEAN NOT NULL DEFAULT FALSE,
    supports_critical_care BOOLEAN NOT NULL DEFAULT FALSE,
    supports_referral BOOLEAN NOT NULL DEFAULT TRUE,
    min_crew_count INT NOT NULL DEFAULT 1 CHECK (min_crew_count >= 1),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ref_capabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    capability_type TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO ref_facility_levels (code, name, rank_no) VALUES
('HCII', 'Health Centre II', 2),
('HCIII', 'Health Centre III', 3),
('HCIV', 'Health Centre IV', 4),
('HOSPITAL', 'Hospital', 5),
('RRH', 'Regional Referral Hospital', 6),
('NRH', 'National Referral Hospital', 7)
ON CONFLICT (code) DO NOTHING;

INSERT INTO ref_incident_types (code, name, description, requires_transport) VALUES
('MATERNAL_EMERGENCY', 'Maternal Emergency', 'Pregnancy and childbirth related emergency', TRUE),
('NEONATAL_EMERGENCY', 'Neonatal Emergency', 'Newborn emergency', TRUE),
('TRAUMA', 'Trauma', 'Trauma case', TRUE),
('ACCIDENT', 'Accident', 'Road traffic or other accident', TRUE),
('GENERAL_MEDICAL', 'General Medical', 'General medical emergency', TRUE),
('REFERRAL', 'Referral', 'Facility referral transfer', TRUE),
('LAB_TRANSPORT', 'Lab Transport', 'Laboratory sample transport', FALSE),
('MEDICINE_DELIVERY', 'Medicine Delivery', 'Medicine dispatch', FALSE)
ON CONFLICT (code) DO NOTHING;

INSERT INTO ref_priority_levels (code, name, color_code, sort_order, target_response_minutes, severity_weight, escalation_note) VALUES
('RED', 'Red Priority', 'RED', 1, 15, 100, 'Immediate life-threatening emergency. Dispatch immediately.'),
('ORANGE', 'Orange Priority', 'ORANGE', 2, 30, 60, 'High-risk urgent case. Prioritize rapid dispatch.'),
('GREEN', 'Green Priority', 'GREEN', 3, 60, 20, 'Stable case. Manage in normal queue unless condition changes.')
ON CONFLICT (code) DO NOTHING;

INSERT INTO ref_severity_levels (code, name, sort_order) VALUES
('SEV1', 'Life Threatening', 1),
('SEV2', 'Severe', 2),
('SEV3', 'Moderate', 3),
('SEV4', 'Minor', 4)
ON CONFLICT (code) DO NOTHING;

INSERT INTO ref_ambulance_categories (
    code, name, description, supports_maternal, supports_neonatal, supports_trauma, supports_critical_care, supports_referral, min_crew_count
) VALUES
('BLS', 'Basic Life Support', 'Basic ambulance', TRUE, FALSE, TRUE, FALSE, TRUE, 2),
('ALS', 'Advanced Life Support', 'Advanced ambulance', TRUE, TRUE, TRUE, TRUE, TRUE, 3),
('MATERNAL', 'Maternal Ambulance', 'Maternal and obstetric ambulance', TRUE, FALSE, FALSE, FALSE, TRUE, 2),
('NEONATAL', 'Neonatal Ambulance', 'Neonatal transport ambulance', FALSE, TRUE, FALSE, TRUE, TRUE, 3),
('TRANSFER', 'Patient Transfer', 'Non-critical patient transfer ambulance', FALSE, FALSE, FALSE, FALSE, TRUE, 1),
('MOTORCYCLE', 'Motorcycle Ambulance', 'Motorcycle emergency responder', FALSE, FALSE, FALSE, FALSE, TRUE, 1),
('BOAT', 'Boat Ambulance', 'Water transport ambulance', TRUE, TRUE, TRUE, FALSE, TRUE, 2),
('RURAL_4X4', 'Rural 4x4 Ambulance', 'Rural terrain ambulance', TRUE, FALSE, TRUE, FALSE, TRUE, 2)
ON CONFLICT (code) DO NOTHING;

INSERT INTO ref_capabilities (code, name, description, capability_type) VALUES
('DISPATCH_COORDINATION', 'Dispatch Coordination', 'Can coordinate dispatch', 'USER'),
('DRIVE_AMBULANCE', 'Drive Ambulance', 'Can drive ambulance', 'USER'),
('EMT_BASIC', 'EMT Basic', 'Basic EMT capability', 'USER'),
('EMT_ADVANCED', 'EMT Advanced', 'Advanced EMT capability', 'USER'),
('MATERNAL_SUPPORT', 'Maternal Support', 'Maternal emergency support', 'USER'),
('NEONATAL_SUPPORT', 'Neonatal Support', 'Neonatal emergency support', 'USER'),
('TRAUMA_SUPPORT', 'Trauma Support', 'Trauma support capability', 'USER'),
('BLS_SUPPORT', 'BLS Support', 'Basic life support capability', 'AMBULANCE'),
('ALS_SUPPORT', 'ALS Support', 'Advanced life support capability', 'AMBULANCE')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS ref_capabilities;
DROP TABLE IF EXISTS ref_ambulance_categories;
DROP TABLE IF EXISTS ref_severity_levels;
DROP TABLE IF EXISTS ref_priority_levels;
DROP TABLE IF EXISTS ref_incident_types;
DROP TABLE IF EXISTS ref_facilities;
DROP TABLE IF EXISTS ref_facility_levels;
DROP TABLE IF EXISTS ref_districts;