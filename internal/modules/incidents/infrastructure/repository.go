package infrastructure

import (
	"context"
	"fmt"
	"strings"

	"dispatch/internal/modules/incidents/application"
	"dispatch/internal/modules/incidents/domain"
	platformdb "dispatch/internal/platform/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var _ application.Repository = (*Repository)(nil)

func (r *Repository) ListIncidents(ctx context.Context, p platformdb.Pagination) ([]domain.Incident, int64, error) {
	allowedSorts := map[string]string{
		"created_at":  "i.created_at",
		"reported_at": "i.reported_at",
		"status":      "i.status",
	}

	where := []string{"1=1"}
	args := make([]any, 0)
	argPos := 1

	if p.Search != "" {
		where = append(where, fmt.Sprintf(`(
			LOWER(i.incident_number) LIKE LOWER($%d) OR
			LOWER(COALESCE(i.caller_name,'')) LIKE LOWER($%d) OR
			LOWER(COALESCE(i.patient_name,'')) LIKE LOWER($%d) OR
			LOWER(COALESCE(i.summary,'')) LIKE LOWER($%d)
		)`, argPos, argPos, argPos, argPos))
		args = append(args, "%"+p.Search+"%")
		argPos++
	}

	for key, value := range p.Filters {
		switch key {
		case "status":
			where = append(where, fmt.Sprintf("i.status = $%d", argPos))
			args = append(args, strings.ToUpper(value))
			argPos++
		case "verification_status":
			where = append(where, fmt.Sprintf("i.verification_status = $%d", argPos))
			args = append(args, strings.ToUpper(value))
			argPos++
		case "incident_type_id":
			where = append(where, fmt.Sprintf("i.incident_type_id = $%d", argPos))
			args = append(args, value)
			argPos++
		case "district_id":
			where = append(where, fmt.Sprintf("i.district_id = $%d", argPos))
			args = append(args, value)
			argPos++
		case "facility_id":
			where = append(where, fmt.Sprintf("i.facility_id = $%d", argPos))
			args = append(args, value)
			argPos++
		case "date_from":
			where = append(where, fmt.Sprintf("i.reported_at >= $%d", argPos))
			args = append(args, value)
			argPos++
		case "date_to":
			where = append(where, fmt.Sprintf("i.reported_at <= $%d", argPos))
			args = append(args, value)
			argPos++
		}
	}

	whereSQL := "WHERE " + strings.Join(where, " AND ")

	countSQL := fmt.Sprintf(`SELECT COUNT(1) FROM incidents i %s`, whereSQL)

	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := platformdb.BuildOrderBy(p, allowedSorts)

	listSQL := fmt.Sprintf(`
		SELECT
			i.id,
			i.incident_number,
			i.source_channel,
			i.caller_name,
			i.caller_phone,
			i.patient_name,
			i.patient_phone,
			i.patient_age_group,
			i.patient_sex,
			i.incident_type_id,
			i.severity_level_id,
			i.priority_level_id,
			i.summary,
			i.description,
			i.district_id,
			i.facility_id,
			i.village,
			i.parish,
			i.subcounty,
			i.landmark,
			i.latitude,
			i.longitude,
			i.verification_status,
			i.status,
			i.reported_at,
			i.created_by_user_id,
			i.triaged_by_user_id,
			i.triaged_at,
			i.assigned_at,
			i.closed_at,
			i.created_at,
			i.updated_at
		FROM incidents i
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereSQL, orderBy, argPos, argPos+1)

	rows, err := r.db.Query(ctx, listSQL, append(args, p.PageSize, p.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.Incident, 0)
	for rows.Next() {
		var in domain.Incident
		if err := rows.Scan(
			&in.ID,
			&in.IncidentNumber,
			&in.SourceChannel,
			&in.CallerName,
			&in.CallerPhone,
			&in.PatientName,
			&in.PatientPhone,
			&in.PatientAgeGroup,
			&in.PatientSex,
			&in.IncidentTypeID,
			&in.SeverityLevelID,
			&in.PriorityLevelID,
			&in.Summary,
			&in.Description,
			&in.DistrictID,
			&in.FacilityID,
			&in.Village,
			&in.Parish,
			&in.Subcounty,
			&in.Landmark,
			&in.Latitude,
			&in.Longitude,
			&in.VerificationStatus,
			&in.Status,
			&in.ReportedAt,
			&in.CreatedByUserID,
			&in.TriagedByUserID,
			&in.TriagedAt,
			&in.AssignedAt,
			&in.ClosedAt,
			&in.CreatedAt,
			&in.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, in)
	}

	return items, total, rows.Err()
}

func (r *Repository) CreateIncident(ctx context.Context, in domain.Incident) (domain.Incident, error) {
	query := `
		INSERT INTO incidents (
			id, incident_number, source_channel,
			caller_name, caller_phone,
			patient_name, patient_phone, patient_age_group, patient_sex,
			incident_type_id, severity_level_id, priority_level_id,
			summary, description,
			district_id, facility_id, village, parish, subcounty, landmark,
			latitude, longitude,
			verification_status, status,
			reported_at, created_by_user_id, triaged_by_user_id, triaged_at,
			assigned_at, closed_at, created_at, updated_at
		) VALUES (
			$1,$2,$3,
			$4,$5,
			$6,$7,$8,$9,
			$10,$11,$12,
			$13,$14,
			$15,$16,$17,$18,$19,$20,
			$21,$22,
			$23,$24,
			$25,$26,$27,$28,
			$29,$30,$31,$32
		)
		RETURNING created_at, updated_at
	`
	if err := r.db.QueryRow(
		ctx,
		query,
		in.ID,
		in.IncidentNumber,
		in.SourceChannel,
		in.CallerName,
		in.CallerPhone,
		in.PatientName,
		in.PatientPhone,
		in.PatientAgeGroup,
		in.PatientSex,
		in.IncidentTypeID,
		in.SeverityLevelID,
		in.PriorityLevelID,
		in.Summary,
		in.Description,
		in.DistrictID,
		in.FacilityID,
		in.Village,
		in.Parish,
		in.Subcounty,
		in.Landmark,
		in.Latitude,
		in.Longitude,
		in.VerificationStatus,
		in.Status,
		in.ReportedAt,
		in.CreatedByUserID,
		in.TriagedByUserID,
		in.TriagedAt,
		in.AssignedAt,
		in.ClosedAt,
		in.CreatedAt,
		in.UpdatedAt,
	).Scan(&in.CreatedAt, &in.UpdatedAt); err != nil {
		return domain.Incident{}, err
	}
	return in, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (domain.Incident, error) {
	const q = `
SELECT
	i.id,
	i.incident_number,
	i.source_channel,
	i.caller_name,
	i.caller_phone,
	i.patient_name,
	i.patient_phone,
	i.patient_age_group,
	i.patient_sex,
	i.incident_type_id,
	i.severity_level_id,
	i.priority_level_id,
	i.summary,
	i.description,
	i.district_id,
	i.facility_id,
	i.village,
	i.parish,
	i.subcounty,
	i.landmark,
	i.latitude,
	i.longitude,
	i.verification_status,
	i.status,
	i.reported_at,
	i.created_by_user_id,
	i.triaged_by_user_id,
	i.triaged_at,
	i.assigned_at,
	i.closed_at,
	i.created_at,
	i.updated_at
FROM incidents i
WHERE i.id = $1`
	var in domain.Incident
	if err := r.db.QueryRow(ctx, q, id).Scan(
		&in.ID,
		&in.IncidentNumber,
		&in.SourceChannel,
		&in.CallerName,
		&in.CallerPhone,
		&in.PatientName,
		&in.PatientPhone,
		&in.PatientAgeGroup,
		&in.PatientSex,
		&in.IncidentTypeID,
		&in.SeverityLevelID,
		&in.PriorityLevelID,
		&in.Summary,
		&in.Description,
		&in.DistrictID,
		&in.FacilityID,
		&in.Village,
		&in.Parish,
		&in.Subcounty,
		&in.Landmark,
		&in.Latitude,
		&in.Longitude,
		&in.VerificationStatus,
		&in.Status,
		&in.ReportedAt,
		&in.CreatedByUserID,
		&in.TriagedByUserID,
		&in.TriagedAt,
		&in.AssignedAt,
		&in.ClosedAt,
		&in.CreatedAt,
		&in.UpdatedAt,
	); err != nil {
		return domain.Incident{}, err
	}
	return in, nil
}

func (r *Repository) UpdateIncident(ctx context.Context, id string, req application.UpdateIncidentRequest) (domain.Incident, error) {
	sets := make([]string, 0)
	args := make([]any, 0)
	pos := 1

	if req.CallerName != nil {
		sets = append(sets, fmt.Sprintf("caller_name = $%d", pos))
		args = append(args, *req.CallerName)
		pos++
	}
	if req.CallerPhone != nil {
		sets = append(sets, fmt.Sprintf("caller_phone = $%d", pos))
		args = append(args, *req.CallerPhone)
		pos++
	}
	if req.PatientName != nil {
		sets = append(sets, fmt.Sprintf("patient_name = $%d", pos))
		args = append(args, *req.PatientName)
		pos++
	}
	if req.PatientPhone != nil {
		sets = append(sets, fmt.Sprintf("patient_phone = $%d", pos))
		args = append(args, *req.PatientPhone)
		pos++
	}
	if req.PatientAgeGroup != nil {
		sets = append(sets, fmt.Sprintf("patient_age_group = $%d", pos))
		args = append(args, *req.PatientAgeGroup)
		pos++
	}
	if req.PatientSex != nil {
		sets = append(sets, fmt.Sprintf("patient_sex = $%d", pos))
		args = append(args, *req.PatientSex)
		pos++
	}
	if req.IncidentTypeID != nil {
		sets = append(sets, fmt.Sprintf("incident_type_id = $%d", pos))
		args = append(args, *req.IncidentTypeID)
		pos++
	}
	if req.SeverityLevelID != nil {
		sets = append(sets, fmt.Sprintf("severity_level_id = $%d", pos))
		args = append(args, *req.SeverityLevelID)
		pos++
	}
	if req.PriorityLevelID != nil {
		sets = append(sets, fmt.Sprintf("priority_level_id = $%d", pos))
		args = append(args, *req.PriorityLevelID)
		pos++
	}
	if req.Summary != nil {
		sets = append(sets, fmt.Sprintf("summary = $%d", pos))
		args = append(args, *req.Summary)
		pos++
	}
	if req.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", pos))
		args = append(args, *req.Description)
		pos++
	}
	if req.DistrictID != nil {
		sets = append(sets, fmt.Sprintf("district_id = $%d", pos))
		args = append(args, *req.DistrictID)
		pos++
	}
	if req.FacilityID != nil {
		sets = append(sets, fmt.Sprintf("facility_id = $%d", pos))
		args = append(args, *req.FacilityID)
		pos++
	}
	if req.Village != nil {
		sets = append(sets, fmt.Sprintf("village = $%d", pos))
		args = append(args, *req.Village)
		pos++
	}
	if req.Parish != nil {
		sets = append(sets, fmt.Sprintf("parish = $%d", pos))
		args = append(args, *req.Parish)
		pos++
	}
	if req.Subcounty != nil {
		sets = append(sets, fmt.Sprintf("subcounty = $%d", pos))
		args = append(args, *req.Subcounty)
		pos++
	}
	if req.Landmark != nil {
		sets = append(sets, fmt.Sprintf("landmark = $%d", pos))
		args = append(args, *req.Landmark)
		pos++
	}
	if req.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", pos))
		args = append(args, strings.ToUpper(*req.Status))
		pos++
	}
	if req.Verification != nil {
		sets = append(sets, fmt.Sprintf("verification_status = $%d", pos))
		args = append(args, strings.ToUpper(*req.Verification))
		pos++
	}
	if len(sets) == 0 {
		return r.GetByID(ctx, id)
	}
	sets = append(sets, "updated_at = now()")
	args = append(args, id)
	query := fmt.Sprintf("UPDATE incidents SET %s WHERE id = $%d", strings.Join(sets, ", "), pos)
	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		return domain.Incident{}, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) DeleteIncident(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM incidents WHERE id = $1`, id)
	return err
}

// ResolveIncidentTypeID resolves either a UUID, a code, or a name (case-insensitive)
// into a concrete incident_type id from ref_incident_types.
func (r *Repository) ResolveIncidentTypeID(ctx context.Context, value string) (string, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return "", fmt.Errorf("incident_type_id is required")
	}

	// If it's already a valid UUID and exists, use it directly.
	var id string
	if _, err := uuid.Parse(v); err == nil {
		if err := r.db.QueryRow(ctx, `SELECT id FROM ref_incident_types WHERE id = $1`, v).Scan(&id); err == nil {
			return id, nil
		}
		// fall through to code/name lookup if not found by id
	}

	// Try by code, then by name (case-insensitive)
	err := r.db.QueryRow(ctx,
		`SELECT id FROM ref_incident_types WHERE code = $1 OR LOWER(name) = LOWER($1) LIMIT 1`,
		v,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("unknown incident_type %q", v)
	}
	return id, nil
}

// ResolveDistrictID resolves a human-friendly district identifier into ref_districts.id.
// Accepts UUID (id), or name/code; empty string results in nil.
func (r *Repository) ResolveDistrictID(ctx context.Context, value string) (*string, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil, nil
	}

	var id string
	// If it's a UUID and exists, use it.
	if _, err := uuid.Parse(v); err == nil {
		if err := r.db.QueryRow(ctx, `SELECT id FROM ref_districts WHERE id = $1`, v).Scan(&id); err == nil {
			return &id, nil
		}
	}

	// Try by code or name.
	err := r.db.QueryRow(ctx,
		`SELECT id FROM ref_districts WHERE code = $1 OR LOWER(name) = LOWER($1) LIMIT 1`,
		v,
	).Scan(&id)
	if err != nil {
		// If we can't resolve, return nil (leave district unset) rather than failing the whole request.
		return nil, nil
	}
	return &id, nil
}

// ResolveFacilityID resolves facility identifier into ref_facilities.id.
// Accepts UUID (id), or code; empty string results in nil. Name resolution is not attempted
// because facility names are frequently non-unique.
func (r *Repository) ResolveFacilityID(ctx context.Context, value string) (*string, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil, nil
	}

	var id string
	// If it's a UUID and exists, use it.
	if _, err := uuid.Parse(v); err == nil {
		if err := r.db.QueryRow(ctx, `SELECT id FROM ref_facilities WHERE id = $1`, v).Scan(&id); err == nil {
			return &id, nil
		}
	}

	// Try by code (our facility UID)
	err := r.db.QueryRow(ctx,
		`SELECT id FROM ref_facilities WHERE code = $1 LIMIT 1`,
		v,
	).Scan(&id)
	if err != nil {
		// If we can't resolve, leave facility unset.
		return nil, nil
	}
	return &id, nil
}
