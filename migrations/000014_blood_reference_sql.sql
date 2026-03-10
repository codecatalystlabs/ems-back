-- +goose Up
-- +goose StatementBegin
CREATE TABLE blood_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO blood_products (code, name, description) VALUES
('WB', 'Whole Blood', 'Whole blood'),
('PRBC', 'Packed Red Blood Cells', 'Packed red blood cells'),
('FFP', 'Fresh Frozen Plasma', 'Fresh frozen plasma'),
('PLT', 'Platelets', 'Platelets')
ON CONFLICT (code) DO NOTHING;

CREATE TABLE blood_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    rhesus TEXT NOT NULL CHECK (rhesus IN ('POSITIVE', 'NEGATIVE')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO blood_groups (code, rhesus) VALUES
('A+', 'POSITIVE'), ('A-', 'NEGATIVE'),
('B+', 'POSITIVE'), ('B-', 'NEGATIVE'),
('AB+', 'POSITIVE'), ('AB-', 'NEGATIVE'),
('O+', 'POSITIVE'), ('O-', 'NEGATIVE')
ON CONFLICT (code) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blood_groups;
DROP TABLE IF EXISTS blood_products;
-- +goose StatementEnd
