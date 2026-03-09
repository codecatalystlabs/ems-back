-- ============================================
-- File: 000013_reporting_views.sql
-- ============================================
-- +goose Up
CREATE VIEW vw_incident_summary AS
SELECT
    i.id,
    i.incident_number,
    i.reported_at,
    i.status,
    i.verification_status,
    it.code AS incident_type_code,
    it.name AS incident_type_name,
    pl.code AS priority_code,
    pl.name AS priority_name,
    sl.code AS severity_code,
    sl.name AS severity_name,
    d.name AS district_name,
    f.name AS facility_name,
    da.id AS dispatch_assignment_id,
    da.ambulance_id,
    da.status AS dispatch_status,
    da.assigned_at,
    da.completed_at,
    CASE
        WHEN da.assigned_at IS NOT NULL THEN EXTRACT(EPOCH FROM (da.assigned_at - i.reported_at)) / 60.0
        ELSE NULL
    END AS minutes_to_assignment,
    CASE
        WHEN da.completed_at IS NOT NULL THEN EXTRACT(EPOCH FROM (da.completed_at - i.reported_at)) / 60.0
        ELSE NULL
    END AS minutes_to_completion
FROM incidents i
LEFT JOIN ref_incident_types it ON it.id = i.incident_type_id
LEFT JOIN ref_priority_levels pl ON pl.id = i.priority_level_id
LEFT JOIN ref_severity_levels sl ON sl.id = i.severity_level_id
LEFT JOIN ref_districts d ON d.id = i.district_id
LEFT JOIN ref_facilities f ON f.id = i.facility_id
LEFT JOIN dispatch_assignments da ON da.incident_id = i.id;

CREATE VIEW vw_ambulance_latest_readiness AS
SELECT DISTINCT ON (arc.ambulance_id)
    arc.ambulance_id,
    arc.mechanical_status,
    arc.equipment_status,
    arc.fuel_status,
    arc.oxygen_status,
    arc.tire_status,
    arc.stretcher_status,
    arc.communication_status,
    arc.cleanliness_status,
    arc.dispatch_readiness,
    arc.checked_by,
    arc.checked_at,
    arc.notes
FROM ambulance_readiness_checks arc
ORDER BY arc.ambulance_id, arc.checked_at DESC;

CREATE VIEW vw_dispatch_performance AS
SELECT
    d.name AS district_name,
    COUNT(i.id) AS total_incidents,
    COUNT(da.id) AS total_dispatches,
    AVG(CASE WHEN da.assigned_at IS NOT NULL THEN EXTRACT(EPOCH FROM (da.assigned_at - i.reported_at)) / 60.0 END) AS avg_minutes_to_assignment,
    AVG(CASE WHEN da.arrived_scene_at IS NOT NULL THEN EXTRACT(EPOCH FROM (da.arrived_scene_at - i.reported_at)) / 60.0 END) AS avg_minutes_to_arrival,
    AVG(CASE WHEN da.completed_at IS NOT NULL THEN EXTRACT(EPOCH FROM (da.completed_at - i.reported_at)) / 60.0 END) AS avg_minutes_to_completion
FROM incidents i
LEFT JOIN ref_districts d ON d.id = i.district_id
LEFT JOIN dispatch_assignments da ON da.incident_id = i.id
GROUP BY d.name;

-- +goose Down
DROP VIEW IF EXISTS vw_dispatch_performance;
DROP VIEW IF EXISTS vw_ambulance_latest_readiness;
DROP VIEW IF EXISTS vw_incident_summary;
