package infrastructure

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	dashboardapp "dispatch/internal/modules/dashboard/application"
	dashboarddomain "dispatch/internal/modules/dashboard/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

var _ dashboardapp.Repository = (*Repository)(nil)

type filterParts struct {
	incidentWhere     string
	availabilityWhere string
	facilityWhere     string
	ambulanceWhere    string
	args              []any
}

func buildFilters(filters dashboarddomain.DashboardFilters) filterParts {
	p := filterParts{
		incidentWhere:     "1=1",
		availabilityWhere: "1=1",
		facilityWhere:     "1=1",
		ambulanceWhere:    "1=1",
		args:              []any{},
	}
	clausesIncident := []string{"1=1"}
	clausesAvailability := []string{"1=1"}
	clausesFacility := []string{"1=1"}
	clausesAmbulance := []string{"1=1"}
	argPos := 1

	if filters.DateFrom != nil {
		clausesIncident = append(clausesIncident, fmt.Sprintf("i.reported_at >= $%d", argPos))
		p.args = append(p.args, *filters.DateFrom)
		argPos++
	}
	if filters.DateTo != nil {
		clausesIncident = append(clausesIncident, fmt.Sprintf("i.reported_at <= $%d", argPos))
		p.args = append(p.args, *filters.DateTo)
		argPos++
	}
	if filters.DistrictID != nil {
		clausesIncident = append(clausesIncident, fmt.Sprintf("i.district_id = $%d", argPos))
		clausesAvailability = append(clausesAvailability, fmt.Sprintf("us.district_id = $%d", argPos))
		clausesFacility = append(clausesFacility, fmt.Sprintf("f.district_id = $%d", argPos))
		clausesAmbulance = append(clausesAmbulance, fmt.Sprintf("a.station_facility_id IN (SELECT id FROM ref_facilities WHERE district_id = $%d)", argPos))
		p.args = append(p.args, *filters.DistrictID)
		argPos++
	}

	if filters.SubcountyID != nil {
		clausesIncident = append(clausesIncident, fmt.Sprintf("i.subcounty_id = $%d", argPos))
		clausesAvailability = append(clausesAvailability, fmt.Sprintf("us.subcounty_id = $%d", argPos))
		clausesFacility = append(clausesFacility, fmt.Sprintf("f.subcounty_id = $%d", argPos))
		clausesAmbulance = append(clausesAmbulance, fmt.Sprintf("a.station_facility_id IN (SELECT id FROM ref_facilities WHERE subcounty_id = $%d)", argPos))
		p.args = append(p.args, *filters.SubcountyID)
		argPos++
	}
	if filters.FacilityID != nil {
		clausesIncident = append(clausesIncident, fmt.Sprintf("i.facility_id = $%d", argPos))
		clausesAvailability = append(clausesAvailability, fmt.Sprintf("us.facility_id = $%d", argPos))
		clausesFacility = append(clausesFacility, fmt.Sprintf("f.id = $%d", argPos))
		clausesAmbulance = append(clausesAmbulance, fmt.Sprintf("a.station_facility_id = $%d", argPos))
		p.args = append(p.args, *filters.FacilityID)
		argPos++
	}
	p.incidentWhere = strings.Join(clausesIncident, " AND ")
	p.availabilityWhere = strings.Join(clausesAvailability, " AND ")
	p.facilityWhere = strings.Join(clausesFacility, " AND ")
	p.ambulanceWhere = strings.Join(clausesAmbulance, " AND ")
	return p
}
func (r *Repository) GetDashboard(ctx context.Context, filters dashboarddomain.DashboardFilters) (dashboarddomain.DashboardResponse, error) {
	resp := dashboarddomain.DashboardResponse{}
	resp.Filters = map[string]interface{}{
		"date_from":   filters.DateFrom,
		"date_to":     filters.DateTo,
		"district_id": filters.DistrictID,
		"facility_id": filters.FacilityID,
	}
	parts := buildFilters(filters)

	if err := r.loadKPIs(ctx, &resp, parts); err != nil {
		return resp, err
	}
	if err := r.loadAmbulanceStatusTable(ctx, &resp, parts); err != nil {
		return resp, err
	}
	if err := r.loadTransfersTrend(ctx, &resp, parts); err != nil {
		return resp, err
	}
	if err := r.loadFacilityCaseTrend(ctx, &resp, parts); err != nil {
		return resp, err
	}
	if err := r.loadCommitteeDonuts(ctx, &resp, parts); err != nil {
		return resp, err
	}
	_ = r.db.QueryRow(ctx, `SELECT MAX(updated_at) FROM incidents`).Scan(&resp.LastUpdatedAt)
	return resp, nil
}

func (r *Repository) loadKPIs(ctx context.Context, resp *dashboarddomain.DashboardResponse, parts filterParts) error {
	// counts by facilities and ambulances
	facilityBase := fmt.Sprintf(`SELECT COUNT(DISTINCT f.id) FROM ref_facilities f WHERE %s`, parts.facilityWhere)
	_ = r.db.QueryRow(ctx, facilityBase, parts.args...).Scan(&resp.KPIs.ConstituenciesCount)

	// Ambulance metrics assume ambulances(category_id, station_facility_id, is_active, functionality_status, vehicle_type)
	qAmb := fmt.Sprintf(`
	SELECT
		COUNT(*) FILTER (WHERE rac.code = 'BLS') AS bls_count,
		COUNT(*) FILTER (WHERE rac.code = 'ALS') AS als_count,
		COUNT(*) AS total_count,
		COUNT(*) FILTER (WHERE rac.code = 'BOAT') AS marine_count
	FROM ambulances a
	LEFT JOIN ref_ambulance_categories rac ON rac.id = a.category_id
	WHERE %s AND COALESCE(a.is_active, TRUE) = TRUE`, parts.ambulanceWhere)
	var bls, als, total, marine int64
	if err := r.db.QueryRow(ctx, qAmb, parts.args...).Scan(&bls, &als, &total, &marine); err != nil {
		return err
	}
	resp.KPIs.BLSAmbulancesCount = bls
	resp.KPIs.ALSAmbulancesCount = als
	if total > 0 {
		resp.KPIs.BLSAmbulancesProportion = (float64(bls) / float64(total)) * 100
		resp.KPIs.ALSAmbulancesProportion = (float64(als) / float64(total)) * 100
		resp.KPIs.MarineAmbulanceProportion = (float64(marine) / float64(total)) * 100
	}

	// Training metrics assume user_training_records(user_id, training_code, completed_at)
	trainingQ := `
	SELECT
		COUNT(DISTINCT user_id) FILTER (WHERE training_code = 'BEC') AS bec,
		COUNT(DISTINCT user_id) FILTER (WHERE training_code = 'BLS') AS bls,
		COUNT(DISTINCT user_id) FILTER (WHERE training_code = 'ALS') AS als,
		COUNT(DISTINCT user_id) FILTER (WHERE training_code = 'CCN') AS ccn,
		COUNT(DISTINCT user_id) FILTER (WHERE training_code = 'EMT') AS emt,
		COUNT(DISTINCT user_id) FILTER (WHERE training_code = 'AMBULANCE_DRIVER') AS drivers
	FROM user_training_records`
	_ = r.db.QueryRow(ctx, trainingQ).Scan(
		&resp.KPIs.HCWsTrainedBEC,
		&resp.KPIs.HCWsTrainedBLS,
		&resp.KPIs.HCWsTrainedALS,
		&resp.KPIs.HCWsTrainedCCN,
		&resp.KPIs.EMTsTrained,
		&resp.KPIs.TrainedAmbulanceDrivers,
	)

	incidentQ := fmt.Sprintf(`
	SELECT
		COUNT(*) FILTER (WHERE i.status IN ('COMPLETED','TRANSPORTING','ENROUTE','AT_SCENE','ASSIGNED')) AS transfers_count,
		COUNT(*) FILTER (WHERE rpl.code = 'RED') AS red_triage_count,
		COUNT(*) FILTER (WHERE it.code = 'HIGHLY_INFECTIOUS' OR i.description ILIKE '%%infectious%%') AS infectious_count,
		COUNT(*) FILTER (WHERE it.code = 'MATERNAL_EMERGENCY') AS mnmci,
		COUNT(*) FILTER (WHERE it.code = 'ACCIDENT' OR i.summary ILIKE '%%rta%%') AS rta
	FROM incidents i
	LEFT JOIN ref_priority_levels rpl ON rpl.id = i.priority_level_id
	LEFT JOIN ref_incident_types it ON it.id = i.incident_type_id
	WHERE %s`, parts.incidentWhere)
	if err := r.db.QueryRow(ctx, incidentQ, parts.args...).Scan(
		&resp.KPIs.TransfersCount,
		&resp.KPIs.RedTriagePatientsCount,
		&resp.KPIs.HighlyInfectiousPatientsCount,
		&resp.KPIs.MNMCI,
		&resp.KPIs.RTA,
	); err != nil {
		return err
	}
	return nil
}

func (r *Repository) loadAmbulanceStatusTable(ctx context.Context, resp *dashboarddomain.DashboardResponse, parts filterParts) error {
	q := fmt.Sprintf(`
	SELECT
		COALESCE(d.name,''),
		COALESCE(f.name,''),
		COALESCE(a.registration_number,''),
		COALESCE(rac.code,''),
		CASE WHEN COALESCE(a.is_functional, FALSE) THEN 'Yes' ELSE 'No' END,
		CASE WHEN COALESCE(a.has_fuel_card, FALSE) THEN 'Yes' ELSE 'No' END,
		CASE WHEN COALESCE(a.has_fuel, FALSE) THEN 'Yes' ELSE 'No' END
	FROM ambulances a
	LEFT JOIN ref_facilities f ON f.id = a.station_facility_id
	LEFT JOIN ref_districts d ON d.id = f.district_id
	LEFT JOIN ref_ambulance_categories rac ON rac.id = a.category_id
	WHERE %s
	ORDER BY d.name, f.name, a.registration_number
	LIMIT 20`, parts.ambulanceWhere)
	rows, err := r.db.Query(ctx, q, parts.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := []dashboarddomain.AmbulanceStatusRow{}
	for rows.Next() {
		var row dashboarddomain.AmbulanceStatusRow
		if err := rows.Scan(&row.District, &row.AmbulanceStation, &row.RegistrationNumber, &row.Category, &row.FunctionalityStatus, &row.FuelCardAvailable, &row.FuelStatus); err != nil {
			return err
		}
		items = append(items, row)
	}
	resp.AmbulanceStatusTable = items
	return rows.Err()
}

func (r *Repository) loadTransfersTrend(ctx context.Context, resp *dashboarddomain.DashboardResponse, parts filterParts) error {
	q := fmt.Sprintf(`
	SELECT to_char(date_trunc('day', i.reported_at), 'YYYY-MM-DD') AS bucket, COUNT(*)
	FROM incidents i
	WHERE %s
	GROUP BY date_trunc('day', i.reported_at)
	ORDER BY date_trunc('day', i.reported_at)
	LIMIT 31`, parts.incidentWhere)
	rows, err := r.db.Query(ctx, q, parts.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := []dashboarddomain.TrendPoint{}
	for rows.Next() {
		var p dashboarddomain.TrendPoint
		if err := rows.Scan(&p.Bucket, &p.Value); err != nil {
			return err
		}
		items = append(items, p)
	}
	resp.TransfersTrend = items
	return rows.Err()
}

func (r *Repository) loadFacilityCaseTrend(ctx context.Context, resp *dashboarddomain.DashboardResponse, parts filterParts) error {
	q := fmt.Sprintf(`
	SELECT COALESCE(f.name, 'Unknown') AS bucket, COUNT(*) AS total_cases
	FROM incidents i
	LEFT JOIN ref_facilities f ON f.id = i.facility_id
	WHERE %s
	GROUP BY f.name
	ORDER BY total_cases DESC
	LIMIT 12`, parts.incidentWhere)
	rows, err := r.db.Query(ctx, q, parts.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := []dashboarddomain.TrendPoint{}
	for rows.Next() {
		var p dashboarddomain.TrendPoint
		if err := rows.Scan(&p.Bucket, &p.Value); err != nil {
			return err
		}
		items = append(items, p)
	}
	resp.FacilityCaseTrend = items
	return rows.Err()
}

func (r *Repository) loadCommitteeDonuts(ctx context.Context, resp *dashboarddomain.DashboardResponse, parts filterParts) error {
	// assumes ambulance_committees and ambulance_financing tables exist; fallback safe if absent should be implemented in production
	_ = r.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE is_functional = TRUE) AS yes_count,
			COUNT(*) FILTER (WHERE is_functional = FALSE) AS no_count
		FROM ambulance_committees`).Scan(&resp.AmbulanceCommittees.Yes, &resp.AmbulanceCommittees.No)
	_ = r.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE is_financed_by_llu = TRUE) AS yes_count,
			COUNT(*) FILTER (WHERE is_financed_by_llu = FALSE) AS no_count
		FROM ambulance_financing`).Scan(&resp.AmbulanceLLUFinancing.Yes, &resp.AmbulanceLLUFinancing.No)
	return nil
}
