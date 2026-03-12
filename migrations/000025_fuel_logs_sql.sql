-- ============================================
-- File: 000025_fuel_logs_sql.sql
-- ============================================
-- +goose Up
CREATE TABLE fuel_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ambulance_id UUID NOT NULL REFERENCES ambulances(id) ON DELETE CASCADE,
    fuel_type TEXT,
    liters NUMERIC(10,2) NOT NULL CHECK (liters > 0),
    cost NUMERIC(12,2) CHECK (cost IS NULL OR cost >= 0),
    odometer_km INT CHECK (odometer_km IS NULL OR odometer_km >= 0),
    station_name TEXT,
    filled_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    filled_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_fuel_logs_ambulance_id_filled_at ON fuel_logs(ambulance_id, filled_at DESC);
CREATE INDEX idx_fuel_logs_filled_at ON fuel_logs(filled_at DESC);

-- +goose Down
DROP TABLE IF EXISTS fuel_logs;

