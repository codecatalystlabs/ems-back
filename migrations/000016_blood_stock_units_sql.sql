-- +goose Up
-- +goose StatementBegin
CREATE TABLE blood_stock_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inventory_site_id UUID NOT NULL REFERENCES blood_inventory_sites(id) ON DELETE CASCADE,
    blood_product_id UUID NOT NULL REFERENCES blood_products(id),
    blood_group_id UUID NOT NULL REFERENCES blood_groups(id),
    unit_count INT NOT NULL DEFAULT 0 CHECK (unit_count >= 0),
    reserved_count INT NOT NULL DEFAULT 0 CHECK (reserved_count >= 0),
    available_count INT GENERATED ALWAYS AS (GREATEST(unit_count - reserved_count, 0)) STORED,
    last_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES users(id),
    UNIQUE(inventory_site_id, blood_product_id, blood_group_id)
);

CREATE INDEX idx_blood_stock_units_site_id ON blood_stock_units(inventory_site_id);
CREATE INDEX idx_blood_stock_units_group_product ON blood_stock_units(blood_group_id, blood_product_id);
CREATE INDEX idx_blood_stock_units_available_count ON blood_stock_units(available_count);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blood_stock_units;
-- +goose StatementEnd
