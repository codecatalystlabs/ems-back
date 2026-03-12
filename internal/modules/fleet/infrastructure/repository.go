package infrastructure

import (
	"context"
	"fmt"
	"strings"

	fleetapp "dispatch/internal/modules/fleet/application"
	"dispatch/internal/modules/fleet/domain"
	platformdb "dispatch/internal/platform/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var _ fleetapp.Repository = (*Repository)(nil)

func (r *Repository) ListAmbulances(ctx context.Context, p platformdb.Pagination) ([]domain.Ambulance, int64, error) {
	allowedSorts := map[string]string{
		"created_at":         "a.created_at",
		"plate_number":       "a.plate_number",
		"status":             "a.status",
		"dispatch_readiness": "a.dispatch_readiness",
	}

	where := []string{"1=1"}
	args := make([]any, 0)
	argPos := 1

	if p.Search != "" {
		where = append(where, fmt.Sprintf(`(
			COALESCE(a.code,'') ILIKE $%d OR
			a.plate_number ILIKE $%d OR
			COALESCE(a.vin,'') ILIKE $%d OR
			COALESCE(a.make,'') ILIKE $%d OR
			COALESCE(a.model,'') ILIKE $%d
		)`, argPos, argPos, argPos, argPos, argPos))
		args = append(args, "%"+p.Search+"%")
		argPos++
	}

	for key, value := range p.Filters {
		switch key {
		case "status":
			where = append(where, fmt.Sprintf("a.status = $%d", argPos))
			args = append(args, strings.ToUpper(value))
			argPos++
		case "dispatch_readiness":
			where = append(where, fmt.Sprintf("a.dispatch_readiness = $%d", argPos))
			args = append(args, strings.ToUpper(value))
			argPos++
		case "district_id":
			where = append(where, fmt.Sprintf("a.district_id = $%d", argPos))
			args = append(args, value)
			argPos++
		case "category_id":
			where = append(where, fmt.Sprintf("a.category_id = $%d", argPos))
			args = append(args, value)
			argPos++
		case "date_from":
			where = append(where, fmt.Sprintf("a.created_at >= $%d", argPos))
			args = append(args, value)
			argPos++
		case "date_to":
			where = append(where, fmt.Sprintf("a.created_at <= $%d", argPos))
			args = append(args, value)
			argPos++
		}
	}

	whereSQL := "WHERE " + strings.Join(where, " AND ")

	countSQL := fmt.Sprintf(`SELECT COUNT(1) FROM ambulances a %s`, whereSQL)

	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := platformdb.BuildOrderBy(p, allowedSorts)

	listSQL := fmt.Sprintf(`
		SELECT
			a.id,
			COALESCE(a.code, ''),
			a.plate_number,
			COALESCE(a.vin, ''),
			COALESCE(a.make, ''),
			COALESCE(a.model, ''),
			a.year_of_manufacture,
			a.category_id,
			COALESCE(a.ownership_type, ''),
			a.station_facility_id,
			a.district_id,
			a.status,
			a.dispatch_readiness,
			a.gps_lat,
			a.gps_lon,
			a.last_seen_at,
			a.is_active,
			a.created_at,
			a.updated_at
		FROM ambulances a
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereSQL, orderBy, argPos, argPos+1)

	rows, err := r.db.Query(ctx, listSQL, append(args, p.PageSize, p.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.Ambulance, 0)
	for rows.Next() {
		var a domain.Ambulance
		var code, vin, makeVal, model, ownershipType *string
		if err := rows.Scan(
			&a.ID,
			&code,
			&a.PlateNumber,
			&vin,
			&makeVal,
			&model,
			&a.YearOfManufacture,
			&a.CategoryID,
			&ownershipType,
			&a.StationFacilityID,
			&a.DistrictID,
			&a.Status,
			&a.DispatchReadiness,
			&a.GPSLat,
			&a.GPSLon,
			&a.LastSeenAt,
			&a.IsActive,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		a.Code = code
		a.VIN = vin
		a.Make = makeVal
		a.Model = model
		a.OwnershipType = ownershipType
		items = append(items, a)
	}

	return items, total, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id string) (domain.Ambulance, error) {
	const q = `
SELECT
	a.id,
	a.code,
	a.plate_number,
	a.vin,
	a.make,
	a.model,
	a.year_of_manufacture,
	a.category_id,
	a.ownership_type,
	a.station_facility_id,
	a.district_id,
	a.status,
	a.dispatch_readiness,
	a.gps_lat,
	a.gps_lon,
	a.last_seen_at,
	a.is_active,
	a.created_at,
	a.updated_at
FROM ambulances a
WHERE a.id = $1`
	var a domain.Ambulance
	if err := r.db.QueryRow(ctx, q, id).Scan(
		&a.ID,
		&a.Code,
		&a.PlateNumber,
		&a.VIN,
		&a.Make,
		&a.Model,
		&a.YearOfManufacture,
		&a.CategoryID,
		&a.OwnershipType,
		&a.StationFacilityID,
		&a.DistrictID,
		&a.Status,
		&a.DispatchReadiness,
		&a.GPSLat,
		&a.GPSLon,
		&a.LastSeenAt,
		&a.IsActive,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		return domain.Ambulance{}, err
	}
	return a, nil
}

func (r *Repository) Create(ctx context.Context, in domain.Ambulance) (domain.Ambulance, error) {
	const q = `
INSERT INTO ambulances (
	id, code, plate_number, vin, make, model, year_of_manufacture,
	category_id, ownership_type, station_facility_id, district_id,
	status, dispatch_readiness, gps_lat, gps_lon, location, last_seen_at,
	is_active, created_at, updated_at
)
VALUES (
	gen_random_uuid(), $1,$2,$3,$4,$5,$6,
	$7,$8,$9,$10,
	$11,$12,NULL,NULL,NULL,NULL,
	TRUE, now(), now()
)
RETURNING id`
	var id string
	if err := r.db.QueryRow(
		ctx,
		q,
		in.Code,
		in.PlateNumber,
		in.VIN,
		in.Make,
		in.Model,
		in.YearOfManufacture,
		in.CategoryID,
		in.OwnershipType,
		in.StationFacilityID,
		in.DistrictID,
		in.Status,
		in.DispatchReadiness,
	).Scan(&id); err != nil {
		return domain.Ambulance{}, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) Update(ctx context.Context, id string, req fleetapp.UpdateAmbulanceRequest) (domain.Ambulance, error) {
	sets := make([]string, 0)
	args := make([]any, 0)
	pos := 1

	if req.Code != nil {
		sets = append(sets, fmt.Sprintf("code = $%d", pos))
		args = append(args, *req.Code)
		pos++
	}
	if req.VIN != nil {
		sets = append(sets, fmt.Sprintf("vin = $%d", pos))
		args = append(args, *req.VIN)
		pos++
	}
	if req.Make != nil {
		sets = append(sets, fmt.Sprintf("make = $%d", pos))
		args = append(args, *req.Make)
		pos++
	}
	if req.Model != nil {
		sets = append(sets, fmt.Sprintf("model = $%d", pos))
		args = append(args, *req.Model)
		pos++
	}
	if req.YearOfManufacture != nil {
		sets = append(sets, fmt.Sprintf("year_of_manufacture = $%d", pos))
		args = append(args, *req.YearOfManufacture)
		pos++
	}
	if req.CategoryID != nil {
		sets = append(sets, fmt.Sprintf("category_id = $%d", pos))
		args = append(args, *req.CategoryID)
		pos++
	}
	if req.OwnershipType != nil {
		sets = append(sets, fmt.Sprintf("ownership_type = $%d", pos))
		args = append(args, *req.OwnershipType)
		pos++
	}
	if req.StationFacilityID != nil {
		sets = append(sets, fmt.Sprintf("station_facility_id = $%d", pos))
		args = append(args, *req.StationFacilityID)
		pos++
	}
	if req.DistrictID != nil {
		sets = append(sets, fmt.Sprintf("district_id = $%d", pos))
		args = append(args, *req.DistrictID)
		pos++
	}
	if req.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", pos))
		args = append(args, strings.ToUpper(*req.Status))
		pos++
	}
	if req.DispatchReadiness != nil {
		sets = append(sets, fmt.Sprintf("dispatch_readiness = $%d", pos))
		args = append(args, strings.ToUpper(*req.DispatchReadiness))
		pos++
	}
	if len(sets) == 0 {
		return r.GetByID(ctx, id)
	}
	sets = append(sets, "updated_at = now()")
	args = append(args, id)
	query := fmt.Sprintf("UPDATE ambulances SET %s WHERE id = $%d", strings.Join(sets, ", "), pos)
	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		return domain.Ambulance{}, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM ambulances WHERE id = $1`, id)
	return err
}
