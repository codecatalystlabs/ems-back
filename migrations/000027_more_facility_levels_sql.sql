-- +goose Up
-- +goose StatementBegin
INSERT INTO ref_facility_levels (code, name, rank_no)
VALUES
('CLINIC', 'Clinic', 1),
('DRUGSHOP', 'Drug Shop', 1),
('BCDP', 'Basic Care/Dispensary', 1),
('NBB', 'Nursing/Boarding?', 1),
('RBB', 'RBB', 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM ref_facility_levels
WHERE code IN (
  'CLINIC',
  'DRUGSHOP',
  'BCDP',
  'NBB',
  'RBB'
);
-- +goose StatementEnd
