-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW mv_ambulance_latest_readiness AS
SELECT
    a.id AS ambulance_id,
    a.code,
    a.plate_number,
    a.category_id,
    a.station_facility_id,
    a.district_id,
    a.status,
    a.dispatch_readiness AS ambulance_dispatch_readiness,
    a.is_active,
    arc.id AS readiness_check_id,
    arc.mechanical_status,
    arc.equipment_status,
    arc.fuel_status,
    arc.oxygen_status,
    arc.tire_status,
    arc.stretcher_status,
    arc.communication_status,
    arc.cleanliness_status,
    arc.dispatch_readiness AS readiness_dispatch_readiness,
    arc.checked_at
FROM ambulances a
LEFT JOIN LATERAL (
    SELECT rc.*
    FROM ambulance_readiness_checks rc
    WHERE rc.ambulance_id = a.id
    ORDER BY rc.checked_at DESC
    LIMIT 1
) arc ON TRUE;

CREATE UNIQUE INDEX uq_mv_ambulance_latest_readiness_ambulance_id
ON mv_ambulance_latest_readiness(ambulance_id);

CREATE INDEX idx_mv_ambulance_latest_readiness_district_id
ON mv_ambulance_latest_readiness(district_id);

CREATE INDEX idx_mv_ambulance_latest_readiness_station_facility_id
ON mv_ambulance_latest_readiness(station_facility_id);

CREATE INDEX idx_mv_ambulance_latest_readiness_category_id
ON mv_ambulance_latest_readiness(category_id);


CREATE MATERIALIZED VIEW mv_incident_daily_stats AS
SELECT
    date_trunc('day', i.reported_at)::date AS stat_date,
    i.district_id,
    i.facility_id,
    COUNT(*) AS incidents_total,
    COUNT(*) FILTER (WHERE rpl.code = 'RED') AS red_triage_count,
    COUNT(*) FILTER (
        WHERE i.status IN ('ASSIGNED', 'ENROUTE', 'AT_SCENE', 'TRANSPORTING', 'COMPLETED')
    ) AS transfers_count,
    COUNT(*) FILTER (
        WHERE rit.code = 'MATERNAL_EMERGENCY'
    ) AS mnmci_count,
    COUNT(*) FILTER (
        WHERE rit.code = 'ACCIDENT'
           OR i.summary ILIKE '%rta%'
           OR i.description ILIKE '%rta%'
    ) AS rta_count,
    COUNT(*) FILTER (
        WHERE rit.code = 'HIGHLY_INFECTIOUS'
           OR i.summary ILIKE '%infectious%'
           OR i.description ILIKE '%infectious%'
    ) AS infectious_count
FROM incidents i
LEFT JOIN ref_priority_levels rpl ON rpl.id = i.priority_level_id
LEFT JOIN ref_incident_types rit ON rit.id = i.incident_type_id
GROUP BY 1, 2, 3;

CREATE UNIQUE INDEX uq_mv_incident_daily_stats
ON mv_incident_daily_stats(stat_date, district_id, facility_id);

CREATE INDEX idx_mv_incident_daily_stats_district_id
ON mv_incident_daily_stats(district_id);

CREATE INDEX idx_mv_incident_daily_stats_facility_id
ON mv_incident_daily_stats(facility_id);


CREATE MATERIALIZED VIEW mv_dashboard_daily_summary AS
WITH ambulance_stats AS (
    SELECT
        COALESCE(f.district_id, m.district_id) AS district_id,
        m.station_facility_id AS facility_id,
        COUNT(*) FILTER (WHERE rac.code = 'BLS' AND m.is_active = TRUE) AS bls_ambulances_count,
        COUNT(*) FILTER (WHERE rac.code = 'ALS' AND m.is_active = TRUE) AS als_ambulances_count,
        COUNT(*) FILTER (WHERE rac.code = 'BOAT' AND m.is_active = TRUE) AS marine_ambulances_count,
        COUNT(*) FILTER (WHERE m.is_active = TRUE) AS total_ambulances_count,
        COUNT(*) FILTER (
            WHERE COALESCE(m.readiness_dispatch_readiness, m.ambulance_dispatch_readiness) = 'DISPATCHABLE'
              AND m.is_active = TRUE
        ) AS dispatchable_ambulances_count
    FROM mv_ambulance_latest_readiness m
    LEFT JOIN ref_ambulance_categories rac ON rac.id = m.category_id
    LEFT JOIN ref_facilities f ON f.id = m.station_facility_id
    GROUP BY 1, 2
)
SELECT
    ids.stat_date,
    ids.district_id,
    ids.facility_id,
    ids.incidents_total,
    ids.red_triage_count,
    ids.transfers_count,
    ids.mnmci_count,
    ids.rta_count,
    ids.infectious_count,
    COALESCE(ast.bls_ambulances_count, 0) AS bls_ambulances_count,
    COALESCE(ast.als_ambulances_count, 0) AS als_ambulances_count,
    COALESCE(ast.marine_ambulances_count, 0) AS marine_ambulances_count,
    COALESCE(ast.total_ambulances_count, 0) AS total_ambulances_count,
    COALESCE(ast.dispatchable_ambulances_count, 0) AS dispatchable_ambulances_count
FROM mv_incident_daily_stats ids
LEFT JOIN ambulance_stats ast
    ON ast.district_id IS NOT DISTINCT FROM ids.district_id
   AND ast.facility_id IS NOT DISTINCT FROM ids.facility_id;

CREATE UNIQUE INDEX uq_mv_dashboard_daily_summary
ON mv_dashboard_daily_summary(stat_date, district_id, facility_id);

CREATE INDEX idx_mv_dashboard_daily_summary_district_id
ON mv_dashboard_daily_summary(district_id);

CREATE INDEX idx_mv_dashboard_daily_summary_facility_id
ON mv_dashboard_daily_summary(facility_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS mv_ambulance_latest_readiness;
DROP MATERIALIZED VIEW IF EXISTS mv_incident_daily_stats;
DROP MATERIALIZED VIEW IF EXISTS mv_dashboard_daily_summary;
-- +goose StatementEnd
