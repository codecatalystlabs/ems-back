package infrastructure

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"dispatch/internal/modules/availability/application/dto"
	availabilitydomain "dispatch/internal/modules/availability/domain"
	platformdb "dispatch/internal/platform/db"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) CreateShift(ctx context.Context, in availabilitydomain.UserShift) (availabilitydomain.UserShift, error) {
	q := `
	INSERT INTO user_shifts (
		id, user_id, shift_date, starts_at, ends_at, shift_type, district_id, facility_id, status, created_by
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, q,
		in.ID, in.UserID, in.ShiftDate, in.StartsAt, in.EndsAt, in.ShiftType, in.DistrictID, in.FacilityID, in.Status, in.CreatedBy,
	).Scan(&in.CreatedAt, &in.UpdatedAt)
	return in, err
}

func (r *Repository) UpdateShift(ctx context.Context, id string, req dto.UpdateShiftRequest) (availabilitydomain.UserShift, error) {
	sets := make([]string, 0)
	args := make([]any, 0)
	pos := 1
	if req.ShiftDate != nil {
		sets = append(sets, fmt.Sprintf("shift_date=$%d", pos))
		args = append(args, *req.ShiftDate)
		pos++
	}
	if req.StartsAt != nil {
		sets = append(sets, fmt.Sprintf("starts_at=$%d", pos))
		args = append(args, *req.StartsAt)
		pos++
	}
	if req.EndsAt != nil {
		sets = append(sets, fmt.Sprintf("ends_at=$%d", pos))
		args = append(args, *req.EndsAt)
		pos++
	}
	if req.ShiftType != nil {
		sets = append(sets, fmt.Sprintf("shift_type=$%d", pos))
		args = append(args, *req.ShiftType)
		pos++
	}
	if req.DistrictID != nil {
		sets = append(sets, fmt.Sprintf("district_id=$%d", pos))
		args = append(args, *req.DistrictID)
		pos++
	}
	if req.FacilityID != nil {
		sets = append(sets, fmt.Sprintf("facility_id=$%d", pos))
		args = append(args, *req.FacilityID)
		pos++
	}
	if req.Status != nil {
		sets = append(sets, fmt.Sprintf("status=$%d", pos))
		args = append(args, *req.Status)
		pos++
	}
	if len(sets) == 0 {
		return r.GetShiftByID(ctx, id)
	}
	sets = append(sets, "updated_at=now()")
	args = append(args, id)
	q := fmt.Sprintf("UPDATE user_shifts SET %s WHERE id=$%d", strings.Join(sets, ", "), pos)
	_, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return availabilitydomain.UserShift{}, err
	}
	return r.GetShiftByID(ctx, id)
}

func (r *Repository) GetShiftByID(ctx context.Context, id string) (availabilitydomain.UserShift, error) {
	var out availabilitydomain.UserShift
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, shift_date::text, starts_at, ends_at, COALESCE(shift_type,''), district_id, facility_id, status, created_by, created_at, updated_at
		FROM user_shifts WHERE id=$1`, id,
	).Scan(&out.ID, &out.UserID, &out.ShiftDate, &out.StartsAt, &out.EndsAt, &out.ShiftType, &out.DistrictID, &out.FacilityID, &out.Status, &out.CreatedBy, &out.CreatedAt, &out.UpdatedAt)
	return out, err
}

func (r *Repository) ListShifts(ctx context.Context, params dto.ListShiftsParams) ([]availabilitydomain.UserShift, int64, error) {
	p := params.Pagination
	allowedSorts := map[string]string{"shift_date": "us.shift_date", "starts_at": "us.starts_at", "status": "us.status", "created_at": "us.created_at"}
	where := []string{"1=1"}
	args := []any{}
	pos := 1
	if params.UserID != nil && *params.UserID != "" {
		where = append(where, fmt.Sprintf("us.user_id=$%d", pos))
		args = append(args, *params.UserID)
		pos++
	}
	if params.DistrictID != nil && *params.DistrictID != "" {
		where = append(where, fmt.Sprintf("us.district_id=$%d", pos))
		args = append(args, *params.DistrictID)
		pos++
	}
	if params.FacilityID != nil && *params.FacilityID != "" {
		where = append(where, fmt.Sprintf("us.facility_id=$%d", pos))
		args = append(args, *params.FacilityID)
		pos++
	}
	if params.ShiftDate != nil && *params.ShiftDate != "" {
		where = append(where, fmt.Sprintf("us.shift_date=$%d", pos))
		args = append(args, *params.ShiftDate)
		pos++
	}
	if params.Status != nil && *params.Status != "" {
		where = append(where, fmt.Sprintf("us.status=$%d", pos))
		args = append(args, strings.ToUpper(*params.Status))
		pos++
	}
	if p.Search != "" {
		where = append(where, fmt.Sprintf("(COALESCE(us.shift_type,'') ILIKE $%d)", pos))
		args = append(args, "%"+p.Search+"%")
		pos++
	}
	whereSQL := "WHERE " + strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(1) FROM user_shifts us `+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := fmt.Sprintf(`SELECT id, user_id, shift_date::text, starts_at, ends_at, COALESCE(shift_type,''), district_id, facility_id, status, created_by, created_at, updated_at FROM user_shifts us %s %s LIMIT $%d OFFSET $%d`, whereSQL, platformdb.BuildOrderBy(p, allowedSorts), pos, pos+1)
	rows, err := r.db.Query(ctx, q, append(args, p.PageSize, p.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []availabilitydomain.UserShift{}
	for rows.Next() {
		var out availabilitydomain.UserShift
		if err := rows.Scan(&out.ID, &out.UserID, &out.ShiftDate, &out.StartsAt, &out.EndsAt, &out.ShiftType, &out.DistrictID, &out.FacilityID, &out.Status, &out.CreatedBy, &out.CreatedAt, &out.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, out)
	}
	return items, total, rows.Err()
}

func (r *Repository) UpsertAvailability(ctx context.Context, in availabilitydomain.UserAvailability) (availabilitydomain.UserAvailability, error) {
	q := `
	INSERT INTO user_availability (
		id, user_id, availability_status, dispatchable, current_incident_id, current_dispatch_assignment_id,
		current_ambulance_id, last_seen_at, source, notes, updated_by, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	ON CONFLICT (user_id) DO UPDATE SET
		availability_status = EXCLUDED.availability_status,
		dispatchable = EXCLUDED.dispatchable,
		current_incident_id = EXCLUDED.current_incident_id,
		current_dispatch_assignment_id = EXCLUDED.current_dispatch_assignment_id,
		current_ambulance_id = EXCLUDED.current_ambulance_id,
		last_seen_at = EXCLUDED.last_seen_at,
		source = EXCLUDED.source,
		notes = EXCLUDED.notes,
		updated_by = EXCLUDED.updated_by,
		updated_at = EXCLUDED.updated_at
	RETURNING id, user_id, availability_status, dispatchable, current_incident_id, current_dispatch_assignment_id,
		current_ambulance_id, last_seen_at, source, COALESCE(notes,''), updated_by, updated_at`
	var out availabilitydomain.UserAvailability
	err := r.db.QueryRow(ctx, q,
		in.ID, in.UserID, in.AvailabilityStatus, in.Dispatchable, in.CurrentIncidentID, in.CurrentDispatchAssignmentID,
		in.CurrentAmbulanceID, in.LastSeenAt, in.Source, in.Notes, in.UpdatedBy, in.UpdatedAt,
	).Scan(&out.ID, &out.UserID, &out.AvailabilityStatus, &out.Dispatchable, &out.CurrentIncidentID, &out.CurrentDispatchAssignmentID,
		&out.CurrentAmbulanceID, &out.LastSeenAt, &out.Source, &out.Notes, &out.UpdatedBy, &out.UpdatedAt)
	return out, err
}

func (r *Repository) GetAvailabilityByUserID(ctx context.Context, userID string) (availabilitydomain.UserAvailability, error) {
	var out availabilitydomain.UserAvailability
	err := r.db.QueryRow(ctx, `SELECT id, user_id, availability_status, dispatchable, current_incident_id, current_dispatch_assignment_id, current_ambulance_id, last_seen_at, source, COALESCE(notes,''), updated_by, updated_at FROM user_availability WHERE user_id=$1`, userID).
		Scan(&out.ID, &out.UserID, &out.AvailabilityStatus, &out.Dispatchable, &out.CurrentIncidentID, &out.CurrentDispatchAssignmentID, &out.CurrentAmbulanceID, &out.LastSeenAt, &out.Source, &out.Notes, &out.UpdatedBy, &out.UpdatedAt)
	return out, err
}

func (r *Repository) ListAvailability(ctx context.Context, params dto.ListAvailabilityParams) ([]availabilitydomain.UserAvailability, int64, error) {
	p := params.Pagination
	allowedSorts := map[string]string{"updated_at": "ua.updated_at", "availability_status": "ua.availability_status", "last_seen_at": "ua.last_seen_at"}
	where := []string{"1=1"}
	args := []any{}
	pos := 1
	if params.Status != nil && *params.Status != "" {
		where = append(where, fmt.Sprintf("ua.availability_status=$%d", pos))
		args = append(args, strings.ToUpper(*params.Status))
		pos++
	}
	if params.Dispatchable != nil {
		where = append(where, fmt.Sprintf("ua.dispatchable=$%d", pos))
		args = append(args, *params.Dispatchable)
		pos++
	}
	whereSQL := "WHERE " + strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(1) FROM user_availability ua `+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := fmt.Sprintf(`SELECT id, user_id, availability_status, dispatchable, current_incident_id, current_dispatch_assignment_id, current_ambulance_id, last_seen_at, source, COALESCE(notes,''), updated_by, updated_at FROM user_availability ua %s %s LIMIT $%d OFFSET $%d`, whereSQL, platformdb.BuildOrderBy(p, allowedSorts), pos, pos+1)
	rows, err := r.db.Query(ctx, q, append(args, p.PageSize, p.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []availabilitydomain.UserAvailability{}
	for rows.Next() {
		var out availabilitydomain.UserAvailability
		if err := rows.Scan(&out.ID, &out.UserID, &out.AvailabilityStatus, &out.Dispatchable, &out.CurrentIncidentID, &out.CurrentDispatchAssignmentID, &out.CurrentAmbulanceID, &out.LastSeenAt, &out.Source, &out.Notes, &out.UpdatedBy, &out.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, out)
	}
	return items, total, rows.Err()
}

func (r *Repository) CreatePresenceLog(ctx context.Context, in availabilitydomain.UserPresenceLog) (availabilitydomain.UserPresenceLog, error) {
	q := `INSERT INTO user_presence_logs (id, user_id, channel, seen_at, ip_address, user_agent, latitude, longitude) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING seen_at`
	err := r.db.QueryRow(ctx, q, in.ID, in.UserID, in.Channel, in.SeenAt, in.IPAddress, in.UserAgent, in.Latitude, in.Longitude).Scan(&in.SeenAt)
	return in, err
}

func (r *Repository) ListPresenceLogs(ctx context.Context, params dto.ListPresenceParams) ([]availabilitydomain.UserPresenceLog, int64, error) {
	p := params.Pagination
	allowedSorts := map[string]string{"seen_at": "upl.seen_at", "channel": "upl.channel"}
	where := []string{"1=1"}
	args := []any{}
	pos := 1
	if params.UserID != nil && *params.UserID != "" {
		where = append(where, fmt.Sprintf("upl.user_id=$%d", pos))
		args = append(args, *params.UserID)
		pos++
	}
	if params.Channel != nil && *params.Channel != "" {
		where = append(where, fmt.Sprintf("upl.channel=$%d", pos))
		args = append(args, strings.ToUpper(*params.Channel))
		pos++
	}
	whereSQL := "WHERE " + strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(1) FROM user_presence_logs upl `+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := fmt.Sprintf(`SELECT id, user_id, channel, seen_at, host(ip_address)::text, COALESCE(user_agent,''), latitude, longitude FROM user_presence_logs upl %s %s LIMIT $%d OFFSET $%d`, whereSQL, platformdb.BuildOrderBy(p, allowedSorts), pos, pos+1)
	rows, err := r.db.Query(ctx, q, append(args, p.PageSize, p.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []availabilitydomain.UserPresenceLog{}
	for rows.Next() {
		var out availabilitydomain.UserPresenceLog
		if err := rows.Scan(&out.ID, &out.UserID, &out.Channel, &out.SeenAt, &out.IPAddress, &out.UserAgent, &out.Latitude, &out.Longitude); err != nil {
			return nil, 0, err
		}
		items = append(items, out)
	}
	return items, total, rows.Err()
}
